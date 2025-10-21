---
description: Distill the session's learning; write to ledger only if truly useful.
argument-hint: [feature-key]
allowed-tools: Read, Edit
---

Aim for 1–3 bullets max; each bullet ≤25 words; concrete, secret-free, reusable.

# Produce a short "Experience" block:
- If there's **nothing worth keeping**, print exactly: `Experience: nothing learned` and **do not edit any file**.
- Otherwise print:

Experience:
- <bullet 1>  tags=[area,tech,risk]  scope=[optional/paths]
- <bullet 2>  tags=[...]

# Append policy:
- Only if there is at least **one** bullet, append to `docs/experience/ledger.md` in this format:
  `- [YYYY-MM-DDTHH:MM:SSZ] <bullet text>  tags=[...]`
- Do **not** duplicate bullets already present (do a quick text search).
- Never include secrets, keys, tokens, or URLs with credentials.
