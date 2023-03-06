//go:build integration

package server_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iciantoine/todo-go-api/server"
	"github.com/stretchr/testify/assert"
)

const serverTimeoutSeconds = 10

func TestListen(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appAddr := addr()
	rootURL := fmt.Sprintf("http://%s", appAddr)

	t.Run("error on wrong options", func(t *testing.T) {
		assert.Error(t, server.Listen(ctx, server.WithLogLevel("test")))
	})

	t.Run("error when database connection fails", func(t *testing.T) {
		assert.Error(t, server.Listen(ctx, server.WithDatabase("todo", "todo", "test.test", "5432", "todo", "disable")))
	})

	go func() {
		assert.NoError(t, server.Listen(ctx,
			server.WithApplicationAddress(appAddr),
			server.WithDatabase("todo", "todo", "127.0.0.1", "5432", "todo", "disable"),
			server.WithLogLevel("debug"),
		))
	}()

	assert.NoError(t, waitForServer(rootURL, serverTimeoutSeconds))

	t.Run("501 response on non configured endpoints", func(t *testing.T) {
		req, _ := http.NewRequest("GET", rootURL, http.NoBody)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
	})
}

// Addr returns a random, free TCP address.
func addr() string {
	lst, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer lst.Close()
	addr, ok := lst.Addr().(*net.TCPAddr)
	if !ok {
		return ""
	}
	return addr.String()
}

func waitForServer(URL string, timeoutSeconds int) error {
	var resp *http.Response
	var err error
	for i := 0; i < timeoutSeconds; i++ {
		time.Sleep(1 * time.Second)

		resp, err = http.DefaultClient.Get(URL)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("server was not reachable in a long time, last-time response and error: %v; %v", resp, err)
}
