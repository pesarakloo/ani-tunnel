# Releasing ani tunnel

Keep VERSION and internal/app/app.go in step. The initial fork release is v1.9.1, based on upstream v1.9.0. Run the tests in VALIDATION.md before publishing.

The Release workflow runs when VERSION changes on main, when a v* tag is pushed, or manually. It tests, cross-builds all supported Linux architectures and publishes archives, install.sh, SHA256SUMS, and the build's module inventory. Version changes must use a new tag; never replace a published tag with different code.

Archives are named ani-tunnel_linux_<architecture>.tar.gz. Their executable is named backpack to retain compatibility with existing systemd units and the updater. The installer exposes the ani-tunnel command alias.

## Optional release signatures

The fork cannot sign using the upstream author's key. Its default ReleasePublicKey is empty; updates still require SHA256 checksums fetched from GitHub HTTPS directly or through the encrypted SOCKS tunnel. Third-party release mirrors are not used.

To enable Ed25519 signatures, generate a pair with `make release-key`, put the private key in this repository's RELEASE_SIGNING_KEY Actions secret and put its public half in internal/app/app.go's ReleasePublicKey. Never commit the private key. Once a key is pinned, missing/invalid signatures are refused by the updater. Ensure both halves match before publishing. The standalone installer checks SHA256; it does not claim Ed25519 verification.
