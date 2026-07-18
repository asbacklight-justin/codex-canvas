//go:build darwin

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func discoverCodex() (CodexInstall, error) {
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
		if plistValue(plist, "CFBundleIdentifier") != "com.openai.codex" {
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
	return exec.Command("/usr/bin/pgrep", "-f", executable).Run() == nil
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
