# Installation

[简体中文](INSTALLATION.zh-CN.md) | **English**

## Requirements

- macOS 12 or later
- Intel or Apple Silicon Mac
- Codex Desktop installed in `/Applications`

## Install the release build

1. Download `Codex-Canvas-1.0.0-macOS-Universal.zip` from [GitHub Releases](../../../releases/latest).
2. Verify its SHA-256 value against the checksum in the release notes.
3. Unzip the archive and move `Codex Canvas.app` into `/Applications`.
4. Control-click the app and choose **Open** on first launch if Gatekeeper displays an unidentified-developer warning.

The community build is ad-hoc signed but not Apple-notarized. Do not bypass security warnings for an archive obtained from anywhere other than this repository's Releases page.

## Apply a background

1. Open Codex Canvas.
2. Select one of the built-in backgrounds on the left, or choose a local image on the right.
3. Adjust readability controls and inspect the preview.
4. Click **Apply to Codex** and confirm the restart.

Codex Canvas asks Codex to quit before relaunching it. Save or send unfinished input first.

## Restore

Click **Restore official appearance**. You can also quit Codex and launch it normally from the Dock; session-only styling will not be present.

## Data removal

Preferences and copied images live in:

```text
~/Library/Application Support/Codex Canvas
```

Quit Codex Canvas before deleting that directory. Removing the app itself does not delete this data automatically.
