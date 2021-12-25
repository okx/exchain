package redis_cgi

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

const (
	MaxIdle     = 3
	IdleTimeout = 240
)

// Use redis URI scheme. URLs should follow the draft IANA specification for the
// scheme (https://www.iana.org/assignments/uri-schemes/prov/redis).
// eg. redis://user:password@localhost:6379, redis://localhost:16379
// return: address, password, error
func ParseRedisURL(redisURL, requirePass string) (string, string, error) {
	u, err := url.Parse(redisURL)
	if err != nil {
		return "", "", err
	}
	if u.Scheme != "redis" && u.Scheme != "rediss" {
		return "", "", fmt.Errorf("invalid redis URL scheme: %s", u.Scheme)
	}

	// As per the IANA draft spec, the host defaults to localhost and
	// the port defaults to 6379.
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		// assume port is missing
		host = u.Host
		port = "6379"
	}
	if host == "" {
		host = "localhost"
	}
	address := net.JoinHostPort(host, port)
	password, _ := u.User.Password()
	if password == "" {
		password = requirePass
	}
	return address, password, nil
}

func NewPool(redisURL string, redisPass string, logger log.Logger) (*redis.Pool, error) {
	address, password, err := ParseRedisURL(redisURL, redisPass)
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("parsed redis url: %s, address: %s, password: %s",
		redisURL, address, password))

	// new redis pool
	pool := &redis.Pool{
		MaxIdle:     MaxIdle,
		IdleTimeout: IdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				logger.Error(fmt.Sprintf("connect to redis failed: %v", err))
				return nil, err
			}

			if password != "" {
				// if password is set, do auth
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					logger.Error(fmt.Sprintf("redis auth failed: %v", err))
					return nil, err
				}
			}
			return c, nil
		},
	}

	// Test pool connection
	conn := pool.Get()
	defer conn.Close()
	if conn.Err() != nil {
		return nil, fmt.Errorf("unable to connect to redis: %v", redisURL)
	}
	return pool, nil
}
