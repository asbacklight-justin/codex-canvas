//go:build !darwin && !windows

package main

import "errors"

func discoverCodex() (CodexInstall, error) {
	return CodexInstall{}, errors.New("当前版本支持 macOS 和 Windows")
}

func quitCodex(CodexInstall) error {
	return errors.New("当前版本支持 macOS 和 Windows")
}

func processRunning(string) bool { return false }

func launchCodexWithDebug(CodexInstall, int) error {
	return errors.New("当前版本支持 macOS 和 Windows")
}

func launchCodexNormal(CodexInstall) error {
	return errors.New("当前版本支持 macOS 和 Windows")
}
