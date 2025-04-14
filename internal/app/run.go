package app

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/things-go/go-socks5"
)

var Version = "dev"

var (
	redisKeyPrefix string
	redisURL       string
	proxyPort      int
	redisTTL       time.Duration
	pingInterval   time.Duration
	proxyUser      string
	proxyPassword  string
)

func Run(ctx context.Context, args []string) error {
	for _, arg := range args {
		if arg == "--version" || arg == "-v" {
			fmt.Println("Redirix version:", Version)
			return nil
		}
	}

	fs := flag.NewFlagSet("redirix", flag.ContinueOnError)

	// Redis configuration
	fs.StringVar(&redisKeyPrefix, "redis-prefix", "redirix:proxy", "Redis key prefix")
	fs.StringVar(&redisURL, "redis-url", "rediss://username:password@localhost:6379", "Redis connection URL")
	fs.DurationVar(&redisTTL, "redis-ttl", 10*time.Second, "TTL for Redis key")
	fs.DurationVar(&pingInterval, "redis-interval", 5*time.Second, "Interval for Redis heartbeat")

	// Proxy configuration
	fs.IntVar(&proxyPort, "proxy-port", 1080, "SOCKS5 proxy port")
	fs.StringVar(&proxyUser, "proxy-user", "redirix", "SOCKS5 proxy username")
	fs.StringVar(&proxyPassword, "proxy-pass", generatePassword(), "SOCKS5 proxy password (if empty, generated at runtime)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	ip := getLocalIP()
	hostname, _ := os.Hostname()

	opts, err := parseRedisOptions(redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse redis-url: %w", err)
	}

	rdb := redis.NewClient(opts)
	go startRedisHeartbeat(ctx, rdb, hostname, ip)

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
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func parseRedisOptions(rawURL string) (*redis.Options, error) {
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

func startRedisHeartbeat(ctx context.Context, rdb *redis.Client, hostname, ip string) {
	key := fmt.Sprintf("%s:%s", redisKeyPrefix, hostname)
	value := fmt.Sprintf("socks5://%s:%s@%s:%d", proxyUser, proxyPassword, ip, proxyPort)

	for {
		err := rdb.Set(ctx, key, value, redisTTL).Err()
		if err != nil {
			fmt.Println("[Redis] Error:", err)
		}
		time.Sleep(pingInterval)
	}
}

func generatePassword() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return "defaultpass"
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}
