package migrations

import (
	"go-seo/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Site{},
		&models.Keyword{},
		&models.Position{},
	)
}

func CreateTables(db *gorm.DB) error {
	if err := AutoMigrate(db); err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_positions_keyword_site_date 
		ON positions (keyword_id, site_id, date DESC);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_positions_site_date 
		ON positions (site_id, date DESC);
	`).Error; err != nil {
		return err
	}

	return nil
}
