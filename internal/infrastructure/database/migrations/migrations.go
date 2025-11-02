package migrations

import (
	"go-seo/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Site{},
		&models.Group{},
		&models.Keyword{},
		&models.Position{},
		&models.TrackingJob{},
		&models.TrackingTask{},
		&models.TrackingResult{},
	)
}

func CreateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Group{}); err != nil {
		return err
	}

	var hasGroupIDColumn bool
	if err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'keywords' 
			AND column_name = 'group_id'
		)
	`).Scan(&hasGroupIDColumn).Error; err != nil {
		return err
	}

	getOrCreateDefaultGroup := func() (int, error) {
		var groupID int
		if err := db.Raw(`SELECT id FROM groups WHERE name = 'Default' LIMIT 1`).Scan(&groupID).Error; err != nil || groupID == 0 {
			if err := db.Exec(`INSERT INTO groups (name, created_at, updated_at) VALUES ('Default', NOW(), NOW())`).Error; err != nil {
				return 0, err
			}
			if err := db.Raw(`SELECT id FROM groups WHERE name = 'Default' LIMIT 1`).Scan(&groupID).Error; err != nil {
				return 0, err
			}
		}
		return groupID, nil
	}

	if !hasGroupIDColumn {
		if err := db.Exec(`ALTER TABLE keywords ADD COLUMN group_id INTEGER`).Error; err != nil {
			return err
		}

		defaultGroupID, err := getOrCreateDefaultGroup()
		if err != nil {
			return err
		}

		if err := db.Exec(`UPDATE keywords SET group_id = ? WHERE group_id IS NULL`, defaultGroupID).Error; err != nil {
			return err
		}
	} else {
		var nullCount int64
		if err := db.Raw(`SELECT COUNT(*) FROM keywords WHERE group_id IS NULL`).Scan(&nullCount).Error; err != nil {
			return err
		}
		if nullCount > 0 {
			defaultGroupID, err := getOrCreateDefaultGroup()
			if err != nil {
				return err
			}
			if err := db.Exec(`UPDATE keywords SET group_id = ? WHERE group_id IS NULL`, defaultGroupID).Error; err != nil {
				return err
			}
		}
	}

	// Проверяем и исправляем столбец domain в tracking_tasks
	var hasDomainColumn bool
	if err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'tracking_tasks' 
			AND column_name = 'domain'
		)
	`).Scan(&hasDomainColumn).Error; err != nil {
		return err
	}

	if hasDomainColumn {
		// Проверяем тип столбца
		var columnType string
		if err := db.Raw(`
			SELECT data_type 
			FROM information_schema.columns 
			WHERE table_name = 'tracking_tasks' 
			AND column_name = 'domain'
		`).Scan(&columnType).Error; err != nil {
			return err
		}

		// Если столбец имеет тип character varying (varchar), удаляем его
		if columnType == "character varying" || columnType == "varchar" {
			if err := db.Exec(`ALTER TABLE tracking_tasks DROP COLUMN IF EXISTS domain`).Error; err != nil {
				return err
			}
		}
	}

	if err := db.AutoMigrate(
		&models.Site{},
		&models.Keyword{},
		&models.Position{},
		&models.TrackingJob{},
		&models.TrackingTask{},
		&models.TrackingResult{},
	); err != nil {
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

	//if err := db.Exec(`
	//	CREATE INDEX IF NOT EXISTS idx_positions_trends
	//	ON positions (keyword_id, date DESC, rank)
	//	WHERE date >= CURRENT_DATE - INTERVAL '30 days';
	//`).Error; err != nil {
	//	return err
	//} - пока без этого, ждем переезда на 16 версию

	return nil
}
