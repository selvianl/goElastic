package models

import "time"

type SortConfig struct {
	Id         string    `gorm:"primaryKey;autoIncrement:false" json:"id"`
	SortOption string    `gorm:"not null" json:"sort_option"`
	SortOrder  string    `gorm:"not null" json:"sort_order"`
	IsActive   bool      `gorm:"not null" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
