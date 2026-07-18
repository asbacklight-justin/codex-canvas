package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxImageBytes = 16 * 1024 * 1024
	styleID       = "codex-canvas-theme"
)

type ThemeConfig struct {
	ImagePath      string `json:"imagePath"`
	ImageName      string `json:"imageName"`
	Overlay        int    `json:"overlay"`
	SurfaceOpacity int    `json:"surfaceOpacity"`
	SidebarOpacity int    `json:"sidebarOpacity"`
	Blur           int    `json:"blur"`
	Radius         int    `json:"radius"`
	Scale          int    `json:"scale"`
	Position       string `json:"position"`
	Accent         string `json:"accent"`
	Active         bool   `json:"active"`
	LastPort       int    `json:"lastPort"`
}

type AppStatus struct {
	Platform      string      `json:"platform"`
	CodexFound    bool        `json:"codexFound"`
	CodexPath     string      `json:"codexPath"`
	CodexVersion  string      `json:"codexVersion"`
	Supported     bool        `json:"supported"`
	Active        bool        `json:"active"`
	SavedTheme    ThemeConfig `json:"savedTheme"`
	PreviewURL    string      `json:"previewUrl"`
	StatusMessage string      `json:"statusMessage"`
	Compatibility string      `json:"compatibility"`
}

type ImageSelection struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	PreviewURL string `json:"previewUrl"`
	Cancelled  bool   `json:"cancelled"`
}

type ActionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Port    int    `json:"port,omitempty"`
}

type CodexInstall struct {
	Path       string
	Executable string
	Version    string
}

type CodexService struct {
	mu sync.Mutex
}

func NewCodexService() *CodexService { return &CodexService{} }

func defaultTheme() ThemeConfig {
	return ThemeConfig{
		Overlay: 48, SurfaceOpacity: 72, SidebarOpacity: 82,
		Blur: 0, Radius: 18, Scale: 100, Position: "center", Accent: "#8b7cff",
	}
}

func (s *CodexService) Status() AppStatus {
	theme, _ := s.loadSettings()
	install, err := discoverCodex()
	status := AppStatus{
		Platform:      goruntime.GOOS,
		Supported:     goruntime.GOOS == "darwin",
		SavedTheme:    theme,
		Active:        theme.Active,
		Compatibility: "运行时主题 · 不修改应用签名",
	}
	if err != nil {
		status.StatusMessage = "尚未找到 Codex Desktop"
		return status
	}
	status.CodexFound = true
	status.CodexPath = install.Path
	status.CodexVersion = install.Version
	status.StatusMessage = "已找到 Codex Desktop"
	if theme.ImagePath != "" {
		if selection, imageErr := loadImageSelection(theme.ImagePath); imageErr == nil {
			status.PreviewURL = selection.PreviewURL
		}
	}
	return status
}

func (s *CodexService) Apply(theme ThemeConfig) (ActionResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if goruntime.GOOS != "darwin" {
		return ActionResult{}, errors.New("当前版本先支持 macOS，Windows 适配层已预留")
	}
	install, err := discoverCodex()
	if err != nil {
		return ActionResult{}, err
	}
	if err := validateTheme(theme); err != nil {
		return ActionResult{}, err
	}

	savedImage, err := s.persistImage(theme.ImagePath)
	if err != nil {
		return ActionResult{}, err
	}
	theme.ImagePath = savedImage
	if theme.ImageName == "" {
		theme.ImageName = filepath.Base(savedImage)
	}

	if err := quitCodex(install); err != nil {
		return ActionResult{}, err
	}
	port, err := availablePort()
	if err != nil {
		return ActionResult{}, fmt.Errorf("无法准备本地主题通道：%w", err)
	}
	if err := launchCodexWithDebug(install, port); err != nil {
		return ActionResult{}, err
	}

	css, err := buildThemeCSS(theme)
	if err != nil {
		return ActionResult{}, err
	}
	if err := waitAndInject(port, install.Version, installStyleScript(css), 25*time.Second); err != nil {
		return ActionResult{}, fmt.Errorf("Codex 已启动，但主题应用失败：%w", err)
	}

	theme.Active = true
	theme.LastPort = port
	if err := s.saveSettings(theme); err != nil {
		return ActionResult{}, fmt.Errorf("主题已应用，但保存设置失败：%w", err)
	}
	return ActionResult{Success: true, Message: "主题已应用，Codex Canvas 没有修改 Codex.app", Port: port}, nil
}

