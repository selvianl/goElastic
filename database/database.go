package database

import (
	"context"
	"insider/config"
	"insider/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db  *gorm.DB
	ctx context.Context
	cfg *config.Config
	log *zap.Logger
}

func New(ctx context.Context, log *zap.Logger, cfg *config.Config) (*Database, error) {
	goCfg := gorm.Config{}
	goCfg.Logger = logger.Default.LogMode(logger.Info)

	conn, err := gorm.Open(postgres.Open(cfg.Database), &goCfg)
	if err != nil {
		return nil, err
	}

	db := &Database{
		db:  conn,
		ctx: ctx,
		log: log,
		cfg: cfg,
	}

	if err := db.migrate(); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Database) migrate() error {
	m := []interface{}{
		&models.ItemOutput{},
		&models.SortConfig{},
	}

	return d.db.AutoMigrate(m...)
}
