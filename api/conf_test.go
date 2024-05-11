package api

import (
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

func TestAPI_listConfs(t *testing.T) {
	conn := sqlite.Open("gorm.db")
	defer func() {
		err := os.Remove("gorm.db")
		assert.NoError(t, err)
	}()

	mc, err := GetMockClient(conn)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/configs", nil)
	rec := httptest.NewRecorder()

	mc.Api.ec.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAPI_createConfig(t *testing.T) {
	conn := sqlite.Open("gorm.db")
	defer func() {
		err := os.Remove("gorm.db")
		assert.NoError(t, err)
	}()

	mc, err := GetMockClient(conn)
	assert.NoError(t, err)

	body := models.ConfigInput{
		SortOption: "test",
		SortOrder:  "test",
		IsActive:   true,
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/configs", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	mc.Api.ec.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var res models.SortConfig
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	dbConn := mc.Api.db.GetConnection()
	var resDb models.SortConfig
	if err := dbConn.Where(
		"lower(id) = lower(?)", res.Id).First(
		&resDb).Error; err != nil {
		assert.Fail(t, "cfg is not created")
	}
	assert.Equal(t, resDb.Id, res.Id)
}
