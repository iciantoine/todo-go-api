package option

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

// ConfigureLogging configures Zerolog's log level.
// See: https://pkg.go.dev/github.com/rs/zerolog#pkg-variables
func ConfigureLogging(lvl string) error {
	switch strings.ToLower(lvl) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		return fmt.Errorf("unknown loglevel: %s", lvl)
	}
	return nil
}

// Endpoint is an HTTP endpoint.
type Endpoint struct {
	Addr      string
}

// Postgres is a connection to PostgreSQL.
type Postgres struct {
	User string
	Pass string
	Host string
	Port string
	Name string
	CA   string
	Cert string
	Key  string
	Mode string
}

// DSN returns the Postgres connection string.
func (pg Postgres) DSN() string {
	builder := new(strings.Builder)

	if pg.User != "" {
		builder.WriteString(fmt.Sprintf("user=%s", pg.User))
		builder.WriteByte(' ')
	}
	if pg.Pass != "" {
		builder.WriteString(fmt.Sprintf("password=%s", pg.Pass))
		builder.WriteByte(' ')
	}
	if pg.Host != "" {
		builder.WriteString(fmt.Sprintf("host=%s", pg.Host))
		builder.WriteByte(' ')
	}
	if pg.Port != "" {
		builder.WriteString(fmt.Sprintf("port=%s", pg.Port))
		builder.WriteByte(' ')
	}
	if pg.Name != "" {
		builder.WriteString(fmt.Sprintf("dbname=%s", pg.Name))
		builder.WriteByte(' ')
	}
	if pg.CA != "" {
		builder.WriteString(fmt.Sprintf("sslrootcert=%s", pg.CA))
		builder.WriteByte(' ')
	}
	if pg.Cert != "" {
		builder.WriteString(fmt.Sprintf("sslcert=%s", pg.Cert))
		builder.WriteByte(' ')
	}
	if pg.Key != "" {
		builder.WriteString(fmt.Sprintf("sslkey=%s", pg.Key))
		builder.WriteByte(' ')
	}
	if pg.Mode != "" {
		builder.WriteString(fmt.Sprintf("sslmode=%s", pg.Mode))
	} else {
		builder.WriteString("sslmode=disable")
	}

	return builder.String()
}
