# 安装说明

**简体中文** | [English](INSTALLATION.md)

## 运行要求

- macOS 12 或更高版本
- Intel 或 Apple Silicon Mac
- Codex Desktop 已安装在 `/Applications`

## 安装发布版本

1. 从 [GitHub Releases](../../../releases/latest) 下载 `Codex-Canvas-1.1.0-macOS-Universal.zip`。
2. 对照发布说明中的校验值验证 SHA-256。
3. 解压并把 `Codex Canvas.app` 移入 `/Applications`。
4. 如果首次启动出现“无法验证开发者”，请按住 Control 点击应用并选择“打开”。

当前社区版本使用 ad-hoc 签名，尚未经过 Apple 公证。不要对来源并非本仓库 Releases 页面的安装包绕过系统安全警告。

## 应用背景

1. 打开 Codex Canvas。
2. 在左侧选择内置背景，或在右侧选择本地图片。
3. 调整阅读参数并查看预览。
4. 点击“应用到 Codex”并确认重新启动。

Codex Canvas 会先请求 Codex 退出，请提前发送或保存尚未完成的输入内容。

## 恢复官方外观

点击“恢复官方外观”。你也可以退出 Codex 后从 Dock 正常启动；仅限会话的主题样式不会继续存在。

## 删除本地数据

偏好设置和复制的图片位于：

```text
~/Library/Application Support/Codex Canvas
```

删除该目录前请退出 Codex Canvas。仅删除应用本身不会自动删除这些数据。
