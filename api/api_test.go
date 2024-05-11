package api

import (
	"context"
	"encoding/json"
	"insider/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
)

func TestAPI_doSearch(t *testing.T) {
	conn := sqlite.Open("gorm.db")
	defer func() {
		err := os.Remove("gorm.db")
		assert.NoError(t, err)
	}()

	mc, err := GetMockClient(conn)
	assert.NoError(t, err)

	mc.Es.DoSearchFunc = func(ctxBg context.Context, queryBytes []byte) models.SearchResponse {
		return models.SearchResponse{
			Hits: struct {
				Total struct {
					Value int `json:"value"`
				} `json:"total"`
				Hits []struct {
					Source map[string]interface{} `json:"_source"`
				} `json:"hits"`
			}{
				Total: struct {
					Value int `json:"value"`
				}{
					Value: 2,
				},
				Hits: []struct {
					Source map[string]interface{} `json:"_source"`
				}{
					{
						Source: map[string]interface{}{
							"item_id":  "399",
							"name":     "Single",
							"locale":   "tr_TR",
							"click":    float64(122),
							"purchase": float64(904),
						},
					},
					{
						Source: map[string]interface{}{
							"item_id":  "1086",
							"name":     "Woo Album #4",
							"locale":   "tr_TR",
							"click":    float64(203),
							"purchase": float64(606),
						},
					},
				},
			},
		}
	}
	expectedLen := 2

	body := models.FilterParams{
		Conditions: []models.FilterCondition{
			{FieldName: "name", Operation: "query", Value: "ogo"},
			{FieldName: "click", Operation: "lt", Value: "280"},
		},
	}

	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/items/", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	mc.Api.ec.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp = &struct {
		Results []models.ItemOutput
		Count   int
	}{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, resp.Count, expectedLen)

	for _, item := range resp.Results {
		isValid := item.Name == "Woo Album #4" || item.Name == "Single"
		assert.True(t, isValid)
	}
}
