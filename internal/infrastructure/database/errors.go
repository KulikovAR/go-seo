package database

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// DatabaseError представляет ошибку базы данных
type DatabaseError struct {
	Code    string
	Message string
	Err     error
}

func (e *DatabaseError) Error() string {
	return e.Message
}

// WrapDatabaseError оборачивает ошибку GORM в пользовательскую ошибку
func WrapDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	// Проверяем, является ли это ошибкой GORM
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &DatabaseError{
			Code:    "NOT_FOUND",
			Message: "Record not found",
			Err:     err,
		}
	}

	// Проверяем на нарушение ограничений
	errStr := err.Error()

	// Нарушение уникального ключа
	if strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "UNIQUE constraint") {
		return &DatabaseError{
			Code:    "DUPLICATE_ENTRY",
			Message: "Record already exists",
			Err:     err,
		}
	}

	// Нарушение внешнего ключа
	if strings.Contains(errStr, "foreign key constraint") || strings.Contains(errStr, "violates foreign key") {
		return &DatabaseError{
			Code:    "FOREIGN_KEY_VIOLATION",
			Message: "Referenced record does not exist",
			Err:     err,
		}
	}

	// Нарушение NOT NULL
	if strings.Contains(errStr, "NOT NULL constraint") || strings.Contains(errStr, "null value") {
		return &DatabaseError{
			Code:    "NULL_CONSTRAINT",
			Message: "Required field cannot be empty",
			Err:     err,
		}
	}

	// Общая ошибка базы данных
	return &DatabaseError{
		Code:    "DATABASE_ERROR",
		Message: "Database operation failed",
		Err:     err,
	}
}

// IsDatabaseError проверяет, является ли ошибка DatabaseError
func IsDatabaseError(err error) bool {
	_, ok := err.(*DatabaseError)
	return ok
}

// GetDatabaseErrorCode возвращает код ошибки базы данных
func GetDatabaseErrorCode(err error) string {
	if de, ok := err.(*DatabaseError); ok {
		return de.Code
	}
	return ""
}
