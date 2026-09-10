package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteXrayConfig(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteXrayConfig(dir, "warning", "vmess", nil)
	if err != nil {
		t.Fatalf("WriteXrayConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "10809") {
		t.Error("config missing HTTP inbound port 10809")
	}
	if !strings.Contains(content, "10808") {
		t.Error("config missing SOCKS inbound port 10808")
	}
	if !strings.Contains(content, "vmess") {
		t.Error("config missing vmess outbound")
	}
	if !strings.Contains(content, "warning") {
		t.Error("config missing log level")
	}
	t.Logf("Xray config generated at %s (%d bytes)", configPath, len(data))
}

func TestWriteMihomoConfig(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteMihomoConfig(dir, "info", "full")
	if err != nil {
		t.Fatalf("WriteMihomoConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "7890") {
		t.Error("config missing mixed-port 7890")
	}
	if !strings.Contains(content, "Proxy") {
		t.Error("config missing Proxy group")
	}
	if !strings.Contains(content, "url-test") {
		t.Error("config missing url-test type")
	}
	if !strings.Contains(content, "fake-ip") {
		t.Error("config missing fake-ip DNS mode")
	}
	t.Logf("Mihomo config generated at %s (%d bytes)", configPath, len(data))
}

func TestWriteSystemdService(t *testing.T) {
	dir := t.TempDir()
	servicePath, err := WriteSystemdService(dir, "xray", "config.json", "Test Xray Service")
	if err != nil {
		t.Fatalf("WriteSystemdService failed: %v", err)
	}

	data, err := os.ReadFile(servicePath)
	if err != nil {
		t.Fatalf("read service: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Test Xray Service") {
		t.Error("service missing description")
	}
	if !strings.Contains(content, "xray") {
		t.Error("service missing binary name")
	}
	if !strings.Contains(content, "Restart=on-failure") {
		t.Error("service missing restart policy")
	}
	if !strings.Contains(content, "[Install]") {
		t.Error("service missing [Install] section")
	}
	t.Logf("Systemd service generated at %s (%d bytes)", servicePath, len(data))
}

func TestEngineReleaseNames(t *testing.T) {
	if XrayRelease.Owner != "XTLS" {
		t.Errorf("XrayRelease.Owner = %q, want XTLS", XrayRelease.Owner)
	}
	if MihomoRelease.Owner != "MetaCubeX" {
		t.Errorf("MihomoRelease.Owner = %q, want MetaCubeX", MihomoRelease.Owner)
	}
}

func TestExtractJSONField(t *testing.T) {
	data := []byte(`{"tag_name": "v1.8.4", "name": "Release"}`)
	result := extractJSONField(data, "tag_name")
	if result != "v1.8.4" {
		t.Errorf("extractJSONField = %q, want v1.8.4", result)
	}
}

func TestWriteXrayConfigCustomLogLevel(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteXrayConfig(dir, "debug", "vmess", nil)
	if err != nil {
		t.Fatalf("WriteXrayConfig failed: %v", err)
	}

	data, _ := os.ReadFile(configPath)
	if !strings.Contains(string(data), "debug") {
		t.Error("custom log level 'debug' not found in config")
	}
}

func TestWriteMihomoConfigCustomLogLevel(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteMihomoConfig(dir, "debug", "full")
	if err != nil {
		t.Fatalf("WriteMihomoConfig failed: %v", err)
	}

	data, _ := os.ReadFile(configPath)
	if !strings.Contains(string(data), "debug") {
		t.Error("custom log level 'debug' not found in config")
	}
}

func TestWriteMultipleServices(t *testing.T) {
	dir := t.TempDir()

	// Xray service
	sp1, err := WriteSystemdService(dir, "xray", "config.json", "Xray-core")
	if err != nil {
		t.Fatal(err)
	}

	// Mihomo service
	sp2, err := WriteSystemdService(dir, "mihomo", "config.yaml", "Mihomo")
	if err != nil {
		t.Fatal(err)
	}

	// Verify both exist
	for _, p := range []string{sp1, sp2} {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("service file %s does not exist", p)
		}
	}
}

func TestConfigFileNames(t *testing.T) {
	dir := t.TempDir()

	xrayPath, _ := WriteXrayConfig(dir, "info", "vmess", nil)
	if filepath.Base(xrayPath) != "config.json" {
		t.Errorf("Xray config name = %q, want config.json", filepath.Base(xrayPath))
	}

	mihomoPath, _ := WriteMihomoConfig(dir, "info", "full")
	if filepath.Base(mihomoPath) != "config.yaml" {
		t.Errorf("Mihomo config name = %q, want config.yaml", filepath.Base(mihomoPath))
	}
}

func TestWriteMihomoChinaSplitConfig(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteMihomoConfig(dir, "info", "china-split")
	if err != nil {
		t.Fatalf("WriteMihomoConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "GEOSITE,cn,DIRECT") {
		t.Error("china-split config missing GEOSITE,cn,DIRECT rule")
	}
	if !strings.Contains(content, "GEOIP,cn,DIRECT") {
		t.Error("china-split config missing GEOIP,cn,DIRECT rule")
	}
	if !strings.Contains(content, "MATCH,Proxy") {
		t.Error("china-split config missing MATCH,Proxy fallback rule")
	}
	t.Logf("China-split config generated at %s (%d bytes)", configPath, len(data))
}

func TestWriteXrayVLessConfig(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteXrayConfig(dir, "info", "vless", nil)
	if err != nil {
		t.Fatalf("WriteXrayConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "vless") {
		t.Error("vless config missing vless protocol")
	}
	if !strings.Contains(content, "xtls-rprx-vision") {
		t.Error("vless config missing xtls-rprx-vision flow")
	}
	t.Logf("VLess config generated at %s (%d bytes)", configPath, len(data))
}

func TestWriteXrayTrojanConfig(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteXrayConfig(dir, "info", "trojan", nil)
	if err != nil {
		t.Fatalf("WriteXrayConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "trojan") {
		t.Error("trojan config missing trojan protocol")
	}
	t.Logf("Trojan config generated at %s (%d bytes)", configPath, len(data))
}

func TestWriteXrayShadowsocksConfig(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteXrayConfig(dir, "info", "shadowsocks", nil)
	if err != nil {
		t.Fatalf("WriteXrayConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "shadowsocks") {
		t.Error("shadowsocks config missing shadowsocks protocol")
	}
	t.Logf("Shadowsocks config generated at %s (%d bytes)", configPath, len(data))
}

func TestWriteXrayVLessRealityConfig(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteXrayConfig(dir, "info", "vless-reality", nil)
	if err != nil {
		t.Fatalf("WriteXrayConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "reality") {
		t.Error("vless-reality config missing reality security")
	}
	if !strings.Contains(content, "vless") {
		t.Error("vless-reality config missing vless protocol")
	}
	t.Logf("VLess+Reality config generated at %s (%d bytes)", configPath, len(data))
}

func TestWriteXrayVLessGRPCConfig(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteXrayConfig(dir, "info", "vless-grpc", nil)
	if err != nil {
		t.Fatalf("WriteXrayConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "grpc") {
		t.Error("vless-grpc config missing grpc network")
	}
	if !strings.Contains(content, "vless") {
		t.Error("vless-grpc config missing vless protocol")
	}
	t.Logf("VLess+gRPC config generated at %s (%d bytes)", configPath, len(data))
}

func TestWriteXrayAnyTLSConfig(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteXrayConfig(dir, "info", "anytls", nil)
	if err != nil {
		t.Fatalf("WriteXrayConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "anytls") {
		t.Error("anytls config missing anytls protocol")
	}
	t.Logf("AnyTLS config generated at %s (%d bytes)", configPath, len(data))
}

func TestWriteXrayShadowsocks2022Config(t *testing.T) {
	dir := t.TempDir()
	configPath, err := WriteXrayConfig(dir, "info", "ss2022", nil)
	if err != nil {
		t.Fatalf("WriteXrayConfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "2022-blake3") {
		t.Error("ss2022 config missing 2022-blake3 cipher")
	}
	t.Logf("Shadowsocks-2022 config generated at %s (%d bytes)", configPath, len(data))
}
