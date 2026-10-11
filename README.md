<p align="center"><img src="img/ani-banner.svg" alt="ani tunnel" width="100%"></p>

# ani tunnel

Tunnel management for Iran and abroad, with a web panel, terminal menu, Telegram alerts and automatic connection recovery.

[فارسی](README_FA.md) · [Releases](https://github.com/pesarakloo/ani-tunnel/releases) · [Telegram: @xray_chrome](https://t.me/xray_chrome) · [Documentation](docs/README.md)

## Install or upgrade

Run as root on **both servers**:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/pesarakloo/ani-tunnel/main/install.sh)
```

Open the menu later:

```bash
sudo ani-tunnel
```

The installer downloads the matching Linux release and checks its SHA256. It stages the executable, saves the previous binary and configuration under `/root/BackPack/backups/ani-*`, and restarts only services that were already running. A failed service restart restores the previous executable. A successful systemd start is not an end-to-end network health check; check tunnel traffic after upgrading.

For compatibility with existing installations, configuration stays in `/etc/backpack`, systemd units keep the `backpack-` prefix, and `/usr/local/bin/backpack` remains the executable path. `ani-tunnel` is its command alias. Existing tokens, port mappings, setup links and credentials are retained. Wire-level identifiers have not been renamed.

### Build from source

This also works before the first release is available. Go 1.26.9 or newer is required; the installer can install the required toolchain.

```bash
git clone https://github.com/pesarakloo/ani-tunnel.git
cd ani-tunnel
sudo BP_BUILD_FROM_SOURCE=1 bash install.sh
```

### Offline installation

Download `install.sh`, `SHA256SUMS` and `ani-tunnel_linux_amd64.tar.gz` (or the matching ARM archive) from the **same release**. Copy them into a private directory such as `/root/ani-install` and run:

```bash
cd /root/ani-install
sha256sum -c SHA256SUMS --ignore-missing
bash install.sh
```

This fork initially uses mandatory SHA256 verification over GitHub HTTPS for online updates. It does not ship upstream's release signing key. Maintainers can configure their own Ed25519 key as documented in [RELEASING.md](RELEASING.md).

## What changed in v1.9.1

- Fixed TCP and WebSocket data connections outliving a disconnected control client, retaining open sockets and consuming the connection cap indefinitely.
- Bind data sockets to their owner at pool admission, so a reconnect cannot accidentally adopt a previous client's connections.
- Detect unresponsive TCP and WS/WSS control peers with periodic RTT probes, including paths where writes still appear successful. Healthy idle clients remain connected; the watchdog allows at least 30 seconds (or three heartbeat intervals) without an RTT reply.
- Preserve upstream's reconnect, write timeout and listener lifecycle fixes.
- Blue ani tunnel branding across the terminal, web panel, login, two-factor screen, app icons and Telegram messages.
- Installer, updater, source and support links use this repository and `@xray_chrome`.

The connection leak is reproduced by regression tests. Without the affected server's logs and configuration, it is not possible to identify every cause of its reported hours-long failure. See [VALIDATION.md](VALIDATION.md) for the actual checks and limits.

## Use

For a reverse tunnel, configure the **Iran machine as server** and the **abroad machine as client**, with the same tunnel token and compatible transport settings. The client dials the server; forwarded ports on Iran reach services abroad. The menu also supports direct and layer-3 tunnels.

- [Choose a transport](docs/choosing-a-transport.md)
- [Web panel](docs/web-panel.md)
- [Telegram monitoring](docs/telegram-bot.md)
- [Configuration reference](docs/config-reference.md)
- [Troubleshooting](docs/troubleshooting.md)

```bash
ani-tunnel version
ani-tunnel tunnel list
ani-tunnel check -c /etc/backpack/your-tunnel.toml
journalctl -u backpack-your-tunnel.service -n 150 --no-pager
```

## Development

```bash
go test ./... -timeout 20m
go test -race ./internal/server/transport ./internal/client/transport ./internal/utils/handlers -timeout 10m
python3 tools/validation/install_test.py
make release
```

The documentation inherited from upstream retains its original verification stamps. Fork-specific changes and installation behavior are described here, in the Persian README and in VALIDATION.md.

## License and origin

AGPL-3.0. The original license, notices and trademark policy are retained.

Based on BackPack by Amin Mohammadi (AminMGMT)
https://github.com/AminMGMT/BackPack

Base: upstream v1.9.0, commit `132467bfdfb2698d74ada2992413e75bfccbfbd3`.
This is an independently branded fork; ani tunnel modifications are maintained in `pesarakloo/ani-tunnel`.
