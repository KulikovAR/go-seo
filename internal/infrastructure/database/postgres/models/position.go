package models

import "time"

type Position struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	KeywordID int       `gorm:"not null;index"`
	SiteID    int       `gorm:"not null;index"`
	Rank      int       `gorm:"not null"`
	URL       string    `gorm:"not null"`
	Title     string    `gorm:"not null"`
	Source    string    `gorm:"not null;index"`
	Device    string    `gorm:"not null"`
	OS        string    `gorm:""`
	Ads       bool      `gorm:"not null"`
	Country   string    `gorm:""`
	Lang      string    `gorm:""`
	Pages     int       `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Keyword Keyword `gorm:"foreignKey:KeywordID"`
	Site    Site    `gorm:"foreignKey:SiteID"`
}
