package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iciantoine/todo-go-api/handler"
	"github.com/iciantoine/todo-go-api/model"
	"github.com/iciantoine/todo-go-api/repository"
	"github.com/stretchr/testify/assert"
)

type stubRepo struct {
	todo     model.Todo
	todoList []model.Todo
	err      error
}

func (sr *stubRepo) GetTodos(ctx context.Context) ([]model.Todo, error) {
	return sr.todoList, sr.err
}

func (sr *stubRepo) GetTodo(ctx context.Context, id uuid.UUID) (model.Todo, error) {
	return sr.todo, sr.err
}

func (sr *stubRepo) AddTodo(ctx context.Context, model model.Todo) (model.Todo, error) {
	return sr.todo, sr.err
}

func TestNewGetTodosHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 on successful empty list call", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request, _ = http.NewRequest("GET", "/todo", http.NoBody)

		expected := []model.Todo{}
		r, err := json.Marshal(expected)
		assert.NoError(t, err)

		hdlr := handler.NewGetTodosHandler(&stubRepo{
			todoList: expected,
		})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, r, body)
	})

	t.Run("returns 200 on successful non empty list call", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request, _ = http.NewRequest("GET", "/todo", http.NoBody)

		expected := []model.Todo{
			{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				IsDone:    true,
				Message:   "Lorem ipsum",
			},
			{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				IsDone:    false,
				Message:   "Test",
			},
		}
		r, err := json.Marshal(expected)
		assert.NoError(t, err)

		hdlr := handler.NewGetTodosHandler(&stubRepo{
			todoList: expected,
		})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, r, body)
	})

	t.Run("returns 500 on list call with repo error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request, _ = http.NewRequest("GET", "/todo", http.NoBody)

		hdlr := handler.NewGetTodosHandler(&stubRepo{
			err: errors.New("test"),
		})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Empty(t, body)
	})

	t.Run("returns 200 on successful todo call", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		uuid := uuid.New()
		ctx.Request, _ = http.NewRequest("GET", fmt.Sprintf("/todo?id=%s", uuid.String()), http.NoBody)

		expected := model.Todo{
			ID:        uuid,
			CreatedAt: time.Now(),
			IsDone:    true,
			Message:   "Lorem ipsum",
		}
		r, err := json.Marshal(expected)
		assert.NoError(t, err)

		hdlr := handler.NewGetTodosHandler(&stubRepo{
			todo: expected,
		})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, r, body)
	})

	t.Run("returns 404 on non existing todo call", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		uuid := uuid.New()
		ctx.Request, _ = http.NewRequest("GET", fmt.Sprintf("/todo?id=%s", uuid.String()), http.NoBody)

		hdlr := handler.NewGetTodosHandler(&stubRepo{
			err: repository.ErrTodoNotFound,
		})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Empty(t, body)
	})

	t.Run("returns 500 on todo call with repo error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request, _ = http.NewRequest("GET", fmt.Sprintf("/todo?id=%s", uuid.New().String()), http.NoBody)

		hdlr := handler.NewGetTodosHandler(&stubRepo{
			err: errors.New("test"),
		})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Empty(t, body)
	})

	t.Run("returns 400 on todo call with empty ID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request, _ = http.NewRequest("GET", "/todo?id=", http.NoBody)

		hdlr := handler.NewGetTodosHandler(&stubRepo{})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Empty(t, body)
	})

	t.Run("returns 400 on todo call with non-valid UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request, _ = http.NewRequest("GET", "/todo?id=1234", http.NoBody)

		hdlr := handler.NewGetTodosHandler(&stubRepo{})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Empty(t, body)
	})
}

func TestNewPostTodoHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 201 on successful call", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		payload, err := json.Marshal(map[string]interface{}{
			"is_done": true,
			"message": "Lorem ipsum",
		})
		assert.NoError(t, err)

		ctx.Request, _ = http.NewRequest("POST", "/todo", bytes.NewReader(payload))

		expected := model.Todo{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			IsDone:    true,
			Message:   "Lorem ipsum",
		}
		r, err := json.Marshal(expected)
		assert.NoError(t, err)

		hdlr := handler.NewPostTodoHandler(&stubRepo{
			todo: expected,
		})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, r, body)
	})

	t.Run("returns 400 on wrong payload", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		payload, err := json.Marshal(map[string]interface{}{
			"is_done": "test",
		})
		assert.NoError(t, err)

		ctx.Request, _ = http.NewRequest("POST", "/todo", bytes.NewReader(payload))

		hdlr := handler.NewPostTodoHandler(&stubRepo{})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Empty(t, body)
	})

	t.Run("returns 500 on repo error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		payload, err := json.Marshal(map[string]interface{}{
			"is_done": true,
			"message": "Lorem ipsum",
		})
		assert.NoError(t, err)

		ctx.Request, _ = http.NewRequest("POST", "/todo", bytes.NewReader(payload))

		hdlr := handler.NewPostTodoHandler(&stubRepo{
			err: errors.New("test"),
		})
		hdlr(ctx)

		resp := rr.Result()
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Empty(t, body)
	})
}
