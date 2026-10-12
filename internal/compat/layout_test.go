package compat

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExistingInstallRemainsReachable(t *testing.T) {
	root := t.TempDir()
	preferred, legacy := filepath.Join(root, "ani-tunnel"), filepath.Join(root, LegacyBinary)
	if got := ExistingPath(preferred, legacy); got != preferred {
		t.Fatal(got)
	}
	if err := os.WriteFile(legacy, []byte("existing"), 0600); err != nil {
		t.Fatal(err)
	}
	if got := BinaryPath(preferred, legacy); got != legacy {
		t.Fatal(got)
	}
	if err := os.Symlink(LegacyBinary, preferred); err != nil {
		t.Fatal(err)
	}
	if got := BinaryPath(preferred, legacy); got != legacy {
		t.Fatalf("alias must update its existing target, got %s", got)
	}
	if err := os.Remove(preferred); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(preferred, []byte("new installation"), 0600); err != nil {
		t.Fatal(err)
	}
	if got := BinaryPath(preferred, legacy); got != preferred {
		t.Fatal(got)
	}
}
