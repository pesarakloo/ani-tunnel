package backup

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pesarakloo/ani-tunnel/internal/compat"
)

func TestLegacyBackupRestoresItsSettings(t *testing.T) {
	archive := archiveOf(t,
		entry{name: compat.LegacyBackupMeta, body: `{"auto_refresh_hours":6}`},
		entry{name: "saved.toml", body: "saved credentials"},
	)
	stage := t.TempDir()
	got, err := stageRestore(bytes.NewReader(archive), t.TempDir(), stage)
	if err != nil {
		t.Fatal(err)
	}
	if got.AutoRefreshHours != 6 {
		t.Fatalf("legacy refresh interval was lost: %+v", got)
	}
	if _, err := os.Stat(filepath.Join(stage, compat.LegacyBackupMeta)); !os.IsNotExist(err) {
		t.Fatal("sidecar leaked into configuration")
	}
	data, err := os.ReadFile(filepath.Join(stage, "saved.toml"))
	if err != nil || string(data) != "saved credentials" {
		t.Fatalf("saved config: %s, %v", data, err)
	}
}

func TestBackupRetentionOrdersBothNamesByTimestamp(t *testing.T) {
	dir := t.TempDir()
	for day := 1; day <= 12; day++ {
		prefix := compat.LegacyBackupPrefix
		if day > 6 {
			prefix = "ani-tunnel-backup-"
		}
		name := fmt.Sprintf("%s2026-10-%02d-120000.tar.gz", prefix, day)
		if err := os.WriteFile(filepath.Join(dir, name), []byte("backup"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	pruneBackups(dir)
	remaining := backupArchives(dir)
	if len(remaining) != 10 {
		t.Fatal(remaining)
	}
	if filepath.Base(remaining[0]) != "ani-tunnel-backup-2026-10-12-120000.tar.gz" || filepath.Base(remaining[9]) != compat.LegacyBackupPrefix+"2026-10-03-120000.tar.gz" {
		t.Fatal(remaining)
	}
}
