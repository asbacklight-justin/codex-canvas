# Codex Canvas

[简体中文](README.zh-CN.md) | **English**

Codex Canvas is a local-first, reversible background and theme manager for Codex Desktop on macOS and Windows. It lets you choose a local image or one of five built-in artworks, tune readability, preview the result, and relaunch Codex through a temporary local theming channel.

**Version 1.1.0** · Support: [asbacklight@gmail.com](mailto:asbacklight@gmail.com)

> Codex Canvas is an independent open-source project. It is not affiliated with, endorsed by, or sponsored by OpenAI. Codex and OpenAI are trademarks of their respective owners.

## Highlights

- Native macOS image picker for PNG, JPEG, WebP, and GIF files up to 16 MB
- Five built-in 16:9 backgrounds: technology, original anime-inspired art, city, nature, and animal
- Chinese and English UI with system-language detection, English fallback, and a persistent manual switch
- Live Codex-style preview
- Readability controls for overlay, content/sidebar opacity, blur, scale, position, radius, and accent color
- High-contrast foreground and input colors for both light and dark Codex themes
- Automatic discovery of `/Applications/Codex.app` and compatible `/Applications/ChatGPT.app` installations
- One-click apply and restore without modifying `Codex.app`, `app.asar`, application resources, or the code signature
- Universal 2 build that runs natively on Intel and Apple Silicon Macs
- Experimental Windows 10/11 x64 build with Microsoft Store Appx discovery
- No uploads, analytics, telemetry, remote theme packages, or arbitrary user-supplied CSS/JavaScript

## Download and install

Download a ready-to-use build from [GitHub Releases](../../releases/latest):

- macOS: `Codex-Canvas-1.1.0-macOS-Universal.zip` for Intel and Apple Silicon
- Windows: `Codex-Canvas-1.1.0-Windows-x64.zip` for Windows 10/11 x64 (preview)

For macOS, unzip the archive and move **Codex Canvas.app** to `/Applications`. For Windows, extract the archive and run `Codex-Canvas.exe`.

The current community build uses an ad-hoc signature and is not Apple-notarized. On first launch, macOS may require **Control-click → Open**. See the [installation guide](docs/INSTALLATION.md) for details and checksums.

macOS requirements:

- macOS 12 or later
- Intel (`x86_64`) or Apple Silicon (`arm64`)
- Codex Desktop installed in `/Applications`

Windows requires the official Microsoft Store ChatGPT desktop app and WebView2. See the [Windows preview guide](docs/WINDOWS.md). The first Windows build is cross-compiled and awaits community real-device confirmation.

## How it works

Codex Canvas does not patch the Codex application. When you apply a theme, it:

1. asks Codex to quit normally;
2. allocates a random loopback-only port on `127.0.0.1`;
3. relaunches Codex with a temporary Chromium DevTools Protocol endpoint;
4. injects a generated style element into Codex pages;
5. keeps the theme active only for that Codex session.

Launching Codex normally from the Dock restores the official appearance. You can also use **Restore official appearance** inside Codex Canvas.

Read the [architecture and security model](docs/ARCHITECTURE.md) and [privacy policy](docs/PRIVACY.md) before use.

## Build from source

Prerequisites: Go 1.23+, Node.js 20+, Wails 2.11+, and macOS Command Line Tools.

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.11.0
cd frontend
npm ci
npm run build
cd ..
go test ./...
./scripts/build-macos.sh
```

The build script compiles both `x86_64` and `arm64` binaries, merges them with `lipo`, and writes the Universal 2 app to `build/bin/Codex Canvas.app`.

For development:

```bash
cd frontend && npm ci && cd ..
wails dev
```

## Documentation

- [Installation](docs/INSTALLATION.md) · [安装说明](docs/INSTALLATION.zh-CN.md)
- [Architecture](docs/ARCHITECTURE.md) · [架构说明](docs/ARCHITECTURE.zh-CN.md)
- [Privacy](docs/PRIVACY.md) · [隐私说明](docs/PRIVACY.zh-CN.md)
- [Windows preview](docs/WINDOWS.md) · [Windows 预览版](docs/WINDOWS.zh-CN.md)
- [Contributing](CONTRIBUTING.md) · [参与贡献](CONTRIBUTING.zh-CN.md)
- [Security policy](SECURITY.md) · [安全政策](SECURITY.zh-CN.md)
- [Changelog](CHANGELOG.md)

## Windows status

Windows installation discovery, graceful lifecycle handling, CDP launch, and x64 packaging are implemented. Remaining release gates are real-device compatibility testing, Authenticode signing, and installer packaging. Contributions and test reports are welcome; see the [Windows guide](docs/WINDOWS.md) and [architecture document](docs/ARCHITECTURE.md#windows-status).

## License

Codex Canvas and its bundled original artwork are released under the [MIT License](LICENSE), unless a file states otherwise.
