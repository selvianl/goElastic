package database

import (
	"insider/models"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

func (d *Database) GetActiveConfig() (*models.SortConfig, error) {
	var config models.SortConfig
	if err := d.db.Where("is_active = ?", true).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

func (d *Database) CreateConfig(in models.ConfigInput) (*models.SortConfig, error) {
	conf := models.SortConfig{
		Id:         ksuid.New().String(),
		SortOption: in.SortOption,
		SortOrder:  in.SortOrder,
		IsActive:   in.IsActive,
	}

	if err := d.db.Create(&conf).Error; err != nil {
		return nil, err
	}

	return &conf, nil
}

func (d *Database) ListConfigs(offset, limit int) (int64, []models.SortConfig, error) {
	var confs []models.SortConfig
	var count int64

	if err := d.db.Find(&confs).Limit(limit).Offset(offset).Count(&count).Error; err != nil {
		return 0, nil, err
	}

	return count, confs, nil
}

func (d *Database) GetConfig(id string) (*models.SortConfig, error) {
	var cfg models.SortConfig
	if err := d.db.First(&cfg, "lower(id) = lower(?)", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (d *Database) UpdateConfig(sc *models.SortConfig, in models.ConfigInputUpdate) (*models.SortConfig, error) {
	sc.IsActive = in.IsActive
	if in.SortOption != "" {
		sc.SortOption = in.SortOption
	}
	if in.SortOrder != "" {
		sc.SortOrder = in.SortOrder
	}

	if err := d.db.Save(sc).Error; err != nil {
		return nil, err
	}

	return sc, nil
}

func (d *Database) DeleteConfig(sc *models.SortConfig) error {
	return d.db.Delete(sc).Error
}
