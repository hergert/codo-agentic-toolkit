#!/usr/bin/env python3
# Append deduped "Experience:" bullets to docs/experience/ledger.md (no external deps)
import sys, os, re, hashlib, datetime, json

E = json.load(sys.stdin)
msg = (E.get("message") or "").strip()

# Stop/SubagentStop events do not include a message field; pull from transcript if provided
if (not msg) and E.get("transcript_path"):
    try:
        with open(E["transcript_path"], "r", encoding="utf-8") as tf:
            entries = [json.loads(line) for line in tf if line.strip()]
        for entry in reversed(entries):
            if entry.get("role") != "assistant":
                continue
            content = entry.get("content")
            if isinstance(content, list):
                chunks = [block.get("text", "") for block in content if block.get("type") == "text"]
                msg = "\n".join(chunks).strip()
            elif isinstance(content, str):
                msg = content.strip()
            if msg:
                break
    except Exception:
        msg = ""

if not msg or "Experience:" not in msg:
    sys.exit(0)

lines = [ln.strip() for ln in msg.splitlines() if ln.strip()]
ix = next((i for i,l in enumerate(lines) if l.lower().startswith("experience")), None)
if ix is None: sys.exit(0)
payload = lines[ix+1:]

def norm(t): return re.sub(r"\s+"," ",t.lower().strip())

# Lightweight SimHash64 for near-dup detection (3-grams)
def simhash64(text):
    toks = re.findall(r"[a-z0-9_/.-]+", text.lower())
    grams = [" ".join(toks[i:i+3]) for i in range(len(toks)-2)] or toks
    v=[0]*64
    for g in grams:
        h=int(hashlib.sha1(g.encode()).hexdigest()[:16],16)
        for b in range(64): v[b]+=1 if (h>>b)&1 else -1
    out=0
    for b,s in enumerate(v):
        if s>0: out|=(1<<b)
    return out

def ham(a,b): return (a^b).bit_count()

# Redact sensitive information before saving
def redact(s: str) -> str:
    s = re.sub(r'(?i)bearer\s+[A-Za-z0-9._-]+', 'bearer ***', s)
    s = re.sub(r'AKIA[0-9A-Z]{16}', 'AKIA****************', s)  # AWS key
    s = re.sub(r'[A-Za-z0-9_-]{20,}\.[A-Za-z0-9_-]{20,}\.[A-Za-z0-9_-]{10,}', '***.***.***', s)  # JWT-ish
    s = re.sub(r'-----BEGIN [^-]+ PRIVATE KEY-----.*?-----END [^-]+ PRIVATE KEY-----', '-----BEGIN *** PRIVATE KEY-----\n...\n-----END *** PRIVATE KEY-----', s, flags=re.DOTALL)
    return s

os.makedirs("docs/experience", exist_ok=True)
ledger = "docs/experience/ledger.md"
indexp = "docs/experience/index.json"

existing=[]
if os.path.exists(ledger):
    with open(ledger,"r",encoding="utf-8") as f:
        existing=[ln.strip() for ln in f if ln.startswith("- [")]

cache=[]
for ln in existing:
    body = re.sub(r"^- \[[^\]]+\]\s*", "", ln).split("<!--")[0].strip()
    body = re.sub(r"\s*tags=\[.*?\]\s*", "", body).strip()
    n = norm(body)
    sim = simhash64(n)
    cache.append((n, sim, ln))

ts = datetime.datetime.utcnow().isoformat(timespec="seconds")+"Z"
cands=[]
for raw in payload:
    if not raw or raw.startswith("#"): continue
    if raw.lower().startswith("experience"): continue
    if raw.startswith("-"): raw = raw[1:].strip()
    if raw.lower().startswith("nothing learned"):
        cands.append(("nothing learned", [])); continue
    m_tags = re.search(r"tags\s*=\s*\[([^\]]*)\]", raw)
    tags = [t.strip().strip('"\'') for t in (m_tags.group(1).split(","))] if m_tags else []
    text = re.sub(r"\s*tags=\[.*?\]\s*", "", raw).strip()
    text = redact(text)  # Apply redaction before saving
    if text: cands.append((text, tags))

appended=0
with open(ledger,"a",encoding="utf-8") as out:
    for text,tags in cands:
        n = norm(text)
        if not n: continue
        s = simhash64(n)
        if any(ham(s, es) <= 3 for en,es,_ in cache):  # near-dup
            continue
        h = hashlib.sha1(n.encode()).hexdigest()[:10]
        tagstr = f'  tags={tags}' if tags else ''
        out.write(f"- [{ts}] {text}{tagstr}  <!-- {h} -->\n")
        cache.append((n,s,text)); appended+=1

# keep a tiny tag index for quick lookups
idx = {}
if os.path.exists(indexp):
    try:
        with open(indexp, "r", encoding="utf-8") as f:
            idx = json.load(f)
    except Exception:
        idx = {}
for text,tags in cands:
    n=norm(text)
    if not n or not tags: continue
    h=hashlib.sha1(n.encode()).hexdigest()[:10]
    for t in tags:
        idx.setdefault(t,{}).setdefault("ids",[])
        if h not in idx[t]["ids"]: idx[t]["ids"].append(h)
        idx[t]["last_seen"]=ts
        idx[t]["count"]=idx[t].get("count",0)+1
with open(indexp, "w", encoding="utf-8") as f:
    json.dump(idx, f, indent=2)
