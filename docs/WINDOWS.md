# Windows preview

[简体中文](WINDOWS.zh-CN.md) | **English**

Codex Canvas 1.1.0 includes an experimental Windows x64 build. OpenAI distributes the current ChatGPT/Codex desktop app through Microsoft Store, so Codex Canvas discovers the installed Appx package and reads its manifest instead of relying on a fixed directory.

## Requirements

- Windows 10 or Windows 11 x64
- The official [ChatGPT desktop app for Windows](https://learn.chatgpt.com/docs/windows/windows-app), installed from Microsoft Store
- Microsoft Edge WebView2 Runtime

Official command-line installation for the ChatGPT app:

```powershell
winget install --id 9PLM9XGG6VKS -s msstore
```

## Install Codex Canvas

Download `Codex-Canvas-1.1.0-Windows-x64.zip` from this repository's GitHub Release, verify the SHA-256 checksum, extract it, and run `Codex-Canvas.exe`.

The preview is not Authenticode-signed, so Windows SmartScreen may warn on first launch. Only bypass that warning for an archive downloaded from the official Release page.

## Implementation

- Discovers the current user's ChatGPT/Codex Appx package through PowerShell.
- Reads `AppxManifest.xml` through `Get-AppxPackageManifest` to find the executable after Store updates.
- Requests graceful window close and never uses `taskkill` to force termination.
- Launches the desktop executable with a random loopback-only CDP port.
- Uses the same controlled CSS generator and restore path as macOS.

## Preview limitation

The binary is cross-compiled and its PE architecture, embedded resources, trim-path output, and Windows build tags are validated. The project currently has no Windows test machine, so Microsoft Store activation behavior still needs real-device confirmation. Please report results in GitHub Issues and never attach private prompts, conversations, tokens, or personal paths.
