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
                for msg in _candidate_messages(ev):
                    if msg.get("role") != "user":
                        continue
                    text = _coerce_text(msg.get("content"))
                    if text:
                        last_text = text
    except Exception:
        pass
    return last_text.strip()


def _candidate_messages(ev):
    """Yield message-like dicts regardless of event shape."""
    if not isinstance(ev, dict):
        return []
    candidates = []
    msg = ev.get("message")
    if isinstance(msg, dict):
        candidates.append(msg)
    if isinstance(ev.get("messages"), list):
        candidates.extend(m for m in ev["messages"] if isinstance(m, dict))
    if isinstance(ev.get("data"), dict):
        data = ev["data"]
        if isinstance(data.get("message"), dict):
            candidates.append(data["message"])
        if isinstance(data.get("messages"), list):
            candidates.extend(m for m in data["messages"] if isinstance(m, dict))
    top = {k: ev.get(k) for k in ("role", "content") if k in ev}
    if top and "role" in top:
        candidates.append(top)
    return candidates


def _coerce_text(content):
    if isinstance(content, str):
        return content
    if isinstance(content, dict):
        text = content.get("text") or content.get("content")
        if isinstance(text, str):
            return text
        if isinstance(text, list):
            return _coerce_text(text)
        if content.get("type") == "text" and isinstance(content.get("text"), str):
            return content.get("text")
    if isinstance(content, list):
        pieces = []
        for item in content:
            if isinstance(item, dict):
                if item.get("type") == "text" and isinstance(item.get("text"), str):
                    pieces.append(item.get("text", ""))
                elif "text" in item and isinstance(item["text"], str):
                    pieces.append(item["text"])
                elif "content" in item and isinstance(item["content"], str):
                    pieces.append(item["content"])
            elif isinstance(item, str):
                pieces.append(item)
        return "\n".join(pieces).strip()
    return ""


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

    try:
        import subprocess

        out = subprocess.check_output(
            ["bash", "-lc", "git diff --name-only | wc -l"], text=True
        ).strip()
        n = int(out or "0")
        if n > 8 and msg.startswith("Next: /execute"):
            print(
                json.dumps(
                    {
                        "systemMessage": "This is a large diff; consider `reader <key>` (scout) before /execute, or `reviewer <key>` as a final audit."
                    }
                )
            )
            sys.exit(0)
    except Exception:
        pass

    out = {"systemMessage": msg}
    print(json.dumps(out))
    sys.exit(0)


if __name__ == "__main__":
    main()
