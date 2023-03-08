package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iciantoine/todo-go-api/model"
	"github.com/iciantoine/todo-go-api/repository"
	"github.com/rs/zerolog/log"
)

type TodoRepo interface {
	GetTodos(ctx context.Context) ([]model.Todo, error)
	GetTodo(ctx context.Context, id uuid.UUID) (model.Todo, error)
	AddTodo(ctx context.Context, model model.Todo) (model.Todo, error)
}

func NewGetTodosHandler(repo TodoRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, exists := ctx.GetQuery("id")

		// ID is given, trying to get the specified todo
		if exists {
			getTodo(ctx, repo, id)
			return
		}

		res, err := repo.GetTodos(ctx)
		if err == nil {
			ctx.JSON(http.StatusOK, res)
			return
		}

		log.Ctx(ctx.Request.Context()).Error().Err(err).Msg("error while getting todos")
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func NewPostTodoHandler(repo TodoRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.Todo
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Ctx(ctx.Request.Context()).Warn().Err(err).Msg("could not bind request body")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		res, err := repo.AddTodo(ctx, req)
		if err != nil {
			log.Ctx(ctx.Request.Context()).Error().Err(err).Msg("error while adding todo")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusCreated, res)
	}
}

func getTodo(ctx *gin.Context, repo TodoRepo, id string) {
	if id == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Ctx(ctx.Request.Context()).Warn().Err(err).Msg("could not parse uuid")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := repo.GetTodo(ctx, uuid)

	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, res)
	case errors.Is(err, repository.ErrTodoNotFound):
		ctx.AbortWithStatus(http.StatusNotFound)
	default:
		log.Ctx(ctx.Request.Context()).Error().Str("id", id).Err(err).Msg("error while getting todo")
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
