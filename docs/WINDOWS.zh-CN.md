# Windows 预览版

**简体中文** | [English](WINDOWS.md)

Codex Canvas 1.1.0 新增实验性 Windows x64 版本。OpenAI 当前通过 Microsoft Store 分发 ChatGPT/Codex 桌面应用，因此 Codex Canvas 会发现已安装的 Appx 包并读取清单，不依赖固定安装目录。

## 运行要求

- Windows 10 或 Windows 11 x64
- 已从 Microsoft Store 安装官方 [Windows 版 ChatGPT 桌面应用](https://learn.chatgpt.com/docs/windows/windows-app)
- Microsoft Edge WebView2 Runtime

官方 ChatGPT 应用命令行安装方式：

```powershell
winget install --id 9PLM9XGG6VKS -s msstore
```

## 安装 Codex Canvas

从本仓库 GitHub Release 下载 `Codex-Canvas-1.1.0-Windows-x64.zip`，验证 SHA-256 后解压并运行 `Codex-Canvas.exe`。

预览版尚未进行 Authenticode 签名，Windows SmartScreen 可能在首次启动时警告。仅对从本项目官方 Release 页面下载的文件绕过警告。

## 实现方式

- 通过 PowerShell 发现当前用户的 ChatGPT/Codex Appx 包。
- 使用 `Get-AppxPackageManifest` 读取清单，在 Microsoft Store 更新后仍能找到可执行文件。
- 只请求窗口正常关闭，不使用 `taskkill` 强制结束进程。
- 使用随机、仅限本机回环地址的 CDP 端口启动桌面应用。
- 与 macOS 共用受控 CSS 生成器和恢复逻辑。

## 预览限制

当前二进制已完成交叉编译，并验证 PE 架构、嵌入资源、路径裁剪和 Windows 构建标签。项目暂时没有 Windows 测试机，因此 Microsoft Store 应用激活行为仍需真机确认。欢迎通过 GitHub Issues 反馈结果，请勿附带隐私提示词、对话、令牌或个人路径。
