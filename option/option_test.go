package option_test

import (
	"errors"
	"testing"

	"github.com/iciantoine/todo-go-api/option"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestConfigureLogging(t *testing.T) {
	assert.NoError(t, option.ConfigureLogging("warn"))
	assert.Equal(t, zerolog.WarnLevel, zerolog.GlobalLevel())
	assert.Error(t, errors.New("unknown loglevel: test"), option.ConfigureLogging("test"))
}

func TestDSN(t *testing.T) {
	pg := option.Postgres{
		User: "user",
		Pass: "pass",
		Host: "host",
		Port: "port",
		Name: "name",
		CA:   "ca",
		Cert: "cert",
		Key:  "key",
		Mode: "mode",
	}

	t.Run("test all fields", func(t *testing.T) {
		dsn := pg.DSN()
		assert.Equal(t, "user=user password=pass host=host port=port dbname=name sslrootcert=ca sslcert=cert sslkey=key sslmode=mode", dsn)
	})

	t.Run("empty sslmode should be disabled sslmode", func(t *testing.T) {
		pg.Mode = ""
		dsn := pg.DSN()
		assert.Equal(t, "user=user password=pass host=host port=port dbname=name sslrootcert=ca sslcert=cert sslkey=key sslmode=disable", dsn)
	})
}
