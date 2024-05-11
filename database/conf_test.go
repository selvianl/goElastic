package database

import (
	"insider/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
)

type MockConfigInput struct {
	SortOption string
	SortOrder  string
	IsActive   bool
}

type ConfigInputUpdate struct {
	SortOption string
	SortOrder  string
	IsActive   bool
}

func TestConfigFunctions_Success(t *testing.T) {
	input := MockConfigInput{
		SortOption: "field1",
		SortOrder:  "asc",
		IsActive:   true,
	}

	conn := sqlite.Open("gorm.db")
	defer func() {
		err := os.Remove("gorm.db")
		assert.NoError(t, err)
	}()

	mc, err := GetMockDb(conn)
	if err != nil {
		assert.Fail(t, "failed")
	}

	cfg, err := mc.Database.CreateConfig(models.ConfigInput(input))
	if err != nil {
		assert.Fail(t, "failed")
	}
	assert.Equal(t, cfg.SortOption, input.SortOption)

	id := cfg.Id

	cfg, err = mc.Database.GetConfig(id)
	if err != nil {
		assert.Fail(t, "failed")
	}
	assert.Equal(t, cfg.Id, id)

	newInput := MockConfigInput{
		SortOption: "changed",
		SortOrder:  "changed",
		IsActive:   true,
	}

	cfg, err = mc.Database.UpdateConfig(cfg, models.ConfigInputUpdate(newInput))
	if err != nil {
		assert.Fail(t, "failed")
	}
	assert.Equal(t, cfg.SortOption, newInput.SortOption)

	activeCfg, err := mc.Database.GetActiveConfig()
	if err != nil {
		assert.Fail(t, "failed")
	}
	assert.Equal(t, cfg.Id, activeCfg.Id)

	if err := mc.Database.DeleteConfig(cfg); err != nil {
		assert.Fail(t, "failed")
	}

	expectedCount := 0
	c, _, err := mc.Database.ListConfigs(0, 20)
	if err != nil {
		assert.Fail(t, "failed")
	}
	assert.Equal(t, expectedCount, int(c))

}
