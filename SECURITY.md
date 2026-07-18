# Security policy

[简体中文](SECURITY.zh-CN.md) | **English**

## Supported version

Security fixes currently target the latest released version.

## Report a vulnerability

Please do not open a public issue for a suspected vulnerability. Email [asbacklight@gmail.com](mailto:asbacklight@gmail.com) with:

- affected version and macOS/Codex versions;
- reproduction steps or a proof of concept;
- expected impact;
- any suggested remediation.

Do not include real prompts, conversations, access tokens, or other private data. You should receive an acknowledgement within seven days. A coordinated disclosure timeline will be discussed after validation.

## Security boundaries

The app intentionally opens a random loopback-only CDP endpoint for a themed Codex session. This is documented behavior, but unintended exposure beyond `127.0.0.1`, arbitrary script execution, bundle modification, or data exfiltration should be reported as security issues.
