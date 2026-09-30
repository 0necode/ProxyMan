package node

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreAddAndList(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)

	a := &Node{Name: "Beta", Type: "vmess", Server: "b.example.com", Port: 443, UUID: "u-b"}
	b := &Node{Name: "Alpha", Type: "trojan", Server: "a.example.com", Port: 443, Password: "p"}

	if !s.Add(a) {
		t.Fatal("first Add should report a new node")
	}
	if !s.Add(b) {
		t.Fatal("second Add should report a new node")
	}
	if s.Add(b) {
		t.Fatal("re-adding the same node should report existing")
	}

	got := s.List()
	if len(got) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(got))
	}
	if got[0].Name != "Alpha" {
		t.Errorf("expected Alpha first (sorted), got %s", got[0].Name)
	}
}

func TestStorePersistenceRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)
	s.Add(&Node{
		Name: "HK", Type: "vless", Server: "hk.example.com", Port: 443,
		UUID: "uuid-1", SNI: "sni.example.com", TLS: true, Network: "ws",
		WSPath: "/ray", WSHost: "cdn.example.com", Source: "manual",
	})
	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	reopened := NewStore(dir)
	nodes := reopened.List()
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node after reload, got %d", len(nodes))
	}
	n := nodes[0]
	if n.UUID != "uuid-1" || !n.TLS || n.WSPath != "/ray" || n.SNI != "sni.example.com" {
		t.Errorf("node fields did not survive the round trip: %+v", n)
	}
}

func TestStoreRemoveAndDelay(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)
	s.Add(&Node{Name: "N1", Type: "vmess", Server: "x", Port: 1, UUID: "u1"})

	s.RecordDelay("N1", 128)
	if got := s.Get("N1"); got == nil || got.LastDelay != 128 {
		t.Fatalf("expected delay 128, got %+v", got)
	}
	if !s.Remove("N1") {
		t.Fatal("Remove should report success")
	}
	if s.Remove("N1") {
		t.Fatal("removing a missing node should report false")
	}
}

func TestStorePathIsUnderWorkDir(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)
	if got := s.Path(); filepath.Dir(got) != filepath.Join(dir, "nodes") {
		t.Errorf("unexpected store path: %s", got)
	}
	if _, err := os.Stat(filepath.Join(dir, "nodes")); err != nil {
		t.Errorf("node directory not created: %v", err)
	}
}

func TestProbeHandlesBadAddresses(t *testing.T) {
	p := NewProber()
	p.Timeout = 200 * time.Millisecond

	results := p.Probe([]*Node{
		{Name: "no-port", Server: "example.com", Port: 0},
		{Name: "no-host", Server: "", Port: 443},
	})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.OK {
			t.Errorf("%s should be unreachable", r.Name)
		}
		if r.Err == "" {
			t.Errorf("%s should carry an error message", r.Name)
		}
	}
}
