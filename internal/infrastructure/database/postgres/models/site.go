package models

import "time"

type Site struct {
	ID            int       `gorm:"primaryKey;autoIncrement"`
	Domain        string    `gorm:"not null"`
	YandexDynamic *int      `gorm:"type:smallint;default:null"`
	GoogleDynamic *int      `gorm:"type:smallint;default:null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
