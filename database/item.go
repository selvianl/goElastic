package database

import (
	"insider/models"
)

func (d *Database) ListItems(filters map[string]interface{}, page, limit uint16) (uint16, []models.ItemOutput, error) {
	var apps []models.ItemOutput
	var count int64

	if err := d.db.Where(filters).Limit(int(limit)).Offset(int(page)).Find(&apps).Count(&count).Error; err != nil {
		return 0, nil, err
	}

	return uint16(count), apps, nil
}
