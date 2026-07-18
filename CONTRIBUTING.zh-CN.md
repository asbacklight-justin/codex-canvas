# 参与贡献

**简体中文** | [English](CONTRIBUTING.md)

感谢帮助改进 Codex Canvas。

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
