package database

import (
	"context"
	"insider/config"

	"gorm.io/gorm"
)

type mockDb struct {
	Database *Database
}

func GetMockDb(conn gorm.Dialector) (*mockDb, error) {
	ctx := context.Background()
	cfg := &config.Config{}

	var db *Database
	var err error
	if conn != nil {
		db, err = NewWithConnection(ctx, conn, cfg)
		if err != nil {
			return nil, err
		}
	}

	return &mockDb{
		Database: db,
	}, nil
}
