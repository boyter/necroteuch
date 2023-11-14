package common

import (
	"github.com/pelletier/go-toml"
)

type Config struct {
	LogLevel   string
	HTTPPort   int
	SqliteName string
}

const (
	DefaultLogLevel   = "info"
	DefaultPort       = 8080
	DefaultSqliteName = ":memory:"
)

func NewConfig() Config {
	config, err := toml.LoadFile("config.toml")

	if err != nil {
		return Config{
			LogLevel:   DefaultLogLevel,
			HTTPPort:   DefaultPort,
			SqliteName: DefaultSqliteName,
		}
	}
	return Config{
		LogLevel:   getOrDefaultString(config, "log.level", DefaultLogLevel),
		HTTPPort:   getOrDefaultInt(config, "server.http_port", DefaultPort),
		SqliteName: getOrDefaultString(config, "data.sqlite_name", DefaultSqliteName),
	}
}

func getOrDefaultString(config *toml.Tree, value, def string) string {
	switch config.Get(value).(type) {
	case string:
		return config.Get(value).(string)
	}

	return def
}

func getOrDefaultInt(config *toml.Tree, value string, def int) int {
	switch config.Get(value).(type) {
	case int:
		return config.Get(value).(int)
	case int64:
		return int(config.Get(value).(int64))
	}

	return def
}
