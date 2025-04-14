package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/proxy"
)

func TestRunWritesToRedis(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate unique test prefix
	testID := fmt.Sprintf("test-%d", time.Now().UnixNano())
	prefix := fmt.Sprintf("redirix:test:%s", testID)

	args := []string{
		"-redis-url=redis://writer:secret@localhost:6380",
		"-redis-prefix=" + prefix,
		"-redis-ttl=5s",
		"-redis-interval=2s",
		"-proxy-port=2000",
		"-proxy-user=testuser",
		"-proxy-pass=testpass",
	}

	go func() {
		err := Run(ctx, args) // replace with actual import alias if needed
		if err != nil {
			t.Errorf("Run failed: %v", err)
		}
	}()

	// Connect as scanner (read-only)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Username: "reader",
		Password: "secret",
		DB:       0,
	})
	defer rdb.Close()

	// Wait and check for key
	time.Sleep(3 * time.Second)
	keys, err := rdb.Keys(ctx, prefix+":*").Result()
	assert.NoError(err)
	assert.NotEmpty(keys, "expected redis key not found")

	val, err := rdb.Get(ctx, keys[0]).Result()
	assert.NoError(err)
	assert.Contains(val, "socks5://testuser:testpass@", "unexpected value format")

	fmt.Println("✅ Found key:", keys[0])
	fmt.Println("🔑 Value:", val)
}

func TestRunAndProxyConnectivity(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	testID := fmt.Sprintf("test-%d", time.Now().UnixNano())
	prefix := fmt.Sprintf("redirix:test:%s", testID)
	port := 2099

	args := []string{
		"-redis-url=redis://writer:secret@localhost:6380",
		"-redis-prefix=" + prefix,
		"-redis-ttl=5s",
		"-redis-interval=2s",
		fmt.Sprintf("-proxy-port=%d", port),
		"-proxy-user=testuser",
		"-proxy-pass=testpass",
	}

	go func() {
		err := Run(ctx, args)
		if err != nil {
			t.Errorf("Run failed: %v", err)
		}
	}()

	time.Sleep(2 * time.Second) // wait for proxy to start

	// test connectivity through the proxy
	dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", port), &proxy.Auth{
		User:     "testuser",
		Password: "testpass",
	}, proxy.Direct)
	assert.NoError(err)

	httpTransport := &http.Transport{}
	httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}
	client := &http.Client{Transport: httpTransport, Timeout: 5 * time.Second}

	resp, err := client.Get("http://example.com")
	assert.NoError(err)
	defer resp.Body.Close()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
