---
description: Distill the session's learning (no writes). Run built-in /compact separately if needed.
argument-hint: [feature-key]
allowed-tools: Read
---

Aim for 2-4 bullets max; each bullet ≤25 words; prefer action verbs and concrete heuristics.

# Experience distiller (no writes)

Summarize this session as ≤120 tokens. Then propose **0-5 bullets** worth remembering long-term, or say "nothing learned".
Each bullet must be:

- Atomic, reusable, and secret-free
- Tagged with 3-5 tags (tech/area/component/risk)
- Optional scope hint: paths=["pkg/x/","cmd/api/"]

# Output (exact)

Experience:

- <bullet 1> tags=[...] scope=[...]
- <bullet 2> tags=[...]
- ...
  (or: Experience: nothing learned)
