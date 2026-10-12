package manage

import (
	"strings"
	"testing"

	"github.com/pesarakloo/ani-tunnel/internal/compat"
)

func TestLegacySetupLinksStillCarryTheSameCredentials(t *testing.T) {
	want := ShareLink{Kind: "reverse", From: "iran", Tok: "existing-server-token", Tr: "tcp", Port: "443", Host: "192.0.2.1"}
	current, err := want.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(current, "ani-tunnel://") {
		t.Fatal(current)
	}
	old := compat.LegacyShareScheme + strings.TrimPrefix(current, shareScheme)
	for _, raw := range []string{old, FindSetupLink("Saved link: " + old)} {
		got, err := DecodeShareLink(raw)
		if err != nil {
			t.Fatal(err)
		}
		if got.Tok != want.Tok || got.Host != want.Host || got.Tr != want.Tr {
			t.Fatalf("old credentials changed: %+v", got)
		}
	}
}
