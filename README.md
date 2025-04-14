# Redirix

**Redirix** is a lightweight, authenticated SOCKS5 proxy that registers itself in Redis for further usage, such as distributed scanning and service discovery.

---

## ✨ Features

- 🔐 SOCKS5 proxy with optional authentication
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
./bin/redirix \
  --redis-url=redis://writer:secret@localhost:6380 \
  --prefix=redirix:proxy \
  --auth-user=red \
  --auth-pass=securepass \
  --port=1080
```

### Docker:

```bash
docker run --rm --network=host redirix \
  --redis-url=redis://writer:secret@localhost:6380 \
  --prefix=redirix:proxy \
  --auth-user=red \
  --auth-pass=securepass \
  --port=1080
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
git tag v0.0.1
git push origin v0.0.1
```

The binary will appear under [Releases](../../releases).

---

## 📝 License

This project is licensed under the [MIT License](./LICENSE).

