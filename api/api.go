package api

import (
	"context"
	"insider/config"
	"insider/database"
	"insider/elastic"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type API struct {
	ec  *echo.Echo
	ctx context.Context
	cfg *config.Config
	val *validator.Validate
	db  *database.Database
	es  *elastic.Elastic
}

func RunAPI(ctx context.Context, wg *sync.WaitGroup, log *zap.Logger, es *elastic.Elastic, db *database.Database, cfg *config.Config, addr string) {
	srv := &API{
		ec:  echo.New(),
		ctx: ctx,
		cfg: cfg,
		db:  db,
		val: validator.New(),
		es:  es,
	}
	srv.ec.Server.IdleTimeout = 120 * time.Second

	srv.setupHandlers()
	srv.setupSwagger()

	// Start server
	wg.Add(1)
	defer wg.Done()

	go func() {
		<-ctx.Done()
		// When the app is shutinsiderg down, shutdown the HTTP server as well
		_ = srv.ec.Shutdown(context.Background())
	}()

	if err := srv.ec.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Error("http", zap.Error(err))
	} else {
		log.Info("Shutinsiderg down", zap.Error(err))
	}
}
