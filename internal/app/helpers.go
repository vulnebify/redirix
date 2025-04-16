package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func ParseRedisOptions(rawURL string) (*redis.Options, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	addr := u.Host
	password := ""
	username := ""
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}

	useTLS := u.Scheme == "rediss"

	return &redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       0,
		TLSConfig: func() *tls.Config {
			if useTLS {
				return &tls.Config{}
			}
			return nil
		}(),
	}, nil
}

func StartRedisHeartbeat(ctx context.Context, rdb *redis.Client, hostname, value, prefix string, ttl, interval time.Duration) {
	key := fmt.Sprintf("%s:%s", prefix, hostname)

	for {
		err := rdb.Set(ctx, key, value, ttl).Err()
		if err != nil {
			fmt.Println("[Redis] Error:", err)
		}
		time.Sleep(interval)
	}
}
