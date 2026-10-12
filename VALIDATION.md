# ani tunnel v1.9.1 validation

Source base: AminMGMT/ani-tunnel `132467bfdfb2698d74ada2992413e75bfccbfbd3`
(v1.9.0). Checks below were run on Linux amd64 with Go 1.26.9 while preparing
this fork. They are local results, not a claim of live Iran–abroad validation.

## Reproduced fault

`TestLostTCPClientReleasesActiveDataQuota` failed against the upstream relay
lifetime: after only the control socket was lost, an idle data connection kept
its connection slot. The test reported `lost client leaked an active data
connection and its quota`. The generation stayed alive, so its context could
not clean up the departed client's data sockets.

TCP and WebSocket data connections now carry the owner's context captured at
admission. Losing or replacing that client closes its flows and releases their
quota. The listener generation stays available for reconnection. This ends old
flows; it does not promise to preserve a user's TCP session across client loss.
A WS pool request spanning client replacement is rejected instead of being
adopted by the new client.

Periodic RTT replies detect a peer that has stopped responding even when local
writes succeed. The production timeout is at least 30 seconds, or three
heartbeat intervals if longer. Legacy TCP clients that echo RTT but do not echo
heartbeats are covered by a regression test.

## Executed checks

| Check | Result |
| --- | --- |
| Five new loss, replacement and control-watchdog regression tests | Passed, including ten repetitions with the race detector |
| Server transports, client transports and relay handlers with `-race` | Passed |
| Ordinary end-to-end suite | Passed in about 305 seconds; includes transport traffic, reconnect, control adoption, failure and compatibility tests |
| Sustained TCP traffic, `ANI_TUNNEL_SOAK_SECONDS=120` | Passed: 1,607,921 round trips, zero failed; goroutines averaged 60.4 early / 63.1 late, descriptors 74.3 / 78.8; both resource-growth assertions passed |
| Installer in isolated temporary paths with fake systemd | Five tests passed: upgrade preserves config, previous binary backup, active-service restart with rollback, checksum refusal, fresh-install alias and empty-binary refusal |
| Web panel, menu, application/version tests after branding | Passed; panel JavaScript tests run through the Go harness |
| `go vet ./...` | Passed |
| `bash -n install.sh`, `git diff --check`, Go formatting | Passed |
| Static Linux amd64 and arm64 compilation | Passed; amd64 executable reports v1.9.1 and this repository |

`go test ./... -timeout 12m` was run. It did **not** pass in its entirety:
`internal/client` cannot inspect interface `lo` through netlink, and
`internal/enginectl` cannot bind Unix-domain sockets in this execution
environment (`operation not permitted`). The relevant failures were reproduced
on the unmodified upstream commit as well. The remaining packages passed.
No test was weakened or disabled to hide these restrictions. GitHub Actions
runs the complete suite on its Linux runner before publishing release assets.

## Reproduce

```bash
go test ./... -timeout 20m
go test -race ./internal/server/transport ./internal/client/transport ./internal/utils/handlers -timeout 10m
go test -race ./internal/server/transport -run 'Test(LostTCP|LostWS|ReplacedWS|SilentControl|RTTReplies)' -count=10
ANI_TUNNEL_SOAK_SECONDS=120 go test ./internal/e2e -run '^TestASustainedRunDoesNotLeak$' -count=1 -v -timeout 3m
python3 tools/validation/install_test.py
go vet ./...
```

The installer tests redirect all installation paths to a temporary directory,
bypass only its root entry guard, and stub downloads and systemd. They do not
modify host configuration or restart real services.

## Remaining field checks

No affected server logs, configuration or SSH access were provided. This
reproduced leak is a plausible contributor to the reported hours-long failure;
it does not establish the sole cause on those servers. The two-minute soak is
not a multi-hour reliability guarantee. Real TUN devices, packet filtering,
live Iran–abroad routing and production systemd recovery need host validation.
Browser screenshot verification was unavailable locally; the panel's automated
tests and embedded branding assets were checked.

Upgrade both peers, verify user traffic after service restart and observe the
tunnel over the previously failing interval. If a failure remains, preserve
the relevant `journalctl -u ani-tunnel-<name>.service` output and redact tokens,
passwords and private keys before sharing it.
