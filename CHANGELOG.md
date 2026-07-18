# Changelog / 更新日志

All notable changes are documented here. / 本文件记录重要版本变化。

## 1.1.0 - 2026-07-18

### English

- Added an experimental Windows 10/11 x64 build with Microsoft Store Appx discovery and graceful process handling.
- Split platform lifecycle code into isolated macOS and Windows adapters.
- Added Windows build scripts, CI cross-compilation, and bilingual Windows documentation.

### 中文

- 增加实验性 Windows 10/11 x64 版本，支持 Microsoft Store Appx 发现和优雅进程处理。
- 将平台生命周期代码拆分为独立的 macOS 与 Windows 适配层。
- 增加 Windows 构建脚本、CI 交叉编译和中英文 Windows 文档。

## 1.0.0 - 2026-07-18

### English

- Initial macOS release with a Universal 2 binary for Intel and Apple Silicon.
- Added local image selection and five original built-in backgrounds.
- Added Chinese/English UI, system-language detection, and persistent manual switching.
- Added live preview and readability controls.
- Added reversible, session-scoped CDP theme injection without modifying the Codex bundle.
- Added high-contrast text and input tokens plus one-click official appearance restore.

### 中文

- 首个 macOS 版本，提供 Intel 与 Apple Silicon 通用的 Universal 2 程序。
- 支持本地图片和五款原创内置背景。
- 支持中英文、系统语言识别与持久化手动切换。
- 提供实时预览和阅读效果控制。
- 通过可逆、仅限会话的 CDP 主题注入实现背景，不修改 Codex 应用包。
- 增加高对比文字与输入框令牌，并支持一键恢复官方外观。
