//go:build integration

package database_test

import (
	"testing"

	"github.com/iciantoine/todo-go-api/database"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	t.Run("successful connection", func(t *testing.T) {
		_, err := database.Connect("pgx", "host=localhost user=todo password=todo")
		assert.NoError(t, err)
	})

	t.Run("error during connection", func(t *testing.T) {
		conn, err := database.Connect("pgx", "")
		assert.Nil(t, conn)
		assert.Errorf(t, err, "could not open DB connection pool: %w")
	})
}
