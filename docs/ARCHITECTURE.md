# Architecture and security model

[简体中文](ARCHITECTURE.zh-CN.md) | **English**

## Module boundaries

```text
Wails UI
  ├─ ThemeConfig (cross-platform)
  ├─ live preview (cross-platform)
  └─ typed App bindings
       └─ CodexService
            ├─ image validation and local persistence (cross-platform)
            ├─ controlled CSS generator (cross-platform)
            ├─ CDP target discovery and WebSocket injection (cross-platform)
            └─ Codex lifecycle adapters (macOS and Windows)
```

The theme configuration is data, not executable code. Codex Canvas does not accept arbitrary CSS, JavaScript, remote URLs, or third-party theme bundles. A selected image is validated, copied into the app data directory, and converted to a local `data:` URL when the style is generated.

## macOS lifecycle

1. Candidate applications are validated using bundle identifier `com.openai.codex`. Both `Codex.app` and the compatible `ChatGPT.app` display name are supported.
2. An Apple Event asks Codex to quit normally. Codex Canvas does not force-kill it.
3. The operating system allocates a random free TCP port on `127.0.0.1`.
4. `/usr/bin/open` launches a new Codex instance with loopback-only CDP flags.
5. The client polls the local `/json` endpoint, connects to eligible page targets, and applies the style to the current page and future reloads.
6. Restore removes the injected style when the endpoint is available or relaunches Codex normally when it is not.

## Trust and threat model

- **Application integrity:** no files inside the Codex bundle are edited, so its code signature remains intact.
- **Local data:** settings and selected backgrounds are stored under `~/Library/Application Support/Codex Canvas` with user-only directory permissions where possible.
- **Network:** theme operations use a random loopback endpoint. No application data is sent to a Codex Canvas server because no such server exists.
- **CDP exposure:** while a themed Codex session is running, another process under the same macOS user could theoretically discover and access the loopback debugging endpoint. Restore the official appearance or launch Codex normally for sensitive sessions.
- **Compatibility:** Codex internal DOM and design tokens can change. The injection is session-scoped and fails without modifying the official app.

## Universal 2 build

`scripts/build-macos.sh` compiles the Go/Wails application twice using the installed macOS SDK:

- `GOARCH=amd64` with Clang target `x86_64`
- `GOARCH=arm64` with Clang target `arm64`

`lipo` combines both Mach-O executables. The resulting bundle is ad-hoc signed for local testing. Public production distribution requires a Developer ID Application signature and Apple notarization.

## Windows status

The Windows adapter discovers Microsoft Store Appx packages through PowerShell and parses the current manifest executable, with `%LOCALAPPDATA%` and `%ProgramFiles%` fallbacks. It requests graceful window closure, launches with equivalent loopback-only CDP flags, and stores data through `os.UserConfigDir()` under `%AppData%\Codex Canvas`.

The x64 preview is cross-compiled as a Windows GUI PE executable with embedded Wails resources. Real-device testing, Authenticode signing, and an NSIS/MSI installer remain release gates.

## Release gates

- Go tests, frontend production build, and Universal architecture checks pass.
- macOS Developer ID signing and notarization pass for production releases.
- SHA-256 checksums and release notes are published.
- Apply, navigation, restore, normal launch, and update compatibility are tested on Intel, Apple Silicon, and Windows x64 hardware.