func (s *CodexService) Restore() (ActionResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	theme, _ := s.loadSettings()
	install, installErr := discoverCodex()
	liveRestored := false
	if theme.LastPort > 0 {
		remove := fmt.Sprintf(`(() => { document.getElementById(%q)?.remove(); document.documentElement.removeAttribute("data-codex-canvas"); return true; })()`, styleID)
		if err := injectAll(theme.LastPort, remove); err == nil {
			liveRestored = true
		}
	}
	theme.Active = false
	theme.LastPort = 0
	if err := s.saveSettings(theme); err != nil {
		return ActionResult{}, err
	}
	if liveRestored {
		return ActionResult{Success: true, Message: "已实时恢复官方外观，下次也会正常启动"}, nil
	}
	if installErr != nil {
		return ActionResult{Success: true, Message: "主题状态已清除；未找到 Codex，未执行重启"}, nil
	}
	if err := quitCodex(install); err != nil {
		return ActionResult{}, err
	}
	if err := launchCodexNormal(install); err != nil {
		return ActionResult{}, err
	}
	return ActionResult{Success: true, Message: "已重新启动 Codex 并恢复官方外观"}, nil
}

func validateTheme(t ThemeConfig) error {
	if t.ImagePath == "" {
		return errors.New("请先选择一张背景图片")
	}
	if _, err := os.Stat(t.ImagePath); err != nil {
		return errors.New("所选图片已不存在，请重新选择")
	}
	if t.Overlay < 0 || t.Overlay > 90 || t.SurfaceOpacity < 15 || t.SurfaceOpacity > 100 || t.SidebarOpacity < 15 || t.SidebarOpacity > 100 {
		return errors.New("透明度参数超出安全范围")
	}
	if t.Blur < 0 || t.Blur > 30 || t.Radius < 0 || t.Radius > 32 || t.Scale < 100 || t.Scale > 150 {
		return errors.New("视觉参数超出安全范围")
	}
	if t.Position != "center" && t.Position != "top" && t.Position != "bottom" {
		return errors.New("不支持的图片位置")
	}
	if len(t.Accent) != 7 || t.Accent[0] != '#' {
		return errors.New("强调色格式不正确")
	}
	if _, err := strconv.ParseUint(t.Accent[1:], 16, 24); err != nil {
		return errors.New("强调色格式不正确")
	}
	return nil
}

