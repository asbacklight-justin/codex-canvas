package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func testTheme(imagePath string) ThemeConfig {
	return ThemeConfig{
		ImagePath: imagePath, Overlay: 48, SurfaceOpacity: 72,
		SidebarOpacity: 82, Blur: 4, Radius: 18, Scale: 108,
		Position: "center", Accent: "#8b7cff",
	}
}

func writeTestPNG(t *testing.T) string {
	t.Helper()
	// A valid transparent 1×1 PNG.
	data, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M/wHwAF/gL+XwFujwAAAABJRU5ErkJggg==")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "wallpaper.png")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestValidateTheme(t *testing.T) {
	path := writeTestPNG(t)
	if err := validateTheme(testTheme(path)); err != nil {
		t.Fatalf("valid theme rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*ThemeConfig)
	}{
		{"missing image", func(v *ThemeConfig) { v.ImagePath = "" }},
		{"unsafe overlay", func(v *ThemeConfig) { v.Overlay = 95 }},
		{"unsafe scale", func(v *ThemeConfig) { v.Scale = 151 }},
		{"invalid position", func(v *ThemeConfig) { v.Position = "left" }},
		{"invalid accent", func(v *ThemeConfig) { v.Accent = "purple" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := testTheme(path)
			test.mutate(&value)
			if err := validateTheme(value); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestBuildThemeCSSIsControlledAndSelfContained(t *testing.T) {
	path := writeTestPNG(t)
	css, err := buildThemeCSS(testTheme(path))
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"data:image/png;base64,", "--color-background-surface", "#8b7cff", "blur(4px)", "scale(1.080)",
		"--foreground: #f5f7fb", "--color-token-input-foreground: #f5f7fb", "--vscode-input-background:",
		`:root.electron-light`, `[contenteditable="true"]`,
	} {
		if !strings.Contains(css, expected) {
			t.Fatalf("CSS does not contain %q", expected)
		}
	}
	if strings.Contains(css, "%!") {
		t.Fatalf("CSS contains an unresolved format marker: %s", css)
	}
	script := installStyleScript(css)
	if !strings.Contains(script, styleID) || !strings.Contains(script, "DOMContentLoaded") {
		t.Fatal("install script is not persistent across initial document load")
	}
}

func TestLoadImageSelection(t *testing.T) {
	path := writeTestPNG(t)
	selection, err := loadImageSelection(path)
	if err != nil {
		t.Fatal(err)
	}
	if selection.Name != "wallpaper.png" || !strings.HasPrefix(selection.PreviewURL, "data:image/png;base64,") {
		t.Fatalf("unexpected selection: %#v", selection)
	}
}

func TestBuiltinPresetsAreCompleteAndEmbeddable(t *testing.T) {
	presets, err := listBuiltinPresets()
	if err != nil {
		t.Fatal(err)
	}
	if len(presets) != 5 {
		t.Fatalf("got %d presets, want 5", len(presets))
	}
	seen := make(map[string]bool)
	for _, preset := range presets {
		if seen[preset.ID] {
			t.Fatalf("duplicate preset %q", preset.ID)
		}
		seen[preset.ID] = true
		if !strings.HasPrefix(preset.PreviewURL, "data:image/jpeg;base64,") || len(preset.PreviewURL) < 10_000 {
			t.Fatalf("preset %q has an invalid preview", preset.ID)
		}
	}
}

func TestInjectAllUsesOnlyPageTargets(t *testing.T) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		if errors.Is(err, syscall.EPERM) {
			t.Skip("sandbox does not allow a loopback listener")
		}
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	received := make(chan string, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/json", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode([]cdpTarget{
			{Type: "service_worker", WebSocketDebuggerURL: "ws://invalid"},
			{Type: "page", Title: "Codex", URL: "file:///app.html", WebSocketDebuggerURL: "ws://127.0.0.1:" + strings.TrimSpace(strings.Split(listener.Addr().String(), ":")[1]) + "/page"},
		})
	})
	mux.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		conn, upgradeErr := upgrader.Upgrade(w, r, nil)
		if upgradeErr != nil {
			return
		}
		defer conn.Close()
		for i := 0; i < 2; i++ {
			_, raw, readErr := conn.ReadMessage()
			if readErr != nil {
				return
			}
			var message struct {
				ID     int    `json:"id"`
				Method string `json:"method"`
			}
			_ = json.Unmarshal(raw, &message)
			if message.ID == 2 {
				received <- message.Method
				_ = conn.WriteJSON(map[string]any{"id": 2, "result": map[string]any{"result": map[string]any{"type": "boolean", "value": true}}})
			}
		}
	})
	server := &http.Server{Handler: mux, ReadHeaderTimeout: time.Second}
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() { _ = server.Close() })

	if err := injectAll(port, "(() => true)()"); err != nil {
		t.Fatal(err)
	}
	select {
	case method := <-received:
		if method != "Runtime.evaluate" {
			t.Fatalf("unexpected method %q", method)
		}
	case <-time.After(time.Second):
		t.Fatal("did not receive Runtime.evaluate")
	}
}
