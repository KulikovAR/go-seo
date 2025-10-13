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

func TestSiteRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repositories.NewSiteRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "sites"`).
		WithArgs("test.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	site := &entities.Site{
		Domain: "test.com",
	}

	err = repo.Create(site)
	assert.NoError(t, err)
	assert.Equal(t, 1, site.ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSiteRepository_GetByDomain(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repositories.NewSiteRepository(gormDB)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "domain", "created_at", "updated_at"}).
		AddRow(1, "test.com", now, now)

	mock.ExpectQuery(`SELECT \* FROM "sites" WHERE domain = \$1 ORDER BY "sites"\."id" LIMIT \$2`).
		WithArgs("test.com", 1).
		WillReturnRows(rows)

	site, err := repo.GetByDomain("test.com")
	assert.NoError(t, err)
	assert.Equal(t, 1, site.ID)
	assert.Equal(t, "test.com", site.Domain)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSiteRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repositories.NewSiteRepository(gormDB)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "domain", "created_at", "updated_at"}).
		AddRow(1, "test1.com", now, now).
		AddRow(2, "test2.com", now, now)

	mock.ExpectQuery(`SELECT \* FROM "sites"`).
		WillReturnRows(rows)

	sites, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, sites, 2)
	assert.Equal(t, "test1.com", sites[0].Domain)
	assert.Equal(t, "test2.com", sites[1].Domain)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSiteRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repositories.NewSiteRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "sites" WHERE "sites"\."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete(1)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
