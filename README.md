# procwatch

Lightweight process supervisor with structured JSON logging and health-check hooks.

---

## Installation

```bash
go install github.com/yourname/procwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/procwatch.git && cd procwatch && go build ./...
```

---

## Usage

Define your processes in a `procwatch.yaml` config file:

```yaml
processes:
  - name: web-server
    command: "./bin/server"
    args: ["--port", "8080"]
    healthcheck:
      url: "http://localhost:8080/health"
      interval: 10s
      timeout: 3s
    restart: on-failure

  - name: worker
    command: "./bin/worker"
    restart: always
```

Then run:

```bash
procwatch start -c procwatch.yaml
```

All output is emitted as structured JSON logs:

```json
{"time":"2024-01-15T10:23:01Z","level":"info","process":"web-server","event":"started","pid":12345}
{"time":"2024-01-15T10:23:11Z","level":"info","process":"web-server","event":"healthcheck_ok"}
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `procwatch.yaml` | Path to config file |
| `-l, --log-level` | `info` | Log level (debug, info, warn, error) |
| `--no-restart` | `false` | Disable automatic restarts |

---

## License

MIT © yourname