package main

import (
	"context"
	"flag"
	"insider/api"
	"insider/config"
	"insider/database"
	"insider/elastic"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

func main() {
	var flagAddr string

	flag.StringVar(&flagAddr, "s", ":8092", "HTTP server listen address")

	flagLogLevel := zap.LevelFlag("l", zapcore.DebugLevel, "Log level")
	flag.Parse()

	// More setup
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	c, err := config.Parse()
	if err != nil {
		log.Fatal("init config failed", zap.Error(err))
		os.Exit(1)
	}

	// log setup
	zc := zap.NewDevelopmentConfig()
	zc.Level = zap.NewAtomicLevelAt(*flagLogLevel)
	zc.OutputPaths = []string{c.Log.Output}

	zlog, err := zc.Build()
	if err != nil {
		panic(err)
	}
	log = zlog

	wg := &sync.WaitGroup{}
	wg.Add(1) // for main

	// start elastic
	es, err := elastic.New(c.IndexName, c.FilePath, c.ElasticUrl)
	if err != nil {
		log.Fatal("init elastic failed", zap.Error(err))
		os.Exit(1)
	}

	// start db
	db, err := database.New(ctx, log.With(zap.String("module", "database")), c)
	if err != nil {
		log.Fatal("init database failed", zap.Error(err))
		os.Exit(1)
	}

	go api.RunAPI(ctx, wg, log.With(zap.String("module", "api")), es, db, c, flagAddr)

	go func() {
		// Wait for all tasks to finish
		<-ctx.Done()
		wg.Done() // for main
	}()

	wg.Wait()
}
