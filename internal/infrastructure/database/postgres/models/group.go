package models

import "time"

type Group struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
