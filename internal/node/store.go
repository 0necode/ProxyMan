package node

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/Geek0ne/ProxyMan/internal/parser"
)

// Node is a persisted proxy node entry
type Node struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Server    string            `json:"server"`
	Port      int               `json:"port"`
	UUID      string            `json:"uuid,omitempty"`
	Password  string            `json:"password,omitempty"`
	Cipher    string            `json:"cipher,omitempty"`
	SNI       string            `json:"sni,omitempty"`
	Network   string            `json:"network,omitempty"`
	WSPath    string            `json:"ws-path,omitempty"`
	WSHost    string            `json:"ws-host,omitempty"`
	TLS       bool              `json:"tls"`
	Flow      string            `json:"flow,omitempty"`
	Alpn      []string          `json:"alpn,omitempty"`
	Params    map[string]string `json:"params,omitempty"`
	Source    string            `json:"source,omitempty"`
	AddedAt   time.Time         `json:"added-at"`
	LastDelay int               `json:"last-delay-ms,omitempty"`
}

// Store persists imported nodes under the ProxyMan work directory
type Store struct {
	mu    sync.RWMutex
	path  string
	nodes map[string]*Node
}

// NewStore opens (or creates) the node store in the given work directory
func NewStore(workDir string) *Store {
	dir := filepath.Join(workDir, "nodes")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		dir = workDir
	}
	s := &Store{
		path:  filepath.Join(dir, "nodes.json"),
		nodes: make(map[string]*Node),
	}
	s.load()
	return s
}

func (s *Store) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	var list []*Node
	if err := json.Unmarshal(data, &list); err != nil {
		return
	}
	for _, n := range list {
		if n == nil || n.Name == "" {
			continue
		}
		s.nodes[key(n)] = n
	}
}

func (s *Store) saveLocked() error {
	list := s.listLocked()
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// key builds a stable identity for a node so re-imports update in place
func key(n *Node) string {
	if n.UUID != "" {
		return fmt.Sprintf("%s|%s|%s|%d", n.Type, n.UUID, n.Server, n.Port)
	}
	return fmt.Sprintf("%s|%s|%s|%d|%s", n.Type, n.Name, n.Server, n.Port, n.Password)
}

// FromProxyNode converts a parsed node into a storable node
func FromProxyNode(p *parser.ProxyNode, source string) *Node {
	return &Node{
		Name:     p.Name,
		Type:     p.Type,
		Server:   p.Server,
		Port:     p.Port,
		UUID:     p.UUID,
		Password: p.Password,
		Cipher:   p.Cipher,
		SNI:      p.SNI,
		Network:  p.Network,
		WSPath:   p.WSPATH,
		WSHost:   p.WSHost,
		TLS:      p.TLS,
		Flow:     p.Flow,
		Alpn:     p.Alpn,
		Params:   p.Params,
		Source:   source,
		AddedAt:  time.Now(),
	}
}

// Add stores a node and reports whether it replaced an existing entry
func (s *Store) Add(n *Node) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(n)
	_, existed := s.nodes[k]
	if existed {
		n.AddedAt = s.nodes[k].AddedAt
	}
	s.nodes[k] = n
	return !existed
}

// List returns all nodes sorted by name
func (s *Store) List() []*Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.listLocked()
}

// listLocked assumes the caller already holds at least a read lock.
func (s *Store) listLocked() []*Node {
	out := make([]*Node, 0, len(s.nodes))
	for _, n := range s.nodes {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Server < out[j].Server
	})
	return out
}

// Get returns a node by exact name
func (s *Store) Get(name string) *Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, n := range s.nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

// Remove deletes a node by name
func (s *Store) Remove(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, n := range s.nodes {
		if n.Name == name {
			delete(s.nodes, k)
			return true
		}
	}
	return false
}

// RecordDelay stores a measured latency for a node
func (s *Store) RecordDelay(name string, ms int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, n := range s.nodes {
		if n.Name == name {
			n.LastDelay = ms
			return
		}
	}
}

// Save persists the store to disk
func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked()
}

// Path returns the on-disk location of the node store
func (s *Store) Path() string { return s.path }
