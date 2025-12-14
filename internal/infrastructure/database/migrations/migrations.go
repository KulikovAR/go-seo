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

	// Удаляем колонку name из таблицы sites, если она существует
	var hasNameColumn bool
	if err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'sites' 
			AND column_name = 'name'
		)
	`).Scan(&hasNameColumn).Error; err != nil {
		return err
	}

	if hasNameColumn {
		if err := db.Exec(`ALTER TABLE sites DROP COLUMN IF EXISTS name`).Error; err != nil {
			return err
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

	var hasSiteIDColumn bool
	if err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'groups' 
			AND column_name = 'site_id'
		)
	`).Scan(&hasSiteIDColumn).Error; err != nil {
		return err
	}

	if !hasSiteIDColumn {
		if err := db.Exec(`ALTER TABLE groups ADD COLUMN site_id INTEGER NOT NULL DEFAULT 1`).Error; err != nil {
			return err
		}
		if err := db.Exec(`ALTER TABLE groups ALTER COLUMN site_id DROP DEFAULT`).Error; err != nil {
			return err
		}
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_groups_site_id 
		ON groups (site_id);
	`).Error; err != nil {
		return err
	}

	// Удаляем уникальное ограничение на поле domain в таблице sites, если оно существует
	// Это позволяет создавать сайты с одинаковыми доменами
	var constraintName string
	if err := db.Raw(`
		SELECT conname 
		FROM pg_constraint 
		WHERE conrelid = 'sites'::regclass 
		AND contype = 'u'
		AND pg_get_constraintdef(oid) LIKE '%domain%'
		LIMIT 1
	`).Scan(&constraintName).Error; err != nil {
		// Игнорируем ошибку, если таблица еще не существует или ограничение не найдено
	}

	if constraintName != "" {
		if err := db.Exec(`ALTER TABLE sites DROP CONSTRAINT IF EXISTS ` + constraintName).Error; err != nil {
			return err
		}
	}

	// Также удаляем уникальные индексы на domain, если они существуют
	var indexName string
	if err := db.Raw(`
		SELECT indexname 
		FROM pg_indexes 
		WHERE schemaname = 'public'
		AND tablename = 'sites' 
		AND indexdef LIKE '%UNIQUE%'
		AND indexdef LIKE '%domain%'
		LIMIT 1
	`).Scan(&indexName).Error; err != nil {
		// Игнорируем ошибку, если индекс не найден
	}

	if indexName != "" {
		if err := db.Exec(`DROP INDEX IF EXISTS ` + indexName).Error; err != nil {
			return err
		}
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
