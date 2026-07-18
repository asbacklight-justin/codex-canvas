# Contributing

[简体中文](CONTRIBUTING.zh-CN.md) | **English**

Thanks for helping improve Codex Canvas.

## Access and collaboration

Anyone can fork this public repository, create a branch in their fork, and open a Pull Request. Opening a Pull Request does not grant direct push or merge access to this repository. Only maintainers and explicitly invited collaborators with the necessary GitHub role can merge changes. The `main` branch requires CI and review according to the repository protection rules.

## Before opening a change

- Search existing issues first.
- Keep changes focused and avoid unrelated formatting rewrites.
- Do not add remote scripts, telemetry, arbitrary theme execution, or modifications to the Codex application bundle.
- Do not commit certificates, signing identities, tokens, local application data, or generated release binaries.

## Development

```bash
cd frontend
npm ci
npm run build
cd ..
go test ./...
```

Run `wails dev` for interactive development. Build the Universal app with `./scripts/build-macos.sh`.

## Pull requests

Describe the user-visible change, security impact, compatibility assumptions, and validation performed. UI changes should include before/after screenshots. Changes to the injection path should include tests and an explanation of why the controlled-data boundary remains intact.

By contributing, you agree that your contribution is licensed under the MIT License.
