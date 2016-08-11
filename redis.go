package redops

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// A Redis model
type Redis struct {
	pool *redis.Pool
}

func dial(network string, addr string, pass string) (redis.Conn, error) {
	c, err := redis.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	if pass != "" {
		if _, err = c.Do("AUTH", pass); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

// New connects to the redis
func New(c Config) *Redis {
	if c.IdleTimeout <= 0 {
		c.IdleTimeout = DefaultRedisIdleTimeout
	}
	if c.Network == "" {
		c.Network = DefaultRedisNetwork
	}
	if c.Addr == "" {
		c.Addr = DefaultRedisAddr
	}

	pool := &redis.Pool{IdleTimeout: DefaultRedisIdleTimeout, MaxIdle: c.MaxIdle, MaxActive: c.MaxActive}
	pool.TestOnBorrow = func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	}
	pool.Dial = func() (redis.Conn, error) {
		conn, err := dial(c.Network, c.Addr, c.Password)
		if err != nil {
			return nil, err
		}
		if _, err = conn.Do("SELECT", c.Database); err != nil {
			conn.Close()
			return nil, err
		}
		return conn, err
	}
	return &Redis{pool: pool}
}

//Ping returns true if pong received otherwise false
func (r *Redis) Ping() (bool, error) {
	c := r.pool.Get()
	defer c.Close()
	s, err := redis.String(c.Do("PING"))
	return (s == "PONG"), err
}

// Get returns value, err by its key
func (r *Redis) Get(key string) (interface{}, error) {
	c := r.pool.Get()
	defer c.Close()
	return c.Do("GET", key)
}

// GetString returns value, err by its key
func (r *Redis) GetString(key string) (string, error) {
	return redis.String(r.Get(key))
}

// HmSet sets a map[string]string by key
func (r *Redis) HmSet(key string, val map[string]string) (interface{}, error) {
	c := r.pool.Get()
	defer c.Close()
	vals := []interface{}{key}
	for k, v := range val {
		vals = append(vals, k, v)
	}
	return c.Do("HMSET", vals...)
}

// HGetAll returns a reply by key
func (r *Redis) HGetAll(key string) (interface{}, error) {
	c := r.pool.Get()
	defer c.Close()
	return c.Do("HGETALL", key)
}

// HGetAllMap returns a map[string]string by key
func (r *Redis) HGetAllMap(key string) (map[string]string, error) {
	return redis.StringMap(r.HGetAll(key))
}

// Del deletes by key
func (r *Redis) Del(key string) (interface{}, error) {
	c := r.pool.Get()
	defer c.Close()
	return c.Do("DEL", key)
}
