// Package compat preserves protocol identifiers needed by already deployed peers.
// These are wire-format constants, not product names. Changing them breaks authentication.
package compat

const (
	SocksPortLabel   = "backpack-socks-v1:"
	PCKLabel         = "backpack-pck-v1:"
	QUICBindingLabel = "EXPORTER-backpack-quic-binding-v1"
	WSSBindingLabel  = "EXPORTER-backpack-wss-binding-v1"
	ICMPLabel        = "backpack-xdi-v1:"
	NoisePrologue    = "backpack-noise-v1"
	NoisePSKLabel    = "backpack-stealth-psk-v1"
	SpoofLabel       = "backpack-spoof-v1:"
	QUICALPN         = "backpack-quic"
	KCPLabel         = "backpack-kcp-v1"
	L3Prologue       = "backpack-l3-v1"
	L3PSKLabel       = "backpack-l3-psk-v1"
	L3ResetLabel     = "backpack l3 quic stateless reset v1\x00"
	DirectLabel      = "backpack-direct-v1|"
	SpoofTestLabel   = "backpack-spooftest-v1:"
	DecoyLabel       = "backpack-decoy\x00"
	ConnTestLabel    = "backpack-conntest/"
	ConnTestHost     = "conntest.backpack.invalid"
	ReleaseDomain    = "backpack release "
	RelayUsername    = "backpack"
	DefaultToken     = "backpack"
)
