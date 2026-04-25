# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected listeners.

---

## Installation

```bash
go install github.com/yourname/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a baseline of allowed ports:

```bash
portwatch --allow 22,80,443 --interval 10s
```

Run a one-time scan and print all open ports:

```bash
portwatch scan
```

Watch for unexpected listeners and send an alert to a webhook:

```bash
portwatch --allow 22,80,443 --webhook https://hooks.example.com/alert
```

**Example output:**

```
[2024-01-15 09:32:11] INFO  Watching ports | allowed: 22, 80, 443
[2024-01-15 09:32:21] WARN  Unexpected listener detected: 0.0.0.0:8080 (pid 4821, nginx)
[2024-01-15 09:32:21] INFO  Alert sent to webhook
```

---

## Flags

| Flag         | Default | Description                          |
|--------------|---------|--------------------------------------|
| `--allow`    | none    | Comma-separated list of allowed ports |
| `--interval` | `30s`   | How often to scan for open ports     |
| `--webhook`  | none    | Webhook URL to POST alerts to        |
| `--verbose`  | false   | Enable verbose logging               |
| `--config`   | none    | Path to a YAML config file           |

---

## Config File

Instead of passing flags each time, you can provide a YAML config file:

```bash
portwatch --config /etc/portwatch/config.yaml
```

```yaml
allow:
  - 22
  - 80
  - 443
interval: 30s
webhook: https://hooks.example.com/alert
verbose: false
```

Command-line flags take precedence over values defined in the config file.

---

## License

MIT © [yourname](https://github.com/yourname)
