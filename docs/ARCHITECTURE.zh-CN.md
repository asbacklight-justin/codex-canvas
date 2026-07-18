# 架构与安全模型

**简体中文** | [English](ARCHITECTURE.md)

## 模块边界

```text
Wails UI
  ├─ ThemeConfig（跨平台）
  ├─ 实时预览（跨平台）
  └─ 类型化 App bindings
       └─ CodexService
            ├─ 图片验证与本地持久化（跨平台）
            ├─ 受控 CSS 生成器（跨平台）
            ├─ CDP 页面发现与 WebSocket 注入（跨平台）
            └─ Codex 生命周期适配（当前为 macOS）
```

主题配置是数据而不是可执行代码。Codex Canvas 不接收任意 CSS、JavaScript、远程 URL 或第三方主题包。选择图片后，应用会先验证并复制到本地数据目录，生成样式时再转为本地 `data:` URL。

## macOS 生命周期

1. 使用 bundle identifier `com.openai.codex` 校验候选应用，兼容 `Codex.app` 和显示名为 `ChatGPT.app` 的版本。
2. 通过 Apple Event 请求 Codex 正常退出，不强制杀进程。
3. 由操作系统在 `127.0.0.1` 分配随机空闲 TCP 端口。
4. 使用 `/usr/bin/open` 和仅限回环地址的 CDP 参数启动新实例。
5. 客户端轮询本机 `/json` 端点，连接符合条件的页面，并同时覆盖当前页面与后续重载。
6. 恢复时优先实时移除样式；端点不可用时按官方方式重新启动 Codex。

## 信任与威胁模型

- **应用完整性：**不修改 Codex 应用包内任何文件，因此不破坏其代码签名。
- **本地数据：**设置和所选背景保存在 `~/Library/Application Support/Codex Canvas`，目录尽可能使用仅当前用户可访问的权限。
- **网络：**主题操作只访问随机的本机回环端点。项目没有 Codex Canvas 服务端，也不会上传应用数据。
- **CDP 暴露：**主题会话运行期间，同一 macOS 用户权限下的其他进程理论上可能发现并访问本机调试端点。敏感会话请恢复官方外观或正常启动 Codex。
- **兼容性：**Codex 内部 DOM 和设计令牌可能变化。注入只影响当前会话，失败时不会修改官方应用。

## Universal 2 构建

`scripts/build-macos.sh` 使用本机 macOS SDK 编译两次 Go/Wails 应用：

- `GOARCH=amd64`，Clang 目标为 `x86_64`
- `GOARCH=arm64`，Clang 目标为 `arm64`

随后通过 `lipo` 合并两个 Mach-O 程序。本地测试包采用 ad-hoc 签名；正式公开分发需要 Developer ID Application 签名和 Apple 公证。

## Windows 路线

1. 发现 Microsoft Store、`%LOCALAPPDATA%` 和 `%ProgramFiles%` 安装，并验证可执行文件信息。
2. 优先请求优雅关闭；任何强制终止回退都必须先获得用户明确同意。
3. 使用等价的仅限回环地址 CDP 参数和随机端口启动。
4. 继续使用 `os.UserConfigDir()`，在 Windows 中自然映射到 `%AppData%\Codex Canvas`。
5. 构建 NSIS/MSI 安装器并应用 Authenticode 签名。
6. 在 Windows 11 x64、Store 与独立安装版 Codex 上完成测试。

## 发布门槛

- Go 测试、前端生产构建和 Universal 架构检查通过。
- 正式版通过 macOS Developer ID 签名和公证。
- 发布 SHA-256 校验值和更新说明。
- 在 Intel 与 Apple Silicon 真机完成应用、页面切换、恢复、正常启动和升级兼容测试。