func buildThemeCSS(t ThemeConfig) (string, error) {
	data, err := os.ReadFile(t.ImagePath)
	if err != nil {
		return "", fmt.Errorf("无法读取背景图片：%w", err)
	}
	if len(data) > maxImageBytes {
		return "", errors.New("图片不能超过 16 MB")
	}
	mime, ok := imageMIME(filepath.Ext(t.ImagePath))
	if !ok {
		return "", errors.New("图片格式不受支持")
	}
	imageURL := makeDataURL(mime, data)
	surface := float64(t.SurfaceOpacity) / 100
	sidebar := float64(t.SidebarOpacity) / 100
	overlay := float64(t.Overlay) / 100
	scale := float64(t.Scale) / 100
	return fmt.Sprintf(`
:root, :root.electron-dark, :root.electron-light {
  color-scheme: dark !important;
  --startup-background: transparent !important;
  --background: rgb(10 13 22 / 0.82) !important;
  --foreground: #f5f7fb !important;
  --surface: rgb(21 24 36 / 0.90) !important;
  --card: rgb(255 255 255 / 0.08) !important;
  --popover: rgb(28 31 44 / 0.98) !important;
  --muted: #a7aeba !important;
  --input: rgb(24 27 39 / 0.94) !important;
  --color-background-surface: rgb(10 13 22 / %.2f) !important;
  --color-background-surface-under: rgb(7 9 17 / %.2f) !important;
  --color-background-elevated-primary: rgb(21 24 36 / %.2f) !important;
  --color-background-elevated-primary-opaque: rgb(21 24 36 / %.2f) !important;
  --color-background-elevated-secondary: rgb(255 255 255 / 0.07) !important;
  --color-background-editor-opaque: rgb(12 15 25 / %.2f) !important;
  --color-background-primary: rgb(10 13 22 / 0.82) !important;
  --color-text-foreground: #f5f7fb !important;
  --color-text-foreground-secondary: #c3c8d2 !important;
  --color-text-foreground-tertiary: #949caa !important;
  --color-text-emphasis: #ffffff !important;
  --color-text-primary: #f5f7fb !important;
  --color-text-secondary: #c3c8d2 !important;
  --color-text-tertiary: #949caa !important;
  --color-text-prose: #edf0f6 !important;
  --color-token-main-surface-primary: rgb(10 13 22 / %.2f) !important;
  --color-token-side-bar-background: rgb(7 9 17 / %.2f) !important;
  --color-token-bg-primary: rgb(10 13 22 / %.2f) !important;
  --color-token-bg-secondary: rgb(21 24 36 / %.2f) !important;
  --color-token-foreground: #f5f7fb !important;
  --color-token-text-primary: #f5f7fb !important;
  --color-token-text-secondary: #c3c8d2 !important;
  --color-token-text-tertiary: #949caa !important;
  --color-token-description-foreground: #aab1bd !important;
  --color-token-disabled-foreground: #717a89 !important;
  --color-token-conversation-body: #d5dae3 !important;
  --color-token-conversation-header: #9da5b3 !important;
  --color-token-icon-foreground: #d7dce5 !important;
  --color-token-input-background: rgb(24 27 39 / 0.94) !important;
  --color-token-input-foreground: #f5f7fb !important;
  --color-token-input-placeholder-foreground: #929aa8 !important;
  --vscode-editor-background: rgb(12 15 25 / %.2f) !important;
  --vscode-editor-foreground: #f5f7fb !important;
  --vscode-sideBar-background: rgb(7 9 17 / %.2f) !important;
  --vscode-sideBar-foreground: #dfe3ea !important;
  --vscode-foreground: #f5f7fb !important;
  --vscode-descriptionForeground: #aab1bd !important;
  --vscode-input-background: rgb(24 27 39 / 0.94) !important;
  --vscode-input-foreground: #f5f7fb !important;
  --vscode-input-placeholderForeground: #929aa8 !important;
  --color-accent: %s !important;
  --codex-canvas-radius: %dpx;
}
html, body, #root { background: transparent !important; }
html::before {
  content: ""; position: fixed; inset: -32px; z-index: -2; pointer-events: none;
  background: url(%q) %s / cover no-repeat fixed;
  filter: blur(%dpx); transform: scale(%.3f);
}
html::after {
  content: ""; position: fixed; inset: 0; z-index: -1; pointer-events: none;
  background: rgb(3 5 12 / %.2f);
}
body { min-height: 100vh; color: var(--color-text-foreground); }
textarea, input, [contenteditable="true"] {
  color: var(--color-token-input-foreground) !important;
  caret-color: #f5f7fb !important;
}
textarea::placeholder, input::placeholder {
  color: var(--color-token-input-placeholder-foreground) !important;
  opacity: 1 !important;
}
[data-radix-popper-content-wrapper] > *, [role="dialog"] { border-radius: min(var(--codex-canvas-radius), 20px) !important; }
`, surface, sidebar, surface, minFloat(1, surface+0.12), surface, surface, sidebar, surface, surface, surface, sidebar, t.Accent, t.Radius, imageURL, t.Position, t.Blur, scale, overlay), nil
}

func installStyleScript(css string) string {
	cssJSON, _ := json.Marshal(css)
	return fmt.Sprintf(`(() => {
  const apply = () => {
    let style = document.getElementById(%q);
    if (!style) { style = document.createElement("style"); style.id = %q; (document.head || document.documentElement).appendChild(style); }
    style.textContent = %s;
    document.documentElement.setAttribute("data-codex-canvas", "active");
    return true;
  };
  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", apply, { once: true }); else apply();
  return true;
})()`, styleID, styleID, cssJSON)
}

