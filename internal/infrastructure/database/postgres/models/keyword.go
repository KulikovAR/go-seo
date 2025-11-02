package models

import "time"

type Keyword struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Value     string    `gorm:"not null"`
	SiteID    int       `gorm:"not null;index"`
	GroupID   int       `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Site  Site  `gorm:"foreignKey:SiteID"`
	Group Group `gorm:"foreignKey:GroupID"`
}
