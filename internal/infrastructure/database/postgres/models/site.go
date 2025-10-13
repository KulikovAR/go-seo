package models

import "time"

type Site struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Domain    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
