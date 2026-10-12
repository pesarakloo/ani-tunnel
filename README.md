<p align="center"><img src="img/ani-banner.svg" alt="ani tunnel" width="100%"></p>

# ani tunnel

High-performance tunneling between Iran and abroad, with an interactive CLI, a web panel, Telegram monitoring and automatic connection recovery.

[فارسی](README_FA.md) · [Releases](https://github.com/pesarakloo/ani-tunnel/releases) · [Telegram: @xray_chrome](https://t.me/xray_chrome) · [Documentation](docs/README.md)

# Quick Start

## 1. Install

Run as root on both servers:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/pesarakloo/ani-tunnel/main/install.sh)
```

Then:

```bash
sudo ani-tunnel
```

The installer downloads the release for the server architecture, verifies its published checksum, installs it and opens the CLI.

Supported release architectures currently include:

```text
x86_64  → amd64
aarch64 → arm64
```

If the server cannot reach GitHub, ani tunnel also supports a completely offline installation.

See [Installation](docs/install.md).

---

## 2. Configure the Iran server

```bash
sudo ani-tunnel
```

Then:

```text
1. Setup Iran
→ Reverse
→ Transport
→ Iran IP Or Domain (the detected one is the default)
→ Tunnel Port
→ Forwarded Ports
→ Tunnel Name
→ Security Token (generated — press Enter)
→ Carry UDP As Well As TCP
→ the transport's own questions (certificate for WSS, flags for PCK, …)
→ How Should The Tunnel Be Tuned?
```

The summary before **Create This Tunnel** shows a one-line **Setup Link** (`ani-tunnel://…`). Copy it.

For a first deployment, **TCP** is the simplest starting point.

---

## 3. Configure the kharej server

```bash
sudo ani-tunnel
```

Then:

```text
2. Setup Kharej
→ Reverse
→ Same transport
→ Setup Link → paste the link → Tunnel Name → Create This Tunnel
```

**Manual** is there too: Iran IP, the same tunnel port, a name, the same token, the same preset.

---

## 4. Check the tunnel

Use:

```text
Manage → Status
```

Then:

```text
Manage → Health Check
```

If you are unsure which transport to use:

```text
Manage → Link Test
```

Link Test measures the route over TCP, including latency, jitter and packet loss, and recommends a suitable transport.

See [Choosing a transport](docs/choosing-a-transport.md).

---


## Features

- Reverse and direct tunnels with multiple transports.
- Automatic reconnection, control-channel liveness checks and cleanup of disconnected clients.
- Web management, two-factor authentication, traffic statistics and logs.
- Telegram reports and alerts through a tunnel peer.
- Verified updates, configuration backups and rollback.

## Documentation

- [Step-by-step tutorials](tutorial/README.md)
- [Web panel](docs/web-panel.md)
- [Telegram monitoring](docs/telegram-bot.md)
- [Configuration reference](docs/config-reference.md)
- [Updates and rollback](docs/updates.md)
- [Troubleshooting](docs/troubleshooting.md)

## License and origin

AGPL-3.0. See [LICENSE](LICENSE), [NOTICE](NOTICE) and [TRADEMARK.md](TRADEMARK.md).

Based on BackPack by Amin Mohammadi (AminMGMT)
https://github.com/AminMGMT/BackPack

This is an independently branded fork. The installation and quick-start instructions are adapted from the original project's README for ani tunnel.
