---
description: Diagnose failing tests and propose minimal fixes
allowed-tools: Edit, Bash(go test ./...), Bash(go test -run*), Bash(pnpm test*), Bash(npm test*), Bash(yarn test*), Bash(bun test*), Bash(vitest*), Bash(uv run pytest*), Bash(pytest*), Bash(mvn test*), Bash(gradle test*), Bash(flutter test*)
---
1) Detect runner from repo (go, JS/TS, Python, JVM, Flutter) and run tests quietly; capture failures.
2) Propose the smallest changes to fix; prefer test updates only when spec dictates.
3) Re-run tests until green; produce a concise rationale for each change.
