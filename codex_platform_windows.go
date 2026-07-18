//go:build windows

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type windowsAppPackage struct {
	Path              string `json:"path"`
	Executable        string `json:"executable"`
	Version           string `json:"version"`
	PackageFamilyName string `json:"packageFamilyName"`
	AppID             string `json:"appId"`
}

func discoverCodex() (CodexInstall, error) {
	if install, ok := discoverPackagedCodex(); ok {
		return install, nil
	}

	candidates := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "ChatGPT", "ChatGPT.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Codex", "Codex.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "OpenAI", "ChatGPT.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "ChatGPT", "ChatGPT.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Codex", "Codex.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "ChatGPT", "ChatGPT.exe"),
	}
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return CodexInstall{Path: filepath.Dir(candidate), Executable: candidate}, nil
		}
	}
	return CodexInstall{}, errors.New("未找到 Windows 版 ChatGPT/Codex；请先从 Microsoft Store 安装官方应用")
}

func discoverPackagedCodex() (CodexInstall, bool) {
	const script = `$ErrorActionPreference = 'SilentlyContinue'
$items = @()
$packages = @(Get-AppxPackage | Where-Object {
  $_.Name -match '(?i)(openai|chatgpt|codex)' -or $_.PackageFamilyName -match '(?i)(openai|chatgpt|codex)'
})
foreach ($package in $packages) {
  $manifest = Get-AppxPackageManifest -Package $package.PackageFullName
  foreach ($app in @($manifest.Package.Applications.Application)) {
    $relative = [string]$app.Executable
    if (-not $relative -or $relative -notmatch '(?i)(chatgpt|codex).*\.exe$') { continue }
    $executable = Join-Path ([string]$package.InstallLocation) $relative
    if (-not (Test-Path -LiteralPath $executable)) { continue }
    $items += [pscustomobject]@{
      path = [string]$package.InstallLocation
      executable = $executable
      version = [string]$package.Version
      packageFamilyName = [string]$package.PackageFamilyName
      appId = [string]$app.Id
    }
  }
}
$items | ConvertTo-Json -Compress`
	output, err := powershell(script)
	if err != nil || strings.TrimSpace(output) == "" {
		return CodexInstall{}, false
	}

	var packages []windowsAppPackage
	if err := json.Unmarshal([]byte(output), &packages); err != nil {
		var single windowsAppPackage
		if json.Unmarshal([]byte(output), &single) != nil {
			return CodexInstall{}, false
		}
		packages = []windowsAppPackage{single}
	}
	for _, item := range packages {
		if item.Executable == "" {
			continue
		}
		return CodexInstall{
			Path: item.Path, Executable: item.Executable, Version: item.Version,
			PackageFamilyName: item.PackageFamilyName, AppID: item.AppID,
		}, true
	}
	return CodexInstall{}, false
}

func quitCodex(install CodexInstall) error {
	if !processRunning(install.Executable) {
		return nil
	}
	processName := strings.TrimSuffix(filepath.Base(install.Executable), filepath.Ext(install.Executable))
	script := fmt.Sprintf(`Get-Process -Name '%s' -ErrorAction SilentlyContinue | ForEach-Object { [void]$_.CloseMainWindow() }`, escapePowerShell(processName))
	if _, err := powershell(script); err != nil {
		return fmt.Errorf("无法请求 Codex 正常退出：%w", err)
	}
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		if !processRunning(install.Executable) {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return errors.New("Codex 仍在运行，请手动退出后重试；Codex Canvas 不会强制结束进程")
}

func processRunning(executable string) bool {
	imageName := filepath.Base(executable)
	output, err := exec.Command("tasklist.exe", "/FI", "IMAGENAME eq "+imageName, "/FO", "CSV", "/NH").CombinedOutput()
	return err == nil && strings.Contains(strings.ToLower(string(output)), strings.ToLower(imageName))
}

func launchCodexWithDebug(install CodexInstall, port int) error {
	args := []string{
		"--remote-debugging-address=127.0.0.1",
		fmt.Sprintf("--remote-debugging-port=%d", port),
		fmt.Sprintf("--remote-allow-origins=http://127.0.0.1:%d", port),
	}
	cmd := exec.Command(install.Executable, args...)
	cmd.Dir = filepath.Dir(install.Executable)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("无法使用主题参数启动 Codex：%w", err)
	}
	return cmd.Process.Release()
}

func launchCodexNormal(install CodexInstall) error {
	if install.PackageFamilyName != "" && install.AppID != "" {
		target := `shell:AppsFolder\` + install.PackageFamilyName + "!" + install.AppID
		cmd := exec.Command("explorer.exe", target)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("无法启动 Codex：%w", err)
		}
		return cmd.Process.Release()
	}
	cmd := exec.Command(install.Executable)
	cmd.Dir = filepath.Dir(install.Executable)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("无法启动 Codex：%w", err)
	}
	return cmd.Process.Release()
}

func powershell(script string) (string, error) {
	cmd := exec.Command("powershell.exe", "-NoLogo", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func escapePowerShell(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
