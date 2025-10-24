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
		&models.TrackingJob{},
		&models.TrackingTask{},
		&models.TrackingResult{},
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

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_positions_keyword_site_source 
		ON positions (keyword_id, site_id, source);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_positions_source_date 
		ON positions (source, date DESC);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_jobs_site_id 
		ON tracking_jobs (site_id);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_jobs_status 
		ON tracking_jobs (status);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_tasks_job_id 
		ON tracking_tasks (job_id);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_tasks_keyword_id 
		ON tracking_tasks (keyword_id);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_tasks_site_id 
		ON tracking_tasks (site_id);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_tasks_status 
		ON tracking_tasks (status);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_results_task_id 
		ON tracking_results (task_id);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_results_job_id 
		ON tracking_results (job_id);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_results_keyword_id 
		ON tracking_results (keyword_id);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracking_results_site_id 
		ON tracking_results (site_id);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_positions_stats_main 
		ON positions (site_id, source, date DESC, rank);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_positions_trends 
		ON positions (keyword_id, date DESC, rank) 
		WHERE date >= CURRENT_DATE - INTERVAL '30 days';
	`).Error; err != nil {
		return err
	}

	return nil
}
