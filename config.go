package redops

import (
	"time"

	"github.com/imdario/mergo"
)

const (
	// DefaultRedisNetwork the redis network option, "tcp"
	DefaultRedisNetwork = "tcp"
	// DefaultRedisAddr the redis address option, "127.0.0.1:6379"
	DefaultRedisAddr = "127.0.0.1:6379"
	// DefaultRedisIdleTimeout the redis idle timeout option, time.Duration(5) * time.Minute
	DefaultRedisIdleTimeout = time.Duration(5) * time.Minute
)

// Config the redis configuration
type Config struct {
	Network     string
	Addr        string
	Password    string
	Database    int
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

// DefaultConfig returns the default redis configuration
func DefaultConfig() Config {
	return Config{
		Network:     DefaultRedisNetwork,
		Addr:        DefaultRedisAddr,
		Password:    "",
		Database:    0,
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: DefaultRedisIdleTimeout,
	}
}

// Merge merges the default with the given config and returns the result
func (c Config) Merge(cfg Config) (config Config) {
	config = cfg
	mergo.Merge(&config, c)
	return
}
