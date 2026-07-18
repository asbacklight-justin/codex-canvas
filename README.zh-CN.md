# Codex Canvas · 幕光

**简体中文** | [English](README.md)

Codex Canvas 是一个面向 macOS 和 Windows、本地优先且可逆的 Codex Desktop 背景与主题管理器。你可以选择本地图片或五款内置原创背景，调整阅读效果、实时预览，并通过临时本地主题通道重新启动 Codex。

**版本 1.1.0** · 联系邮箱：[asbacklight@gmail.com](mailto:asbacklight@gmail.com)

> Codex Canvas 是独立开源项目，与 OpenAI 无隶属、背书或赞助关系。Codex 和 OpenAI 是其各自权利人的商标。

<p align="center">
  <img src="docs/images/codex-canvas-overview.png" alt="Codex Canvas 内置背景、实时预览和阅读效果控制界面" width="100%">
</p>

## 主要功能

- 原生 macOS 图片选择器，支持 PNG、JPEG、WebP、GIF，最大 16 MB
- 内置科技、原创二次元、城市、自然、动物五款 16:9 背景
- 中英文界面：默认跟随系统，识别失败回退英文，也可手动切换并记住选择
- Codex 风格实时预览
- 暗色遮罩、内容/侧栏透明度、模糊、缩放、位置、圆角和强调色控制
- 深浅色 Codex 均保持高对比文字与输入框可读性
- 自动发现 `/Applications/Codex.app` 和兼容的 `/Applications/ChatGPT.app`
- 一键应用与恢复，不修改 `Codex.app`、`app.asar`、应用资源或代码签名
- Universal 2 安装包，Intel 与 Apple Silicon Mac 均可原生运行
- 实验性 Windows 10/11 x64 版本，支持动态发现 Microsoft Store Appx 安装
- 不上传图片，不含分析、遥测、远程主题包，也不执行用户提供的任意 CSS/JavaScript

## 下载与安装

从 [GitHub Releases](../../releases/latest) 下载可直接使用的成品：

- macOS：`Codex-Canvas-1.1.0-macOS-Universal.zip`，支持 Intel 和 Apple Silicon
- Windows：`Codex-Canvas-1.1.0-Windows-x64.zip`，支持 Windows 10/11 x64（预览版）

macOS 解压后把 **Codex Canvas.app** 移入 `/Applications`；Windows 解压后运行 `Codex-Canvas.exe`。

当前社区构建采用 ad-hoc 签名，尚未经过 Apple 公证。首次启动时 macOS 可能要求使用 **按住 Control 点击 → 打开**。详细步骤和校验方式见[安装说明](docs/INSTALLATION.zh-CN.md)。

macOS 运行要求：

- macOS 12 或更高版本
- Intel（`x86_64`）或 Apple Silicon（`arm64`）
- Codex Desktop 已安装到 `/Applications`

Windows 需要已安装 Microsoft Store 官方 ChatGPT 桌面应用和 WebView2，详见 [Windows 预览版说明](docs/WINDOWS.zh-CN.md)。首个 Windows 包为交叉编译版本，等待社区真机确认。

## 工作原理

Codex Canvas 不会给 Codex 应用打补丁。应用主题时，它会：

1. 请求 Codex 正常退出；
2. 在 `127.0.0.1` 随机分配仅限本机的端口；
3. 使用临时 Chromium DevTools Protocol 端点重新启动 Codex；
4. 向 Codex 页面注入由受控参数生成的样式元素；
5. 仅在本次 Codex 会话中保持主题。

从 Dock 正常启动 Codex 即可回到官方外观，也可在 Codex Canvas 中点击“恢复官方外观”。使用前建议阅读[架构与安全模型](docs/ARCHITECTURE.zh-CN.md)和[隐私说明](docs/PRIVACY.zh-CN.md)。

## 从源码构建

需要 Go 1.23+、Node.js 20+、Wails 2.11+ 和 macOS Command Line Tools。

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.11.0
cd frontend
npm ci
npm run build
cd ..
go test ./...
./scripts/build-macos.sh
```

构建脚本会分别编译 `x86_64` 和 `arm64`，再通过 `lipo` 合并为 Universal 2 应用，输出到 `build/bin/Codex Canvas.app`。

开发模式：

```bash
cd frontend && npm ci && cd ..
wails dev
```

## 文档

- [安装说明](docs/INSTALLATION.zh-CN.md) · [Installation](docs/INSTALLATION.md)
- [架构说明](docs/ARCHITECTURE.zh-CN.md) · [Architecture](docs/ARCHITECTURE.md)
- [隐私说明](docs/PRIVACY.zh-CN.md) · [Privacy](docs/PRIVACY.md)
- [Windows 预览版](docs/WINDOWS.zh-CN.md) · [Windows preview](docs/WINDOWS.md)
- [参与贡献](CONTRIBUTING.zh-CN.md) · [Contributing](CONTRIBUTING.md)
- [安全政策](SECURITY.zh-CN.md) · [Security policy](SECURITY.md)
- [更新日志](CHANGELOG.md)

## Windows 状态

Windows 安装发现、优雅退出、CDP 启动和 x64 打包已经实现。剩余发布门槛是真机兼容测试、Authenticode 签名和安装器打包，详见 [Windows 预览版说明](docs/WINDOWS.zh-CN.md)与[架构说明](docs/ARCHITECTURE.zh-CN.md#windows-状态)。欢迎贡献代码和真机反馈。

## 许可证

除非文件另有说明，Codex Canvas 及其内置原创背景均采用 [MIT License](LICENSE) 开源。
