from pathlib import Path
import sys,json,hashlib
p=Path(sys.argv[1]); b=p.read_bytes(); lines=b.splitlines(keepends=True)
a=int(sys.argv[2]) if len(sys.argv)>2 else 1
z=int(sys.argv[3]) if len(sys.argv)>3 else len(lines)
assert 1<=a<=z<=len(lines)
out=b"".join(lines[a-1:z])
record={"path":str(p),"start_line":a,"end_line":z,"total_lines":len(lines),"file_bytes":len(b),"file_sha256":hashlib.sha256(b).hexdigest(),"output_bytes":len(out),"output_sha256":hashlib.sha256(out).hexdigest(),"argv":sys.argv,"exit":0}
with (Path(__file__).parent/"reads.jsonl").open("a") as f: f.write(json.dumps(record)+"\n")
print(f"--- {p}: lines {a}-{z} of {len(lines)} ---",flush=True)
sys.stdout.buffer.write(out)
