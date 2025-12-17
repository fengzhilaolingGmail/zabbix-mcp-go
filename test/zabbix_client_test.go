package zabbix_test

import (
	"testing"
	"time"
	zabbix "zabbixMcp/zabbix"
)

func TestNewZabbixClient_DefaultTimeout(t *testing.T) {
	c := zabbix.NewZabbixClient("http://example", "u", "p", 0)
	if c.HTTPClient == nil {
		t.Fatal("HTTPClient should not be nil")
	}
	if c.HTTPClient.Timeout != 120*time.Second {
		t.Fatalf("expected default timeout 120s, got %v", c.HTTPClient.Timeout)
	}
}

func TestNewZabbixClient_CustomTimeout(t *testing.T) {
	c := zabbix.NewZabbixClient("http://example", "u", "p", 5)
	if c.HTTPClient.Timeout.Seconds() != 5 {
		t.Fatalf("expected timeout 5s, got %v", c.HTTPClient.Timeout)
	}
}

func TestSetServerTimezone(t *testing.T) {
	c := zabbix.NewZabbixClient("u", "u", "p", 0)
	c.SetServerTimezone("Asia/Shanghai")
	if c.ServerTZ != "Asia/Shanghai" {
		t.Fatalf("expected Asia/Shanghai, got %q", c.ServerTZ)
	}

	// empty string should set to non-empty local string
	c.SetServerTimezone("")
	if c.ServerTZ == "" {
		t.Fatalf("expected non-empty ServerTZ when setting empty, got empty")
	}
}

func TestCachedVersion_CopyBehavior(t *testing.T) {
	c := zabbix.NewZabbixClient("u", "u", "p", 0)
	v := &zabbix.VersionInfo{Major: 6, Minor: 0, Patch: 5, Full: "6.0.5"}
	c.SetCachedVersion(v)

	got := c.GetCachedVersion()
	if got == nil {
		t.Fatal("expected cached version, got nil")
	}
	// modify returned copy
	got.Major = 1

	// ensure original cached version not affected
	stored := c.GetCachedVersion()
	if stored.Major != 6 {
		t.Fatalf("expected stored Major 6, got %d", stored.Major)
	}

	// clear cache
	c.ClearCachedVersion()
	if c.GetCachedVersion() != nil {
		t.Fatalf("expected nil after ClearCachedVersion")
	}
}

func TestRPCError_ErrorString(t *testing.T) {
	e := &zabbix.RPCError{Code: 123, Message: "bad", Data: "details"}
	s := e.Error()
	if s == "" {
		t.Fatalf("expected non-empty error string")
	}
	if !(contains(s, "123") && contains(s, "bad") && contains(s, "details")) {
		t.Fatalf("error string missing expected parts: %s", s)
	}
}

func TestParseVersion(t *testing.T) {
	vd := &zabbix.VersionDetector{}

	v, err := vd.ParseVersion("v6.0.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 6 || v.Minor != 0 || v.Patch != 5 || v.Full != "6.0.5" {
		t.Fatalf("unexpected parse result: %+v", v)
	}

	v2, err := vd.ParseVersion("5.4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v2.Major != 5 || v2.Minor != 4 || v2.Patch != 0 || v2.Full != "5.4" {
		t.Fatalf("unexpected parse result: %+v", v2)
	}
}

// small helper to avoid importing strings in tests repeatedly
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(s) > len(sub) && (indexOf(s, sub) >= 0)))
}

func indexOf(s, sub string) int {
	// simple fallback to avoid extra imports; use strings.Index normally
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
