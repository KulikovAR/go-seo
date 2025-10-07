package models

import "time"

type Site struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`
	Domain    string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
