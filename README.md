# Redirix

**Redirix** is a lightweight, authenticated SOCKS5 proxy server that registers itself in Redis for further usage, such as distributed scanning and service discovery.

[![Test Redirix](https://github.com/vulnebify/redirix/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/vulnebify/redirix/actions/workflows/test.yaml)
[![Release Redirix](https://github.com/vulnebify/redirix/actions/workflows/release.yaml/badge.svg)](https://github.com/vulnebify/redirix/actions/workflows/release.yaml)

---

## ✨ Features

- 🔐 SOCKS5 proxy with authentication
- 📡 Auto-registers to Redis with TTL
- 🐳 Docker-ready & GitHub Actions integrated
- ⚡ Built with Go, deploy anywhere

---

## 📦 Installation

### Build locally:

```bash
make build
```

### Or use Docker:

```bash
docker build -t redirix .
```

---

## 🚀 Usage

### Local binary:

```bash
./bin/redirix serve --redis-url=redis://writer:secret@localhost:6380 
```

### Docker:

```bash
docker run --rm --network=host redirix serve --redis-url=redis://writer:secret@localhost:6380
```

---

### Flags:

| Flag               | Description                       | Default         |
|--------------------|-----------------------------------|-----------------|
| `--redis-url`      | Redis connection string           | *required*      |
| `--redis-prefix`   | Redis key prefix                  | `redirix:proxy` |
| `--redis-ttl`      | TTL for Redis registration        | `10s`           |
| `--redis-interval` | Ping interval to Redis            | `5s`            |
| `--proxy-port`     | SOCKS5 proxy port                 | 1080            |
| `--proxy-user`     | Proxy username                    | `redirix`       |
| `--proxy-pass`     | Proxy password                    | *(generated)*   |

---

## 🧪 Testing

Start Redis (with ACL) via Docker Compose:

```bash
docker compose up -d
```

Then run tests:

```bash
go test ./internal/app
```

---

## 📥 GitHub Release

Create a versioned release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The binary will appear under [Releases](../../releases).

---

## 📝 License

This project is licensed under the [MIT License](./LICENSE).

