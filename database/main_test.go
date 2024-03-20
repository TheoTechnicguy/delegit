/**
 * file: database/main_test.go
 * author: theo technicguy
 * license: apache-2.0
 */

package database_test

import (
	"testing"

	"git.licolas.net/delegit/delegit/database"
	"github.com/stretchr/testify/assert"
)

func TestNewDatabaseSQLite(t *testing.T) {
	dir := t.TempDir()
	db, err := database.NewDatabase("sqlite", dir+"/test.db")

	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestNewDatabasePGSQL(t *testing.T) {
	db, err := database.NewDatabase("pgsql", "host=localhost user=test password=test dbname=test port=5432 sslmode=disable")
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestNewDatabaseDialectorError(t *testing.T) {
	db, err := database.NewDatabase("invalid", "invalid")
	assert.Error(t, err)
	assert.ErrorIs(t, err, database.ErrInvalidDatabaseKind)
	assert.Nil(t, db)
}

func TestNewDatabaseDSNError(t *testing.T) {
	db, err := database.NewDatabase("pgsql", "host=localhost user=bad password=test dbname=test port=5432 sslmode=disable")
	assert.Error(t, err)
	assert.Nil(t, db)
}
