#!/usr/bin/env python3
# Minimal end-of-turn nudge: runs on Stop/SubagentStop.
# Reads the transcript JSONL, looks at the last user prompt,
# and prints a small systemMessage with the next suggested step.
import json, sys, os


def read_last_user_text(transcript_path: str) -> str:
    if not transcript_path or not os.path.exists(transcript_path):
        return ""
    last_text = ""
    try:
        with open(transcript_path, "r", encoding="utf-8") as f:
            for line in f:
                try:
                    ev = json.loads(line)
                except Exception:
                    continue
                role = (ev.get("message") or {}).get("role")
                content = (ev.get("message") or {}).get("content")
                if role == "user":
                    # content can be str or array-of-blocks
                    if isinstance(content, str):
                        last_text = content
                    elif isinstance(content, list):
                        pieces = []
                        for b in content:
                            if isinstance(b, dict) and b.get("type") == "text":
                                pieces.append(b.get("text", ""))
                        if pieces:
                            last_text = "\n".join(pieces)
    except Exception:
        pass
    return last_text.strip()


def hint_for(text: str):
    t = (text or "").strip()
    if not t:
        return None
    first = t.splitlines()[0].strip()
    if first.startswith("/prime"):
        return 'Next: /plan "<key>" to map & plan, or /execute "<key>" if a plan exists.'
    if first.startswith("/plan"):
        return 'Next: /execute "<key>" — tests first, smallest viable diff.'
    if first.startswith("/execute"):
        return 'Next: /review "<key>" — aim APPROVE; then /prepare-commit.'
    if first.startswith("/review"):
        return 'Next: /prepare-commit — human reviews & commits.'
    if first.startswith("/prepare-commit"):
        return "Reminder: human commits; keep diffs tight and messages clear."
    return None


def main():
    try:
        E = json.load(sys.stdin)
    except Exception:
        sys.exit(0)

    transcript = E.get("transcript_path") or ""
    last_user = read_last_user_text(transcript)
    msg = hint_for(last_user)
    if not msg:
        sys.exit(0)

    out = {"systemMessage": msg}
    print(json.dumps(out))
    sys.exit(0)


if __name__ == "__main__":
    main()
