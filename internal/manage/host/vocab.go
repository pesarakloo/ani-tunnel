package host

import (
	"github.com/pesarakloo/ani-tunnel/internal/manage/core"
	"github.com/pesarakloo/ani-tunnel/internal/manage/spec"
)

// The tunnel listing and the address vocabulary, under the names the code in
// this package was written with.
type tunnel = core.Tunnel

var (
	listTunnels    = core.List
	fileExists     = core.FileExists
	addrPort       = spec.AddrPort
	isDatagram     = spec.IsDatagram
	isWildcardBind = spec.IsWildcardBind
)
