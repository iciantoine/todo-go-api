package server_test

import (
	"testing"

	"github.com/iciantoine/todo-go-api/server"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.NotNil(t, server.WithApplicationAddress("127.0.0.1:8080"))
	assert.NotNil(t, server.WithDatabase("todo", "todo", "127.0.0.1", "5432", "todo", "disable"))
	assert.NotNil(t, server.WithLogLevel("debug"))
}
