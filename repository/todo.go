package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iciantoine/todo-go-api/model"
)

var ErrTodoNotFound = errors.New("todo not found")

// TodoRepo is the todo repository.
type TodoRepo struct {
	db DBTX
}

// NewTodoRepo instantiates TodoRepo.
func NewTodoRepo(db DBTX) TodoRepo {
	return TodoRepo{
		db: db,
	}
}

// GetTodos gets all the todos ordered by creation date from most newest to oldest.
func (repo TodoRepo) GetTodos(ctx context.Context) ([]model.Todo, error) {
	const q = `
		SELECT id, created_at, is_done, message
		FROM todo
		ORDER BY created_at DESC
	`

	return list(scan)(repo.db.QueryContext(ctx, q))
}

// GetTodos retrives one todo by its ID or throws an error.
func (repo TodoRepo) GetTodo(ctx context.Context, id uuid.UUID) (model.Todo, error) {
	const q = `
		SELECT id, created_at, is_done, message
		FROM todo
		WHERE id = $1
	`

	res, err := scan(repo.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return res, ErrTodoNotFound
	}

	return res, err
}

// AddTodo adds a todo model.
func (repo TodoRepo) AddTodo(ctx context.Context, model model.Todo) (model.Todo, error) {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()

	const q = `
		INSERT INTO todo (id, created_at, is_done, message)
		VALUES ($1, $2, $3, $4)
	`

	_, err := repo.db.ExecContext(ctx, q, model.ID, model.CreatedAt, model.IsDone, model.Message)

	return model, err
}

func scan(row scanner) (model.Todo, error) {
	var val model.Todo
	err := row.Scan(&val.ID, &val.CreatedAt, &val.IsDone, &val.Message)
	return val, err
}
