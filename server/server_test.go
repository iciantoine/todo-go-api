//go:build integration

package server_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iciantoine/todo-go-api/server"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
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

	t.Run("200 response on getting todos", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/todo", rootURL), http.NoBody)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, int64(2), gjson.GetBytes(body, "#").Int())
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, validateSchema(t, "../schema/todos.json", body))
	})

	t.Run("200 response on getting todo", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/todo?id=%s", rootURL, "038863e4-2fbe-4bc3-9e38-1e62e93659f5"), http.NoBody)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, validateSchema(t, "../schema/todo.json", body))
	})

	t.Run("400 response on getting todo without id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/todo?id=", rootURL), http.NoBody)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("400 response on getting todo without a valid uuid as id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/todo?id=%s", rootURL, "038863e4-2fbe-4bc3-9e38"), http.NoBody)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("404 response on getting non existing todo", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/todo?id=%s", rootURL, "038863e4-2fbe-4bc3-9e38-1e62e93659f6"), http.NoBody)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// addr returns a random, free TCP address.
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

// waitForServer sleeps for a max period until URL is responding
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

// validateSchema is a validation function to ensure the response matches a given JSON schema.
func validateSchema(t *testing.T, path string, data []byte) bool {
	schema, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Log(err)
		return false
	}

	sl := gojsonschema.NewBytesLoader(schema)
	dl := gojsonschema.NewBytesLoader(data)

	result, err := gojsonschema.Validate(sl, dl)
	if err != nil {
		t.Log(err)
		return false
	}

	t.Log(result.Errors())
	return result.Valid()
}
