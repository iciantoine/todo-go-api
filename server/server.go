package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iciantoine/todo-go-api/database"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/rs/zerolog/log"
)

// Listen starts the HTTP server.
func Listen(parent context.Context, opts ...Option) error {
	cfg := new(config)

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			log.Error().Err(err).Msg("could not configure application")
			return err
		}
	}

	conn, err := database.Connect("pgx", cfg.Database.DSN())
	if err != nil {
		log.Error().Err(err).Msg("could not connect to Lydia database")
		return err
	}
	defer conn.Close()

	return router().Run(cfg.Application.Addr)
}

func router() *gin.Engine {
	router := gin.Default()

	// default handler for unknown routes
	router.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	})

	return router
}
