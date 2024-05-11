package api

import (
	"context"
	"insider/config"
	"insider/database"
	"insider/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type mockClient struct {
	Api      API
	Es       *mockElasticClient
	Database *database.Database
}

type mockElasticClient struct {
	DoSearchFunc func(ctxBg context.Context, queryBytes []byte) models.SearchResponse
}

func newMockElasticClient() *mockElasticClient {
	return &mockElasticClient{}
}

func (me *mockElasticClient) DoSearch(ctxBg context.Context, queryBytes []byte) models.SearchResponse {
	if me.DoSearchFunc != nil {
		return me.DoSearchFunc(ctxBg, queryBytes)
	}
	return models.SearchResponse{}
}

func GetMockClient(conn gorm.Dialector) (*mockClient, error) {
	ctx := context.Background()
	cfg := &config.Config{}
	esClient := newMockElasticClient()

	var db *database.Database
	var err error
	if conn != nil {
		db, err = database.NewWithConnection(ctx, conn, cfg)
		if err != nil {
			return nil, err
		}
	}

	api := API{
		ec:  echo.New(),
		cfg: cfg,
		val: validator.New(),
		db:  db,
		ctx: ctx,
		es:  esClient,
	}

	api.setupHandlers()
	api.setupSwagger()

	return &mockClient{
		Api:      api,
		Es:       esClient,
		Database: db,
	}, nil
}