type cdpTarget struct {
	Type                 string `json:"type"`
	Title                string `json:"title"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

func waitAndInject(port int, version, expression string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		if err := injectAll(port, expression); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(450 * time.Millisecond)
	}
	if lastErr == nil {
		lastErr = errors.New("调试端点未就绪")
	}
	return fmt.Errorf("版本 %s：%w", version, lastErr)
}

func injectAll(port int, expression string) error {
	client := &http.Client{Timeout: 1200 * time.Millisecond}
	response, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/json", port))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("本地端点返回 %d", response.StatusCode)
	}
	var targets []cdpTarget
	if err := json.NewDecoder(io.LimitReader(response.Body, 2<<20)).Decode(&targets); err != nil {
		return err
	}
	count := 0
	for _, target := range targets {
		if target.Type != "page" || target.WebSocketDebuggerURL == "" || strings.HasPrefix(target.URL, "devtools://") {
			continue
		}
		if err := injectTarget(target, expression, port); err != nil {
			return err
		}
		count++
	}
	if count == 0 {
		return errors.New("还没有发现 Codex 页面")
	}
	return nil
}

func injectTarget(target cdpTarget, expression string, port int) error {
	header := http.Header{}
	header.Set("Origin", fmt.Sprintf("http://127.0.0.1:%d", port))
	conn, _, err := websocket.DefaultDialer.Dial(target.WebSocketDebuggerURL, header)
	if err != nil {
		return err
	}
	defer conn.Close()
	_ = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	newDocument := map[string]any{"id": 1, "method": "Page.addScriptToEvaluateOnNewDocument", "params": map[string]any{"source": expression}}
	evaluate := map[string]any{"id": 2, "method": "Runtime.evaluate", "params": map[string]any{"expression": expression, "returnByValue": true}}
	if err := conn.WriteJSON(newDocument); err != nil {
		return err
	}
	if err := conn.WriteJSON(evaluate); err != nil {
		return err
	}
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		var message struct {
			ID     int             `json:"id"`
			Error  json.RawMessage `json:"error"`
			Result struct {
				ExceptionDetails json.RawMessage `json:"exceptionDetails"`
			} `json:"result"`
		}
		if json.Unmarshal(raw, &message) != nil || message.ID != 2 {
			continue
		}
		if len(message.Error) > 0 || len(message.Result.ExceptionDetails) > 0 {
			return errors.New("Codex 拒绝了样式注入")
		}
		return nil
	}
}

func discoverCodex() (CodexInstall, error) {
	if goruntime.GOOS != "darwin" {
		return CodexInstall{}, errors.New("当前平台尚未实现 Codex 自动发现")
	}
	home, _ := os.UserHomeDir()
	candidates := []string{
		"/Applications/Codex.app", filepath.Join(home, "Applications", "Codex.app"),
		"/Applications/ChatGPT.app", filepath.Join(home, "Applications", "ChatGPT.app"),
	}
	for _, candidate := range candidates {
		plist := filepath.Join(candidate, "Contents", "Info.plist")
		if _, err := os.Stat(plist); err != nil {
			continue
		}
		identifier := plistValue(plist, "CFBundleIdentifier")
		if identifier != "com.openai.codex" {
			continue
		}
		executableName := plistValue(plist, "CFBundleExecutable")
		if executableName == "" {
			executableName = strings.TrimSuffix(filepath.Base(candidate), ".app")
		}
		executable := filepath.Join(candidate, "Contents", "MacOS", executableName)
		if info, err := os.Stat(executable); err != nil || info.Mode()&0o111 == 0 {
			continue
		}
		return CodexInstall{Path: candidate, Executable: executable, Version: plistValue(plist, "CFBundleShortVersionString")}, nil
	}
	return CodexInstall{}, errors.New("未找到 Codex Desktop；请先安装官方应用")
}

func plistValue(plist, key string) string {
	output, err := exec.Command("/usr/bin/plutil", "-extract", key, "raw", "-o", "-", plist).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func quitCodex(install CodexInstall) error {
	if !processRunning(install.Executable) {
		return nil
	}
	cmd := exec.Command("/usr/bin/osascript", "-e", `tell application id "com.openai.codex" to quit`)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("无法关闭 Codex：%s", strings.TrimSpace(string(output)))
	}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if !processRunning(install.Executable) {
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	return errors.New("Codex 仍在运行，请手动退出后重试")
}

func processRunning(executable string) bool {
	cmd := exec.Command("/usr/bin/pgrep", "-f", executable)
	return cmd.Run() == nil
}

func launchCodexWithDebug(install CodexInstall, port int) error {
	args := []string{"-na", install.Path, "--args",
		"--remote-debugging-address=127.0.0.1",
		fmt.Sprintf("--remote-debugging-port=%d", port),
		fmt.Sprintf("--remote-allow-origins=http://127.0.0.1:%d", port),
	}
	if output, err := exec.Command("/usr/bin/open", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("无法启动 Codex：%s", strings.TrimSpace(string(output)))
	}
	return nil
}

func launchCodexNormal(install CodexInstall) error {
	if output, err := exec.Command("/usr/bin/open", "-a", install.Path).CombinedOutput(); err != nil {
		return fmt.Errorf("无法启动 Codex：%s", strings.TrimSpace(string(output)))
	}
	return nil
}

func availablePort() (int, error) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func (s *CodexService) dataDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "Codex Canvas"), nil
}

func (s *CodexService) settingsPath() (string, error) {
	dir, err := s.dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

func (s *CodexService) loadSettings() (ThemeConfig, error) {
	theme := defaultTheme()
	path, err := s.settingsPath()
	if err != nil {
		return theme, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return theme, nil
	}
	if err != nil {
		return theme, err
	}
	if err := json.Unmarshal(data, &theme); err != nil {
		return defaultTheme(), err
	}
	return theme, nil
}

func (s *CodexService) saveSettings(theme ThemeConfig) error {
	path, err := s.settingsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		return err
	}
	temp := path + ".tmp"
	if err := os.WriteFile(temp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(temp, path)
}

func (s *CodexService) persistImage(source string) (string, error) {
	info, err := os.Stat(source)
	if err != nil {
		return "", fmt.Errorf("无法读取图片：%w", err)
	}
	if info.Size() > maxImageBytes {
		return "", errors.New("图片不能超过 16 MB")
	}
	ext := strings.ToLower(filepath.Ext(source))
	if _, ok := imageMIME(ext); !ok {
		return "", errors.New("图片格式不受支持")
	}
	dir, err := s.dataDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	destination := filepath.Join(dir, "wallpaper"+ext)
	if samePath(source, destination) {
		return destination, nil
	}
	input, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer input.Close()
	temp := destination + ".tmp"
	output, err := os.OpenFile(temp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return "", err
	}
	_, copyErr := io.Copy(output, io.LimitReader(input, maxImageBytes+1))
	closeErr := output.Close()
	if copyErr != nil {
		return "", copyErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	if err := os.Rename(temp, destination); err != nil {
		return "", err
	}
	return destination, nil
}

func imageMIME(ext string) (string, bool) {
	switch strings.ToLower(ext) {
	case ".png":
		return "image/png", true
	case ".jpg", ".jpeg":
		return "image/jpeg", true
	case ".webp":
		return "image/webp", true
	case ".gif":
		return "image/gif", true
	default:
		return "", false
	}
}

func makeDataURL(mime string, data []byte) string {
	var buffer bytes.Buffer
	buffer.Grow(len(data)*4/3 + 32)
	buffer.WriteString("data:")
	buffer.WriteString(mime)
	buffer.WriteString(";base64,")
	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
	_, _ = encoder.Write(data)
	_ = encoder.Close()
	return buffer.String()
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func samePath(a, b string) bool {
	aa, errA := filepath.Abs(a)
	bb, errB := filepath.Abs(b)
	return errA == nil && errB == nil && aa == bb
}
