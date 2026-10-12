package compat

import (
	"os"
	"path/filepath"
	"strings"
)

// Legacy names are accepted only for upgrades, existing data and older peers.
// Fresh installations use the ani-tunnel namespace throughout.
const (
	LegacyBinary        = "backpack"
	LegacyBinPath       = "/usr/local/bin/backpack"
	LegacyConfigDir     = "/etc/backpack"
	LegacyInstallDir    = "/root/BackPack"
	LegacyServicePrefix = "backpack-"
	LegacyHelperDir     = "/usr/local/lib/backpack"
	LegacyRuntimeDir    = "/run/backpack"
	LegacyShareScheme   = "backpack://"
	LegacyBackupPrefix  = "backpack-backup-"
	LegacyBackupMeta    = ".backpack-backup.json"
)

// ExistingPath uses the new name unless only the legacy path exists. It does
// not move directories: running units, backups and explicit config paths stay valid.
func ExistingPath(preferred, legacy string) string {
	if _, err := os.Stat(preferred); os.IsNotExist(err) {
		if _, err := os.Stat(legacy); err == nil {
			return legacy
		}
	}
	return preferred
}

// BinaryPath retains the real target of the command alias installed by v1.9.1.
// Updating the alias itself would leave existing systemd units on the old binary.
func BinaryPath(preferred, legacy string) string {
	if target, err := os.Readlink(preferred); err == nil {
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(preferred), target)
		}
		if filepath.Clean(target) == filepath.Clean(legacy) {
			return legacy
		}
	}
	return ExistingPath(preferred, legacy)
}

func NormalizeLink(s string) string {
	return strings.ReplaceAll(s, LegacyShareScheme, "ani-tunnel://")
}
