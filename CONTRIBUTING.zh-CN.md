# 参与贡献

**简体中文** | [English](CONTRIBUTING.md)

感谢帮助改进 Codex Canvas。

## 权限与协作

任何人都可以 Fork 这个公开仓库、在自己的 Fork 中创建分支并提交 Pull Request。提交 Pull Request 不会自动获得本仓库的直接推送或合并权限。只有维护者及你明确邀请、并授予相应 GitHub 角色的协作者可以合并代码。`main` 分支会按照仓库保护规则要求 CI 和代码审核。

## 提交修改前

- 先搜索已有 Issue。
- 保持改动聚焦，避免夹带无关格式化。
- 不要加入远程脚本、遥测、任意主题代码执行或修改 Codex 应用包的逻辑。
- 不要提交证书、签名身份、令牌、本地应用数据或生成的发布二进制。

## 开发与测试

```bash
cd frontend
npm ci
npm run build
cd ..
go test ./...
```

使用 `wails dev` 进行交互开发，使用 `./scripts/build-macos.sh` 构建 Universal 应用。

## Pull Request

请说明用户可见变化、安全影响、兼容性假设和验证方法。界面改动请提供前后截图；注入路径改动必须包含测试，并说明如何继续保证只接受受控数据。

提交贡献即表示你同意按 MIT License 授权该贡献。
