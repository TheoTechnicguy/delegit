/**
 * file: database/main.go
 * author: theo technicguy
 * license: apache-2.0
 *
 * This file contains generic database utility parts
 * for the application.
 */

// The database package implements the logic
// for the data persistance plane of the application.
// It should be as extensible as possible, so that
// it can be used with any database type.
package database

import (
	"errors"

	"git.licolas.net/delegit/delegit/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrInvalidDatabaseKind error = errors.New("invalid database type")
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(kind, dsn string) (*Database, error) {
	var dialect gorm.Dialector
	switch kind {
	case "sqlite":
		dialect = sqlite.Open(dsn)
	case "pgsql":
		dialect = postgres.Open(dsn)
	default:
		return nil, ErrInvalidDatabaseKind
	}

	db, err := NewDatabaseFromDialector(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate()
	return db, err
}

func NewDatabaseFromDialector(dialect gorm.Dialector, config *gorm.Config) (*Database, error) {
	db := new(Database)

	var err error
	db.db, err = gorm.Open(dialect, config)
	if err != nil {
		return nil, err
	}

	return db, err
}

func (db *Database) AutoMigrate() error {
	t := []any{
		models.Feedback{},
	}

	for _, v := range t {
		if err := db.db.AutoMigrate(&v); err != nil {
			return err
		}
	}

	return nil
}
