//go:build integration

package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iciantoine/todo-go-api/model"
	"github.com/iciantoine/todo-go-api/repository"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

const pgTimestamptzHourFormat = "2006-01-02 15:04:05.999999999Z07"

func TestGetTodos(t *testing.T) {
	t.Run("it should return a list of todos", func(t *testing.T) {
		SUT, teardown := setup(t)
		defer teardown()

		res, err := SUT.GetTodos(context.Background())
		assert.NoError(t, err)
		assert.Len(t, res, 2)

		assert.Equal(t, uuid.MustParse("169e84e3-35d9-4476-8295-2c28c54d50fc"), res[0].ID)
		cDate, _ := time.Parse(pgTimestamptzHourFormat, "2023-03-06 14:00:00.000000+00")
		assert.Equal(t, cDate.Local(), res[0].CreatedAt)
		assert.True(t, res[0].IsDone)
		assert.Equal(t, "Lorem ipsum", res[0].Message)

		assert.Equal(t, uuid.MustParse("038863e4-2fbe-4bc3-9e38-1e62e93659f5"), res[1].ID)
		cDate, _ = time.Parse(pgTimestamptzHourFormat, "2023-03-06 12:00:00.000000+00")
		assert.Equal(t, cDate.Local(), res[1].CreatedAt)
		assert.False(t, res[1].IsDone)
		assert.Equal(t, "Test", res[1].Message)
	})
}

func TestGetTodo(t *testing.T) {
	t.Run("it should return a todo", func(t *testing.T) {
		SUT, teardown := setup(t)
		defer teardown()

		res, err := SUT.GetTodo(context.Background(), uuid.MustParse("038863e4-2fbe-4bc3-9e38-1e62e93659f5"))
		assert.NoError(t, err)
		cDate, _ := time.Parse(pgTimestamptzHourFormat, "2023-03-06 12:00:00.000000+00")
		assert.Equal(t, cDate.Local(), res.CreatedAt)
		assert.False(t, res.IsDone)
		assert.Equal(t, "Test", res.Message)
	})

	t.Run("it should return a not found error", func(t *testing.T) {
		SUT, teardown := setup(t)
		defer teardown()

		_, err := SUT.GetTodo(context.Background(), uuid.New())
		assert.ErrorIs(t, repository.ErrTodoNotFound, err)
	})
}

func TestAddTodo(t *testing.T) {
	t.Run("it should add a todo and return it", func(t *testing.T) {
		SUT, teardown := setup(t)
		defer teardown()

		res, err := SUT.AddTodo(context.Background(), model.Todo{
			IsDone:  true,
			Message: "test",
		})
		assert.NoError(t, err)

		assert.NotEmpty(t, res.ID.String())
		assert.True(t, res.IsDone)
		assert.Equal(t, "test", res.Message)
		assert.NotZero(t, res.CreatedAt)
	})
}

func setup(t *testing.T) (repository.TodoRepo, func()) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=todo password=todo dbname=todo sslmode=disable")
	assert.NoError(t, err)

	tx, err := db.BeginTx(context.Background(), nil)
	assert.NoError(t, err)

	SUT := repository.NewTodoRepo(tx)

	return SUT, func() {
		assert.NoError(t, tx.Rollback())
		assert.NoError(t, db.Close())
	}
}
