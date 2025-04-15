package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/things-go/go-socks5"
	"github.com/vulnebify/redirix/internal/app"
)

var (
	redisKeyPrefix string
	redisURL       string
	proxyPort      int
	redisTTL       time.Duration
	pingInterval   time.Duration
	proxyUser      string
	proxyPassword  string
)

func getRedisURL() string {
	return os.Getenv("TUNNELIER_REDIS_URL")
}

func generatePassword() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return "defaultpass"
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Server to a working WireGuard VPN from MongoDB",
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := app.GetLocalIP()
		hostname, _ := os.Hostname()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		opts, err := app.ParseRedisOptions(redisURL)

		if err != nil {
			return fmt.Errorf("failed to parse redis-url: %w", err)
		}

		proxy_url := fmt.Sprintf("socks5://%s:%s@%s:%d", proxyUser, proxyPassword, ip, proxyPort)

		rdb := redis.NewClient(opts)
		go app.StartRedisHeartbeat(ctx, rdb, hostname, ip, redisKeyPrefix, proxy_url, redisTTL, pingInterval)

		auth := socks5.StaticCredentials{
			proxyUser: proxyPassword,
		}
		creds := &socks5.UserPassAuthenticator{Credentials: auth}
		server := socks5.NewServer(
			socks5.WithAuthMethods([]socks5.Authenticator{creds}),
		)

		fmt.Printf("[Redirix] Starting SOCKS5 proxy on :%d\n", proxyPort)
		fmt.Printf("[Redirix] Proxy Auth → user: %s  pass: %s\n", proxyUser, proxyPassword)
		return server.ListenAndServe("tcp", fmt.Sprintf(":%d", proxyPort))
	},
}

func init() {
	// Redis configuration
	serveCmd.Flags().StringVar(&redisURL, "redis-url", getRedisURL(), "Redis connection URL")
	serveCmd.Flags().StringVar(&redisKeyPrefix, "redis-prefix", "redirix:proxy", "Redis key prefix")
	serveCmd.Flags().DurationVar(&redisTTL, "redis-ttl", 10*time.Second, "TTL for Redis key")
	serveCmd.Flags().DurationVar(&pingInterval, "redis-interval", 5*time.Second, "Interval for Redis heartbeat")

	// Proxy configuration
	serveCmd.Flags().IntVar(&proxyPort, "proxy-port", 1080, "SOCKS5 proxy port")
	serveCmd.Flags().StringVar(&proxyUser, "proxy-user", "redirix", "SOCKS5 proxy username")
	serveCmd.Flags().StringVar(&proxyPassword, "proxy-pass", generatePassword(), "SOCKS5 proxy password (default: generated)")

	_ = serveCmd.MarkFlagRequired("redis-url")
}
