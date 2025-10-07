package integration

import (
	"testing"
	"time"

	"go-seo/internal/domain/entities"
	"go-seo/internal/infrastructure/database/postgres/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestKeywordRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repositories.NewKeywordRepository(gormDB)

	// Ожидаем INSERT запрос
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "keywords"`).
		WithArgs("купить чай", 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	keyword := &entities.Keyword{
		Value:  "купить чай",
		SiteID: 1,
	}

	err = repo.Create(keyword)
	assert.NoError(t, err)
	assert.Equal(t, 1, keyword.ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestKeywordRepository_GetByValueAndSite(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repositories.NewKeywordRepository(gormDB)

	// Ожидаем SELECT запрос
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "value", "site_id", "created_at", "updated_at"}).
		AddRow(1, "купить чай", 1, now, now)

	mock.ExpectQuery(`SELECT \* FROM "keywords" WHERE value = \$1 AND site_id = \$2 ORDER BY "keywords"\."id" LIMIT \$3`).
		WithArgs("купить чай", 1, 1).
		WillReturnRows(rows)

	keyword, err := repo.GetByValueAndSite("купить чай", 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, keyword.ID)
	assert.Equal(t, "купить чай", keyword.Value)
	assert.Equal(t, 1, keyword.SiteID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestKeywordRepository_GetBySiteID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repositories.NewKeywordRepository(gormDB)

	// Ожидаем SELECT запрос
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "value", "site_id", "created_at", "updated_at"}).
		AddRow(1, "купить чай", 1, now, now).
		AddRow(2, "купить кофе", 1, now, now)

	mock.ExpectQuery(`SELECT \* FROM "keywords" WHERE site_id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	keywords, err := repo.GetBySiteID(1)
	assert.NoError(t, err)
	assert.Len(t, keywords, 2)
	assert.Equal(t, "купить чай", keywords[0].Value)
	assert.Equal(t, "купить кофе", keywords[1].Value)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestKeywordRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repositories.NewKeywordRepository(gormDB)

	// Ожидаем DELETE запрос
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "keywords" WHERE "keywords"\."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete(1)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
