#!/usr/bin/env python3
"""Package only sanitized textual evidence, preserving exact source bytes and exclusions."""
import hashlib,json,pathlib,re,shutil
source=pathlib.Path('/tmp/cadence-ch05-evidence-r1')
notes=pathlib.Path('/Users/bill/projects/ensemble/docs/edition-3-notes')
target=notes/'ch05-student-evidence';report=notes/'ch05-student-review.md'
def sha(raw):return hashlib.sha256(raw).hexdigest()
entries={};excluded={};payloads={}
token=re.compile(rb'(?:sk-ant-[A-Za-z0-9_-]{20,}|sk-proj-[A-Za-z0-9_-]{20,}|AIza[0-9A-Za-z_-]{30,})')
media={'.png','.jpg','.jpeg','.gif','.pdf','.mp4','.mp3','.wav','.tar','.gz','.zip','.pyc','.bin'}
for path in sorted(source.rglob('*')):
    if not path.is_file():continue
    name=str(path.relative_to(source));raw=path.read_bytes();meta={'sha256':sha(raw),'bytes':len(raw)}
    reason=None
    if '__pycache__' in path.parts:reason='generated interpreter cache, not acceptance evidence'
    elif path.suffix.lower() in media:reason='binary/media/archive retained outside Git'
    elif b'\0' in raw:reason='binary content retained outside Git'
    else:
        try:raw.decode('utf-8')
        except UnicodeDecodeError:reason='non-UTF-8 binary retained outside Git'
    if reason:
        excluded[name]={**meta,'local_path':str(path),'reason':reason};continue
    if token.search(raw):raise SystemExit('suspected token-shaped secret; packaging stopped: '+name)
    entries[name]=meta;payloads[name]=raw
# Live histories are already filtered by the tested capture path. Reject actual residual payloads.
checked=0
def inspect(value,name):
    if isinstance(value,dict):
        if value.get('type') in ['thinking','opaque'] and set(value)-{'type','omitted'}:raise SystemExit('live provider payload: '+name)
        for key,item in value.items():
            if key.lower() in ['signature','thoughtsignature','thought_signature','thinking','api_key','apikey'] and item:raise SystemExit('live sensitive field: '+name)
            inspect(item,name)
    elif isinstance(value,list):
        for item in value:inspect(item,name)
for name,raw in payloads.items():
    if name.startswith(('live-r1/','live-remaining-r1/')) and name.endswith('.safe'):
        for line in raw.decode().splitlines():
            if line.strip():inspect(json.loads(line),name)
        checked+=1
if token.search(report.read_bytes()):raise SystemExit('suspected token in report')
# Refuse unrelated existing material; never replace a different file silently.
if target.exists():
    existing={str(f.relative_to(target)) for f in target.rglob('*') if f.is_file()}
    if existing-set(entries):raise SystemExit('unexpected existing files in portable destination')
    for name in existing:
        if (target/name).read_bytes()!=payloads[name]:raise SystemExit('different existing portable file: '+name)
target.mkdir(parents=True,exist_ok=True)
for name,raw in payloads.items():
    path=target/name;path.parent.mkdir(parents=True,exist_ok=True);path.write_bytes(raw)
state=json.loads((source/'final-candidate-state.json').read_text())
index={'schema':1,'portable_directory':target.name,'text_files':entries,'excluded_artifacts':excluded,'report':{'relative_path':report.name,'sha256':sha(report.read_bytes()),'bytes':report.stat().st_size},'candidate':state,'privacy_checks':{'token_shaped_matches':0,'live_filtered_histories_checked':checked,'residual_live_sensitive_payloads':0,'real_credential_reads_for_packaging':0},'evidence_credit':{'mocked_route_control':'uncredited; preserved as attempt','actual_route_rehearsal':'local-remaining-r4','original_human_launcher':'failed; never rewritten as PASS','human_proof':'independently reconciled from actual immutable artifacts','final_remaining_live':'passed; separate authorized invocation','earlier_disk_mutations':'five inconclusive; do not credit as detections','final_mutations':'six actual test assertion detections'},'relocation':'verify-portable.py checks relative text paths and adjacent report only; original absolute invocation paths remain historical receipts; binaries/archives/media stay outside Git'}
indexpath=notes/'ch05-student-evidence-manifest.json'
if indexpath.exists():raise SystemExit('refusing to overwrite existing evidence index')
indexpath.write_text(json.dumps(index,indent=2)+'\n')
print(json.dumps({'copied_text_files':len(entries),'excluded_files':len(excluded),'text_bytes':sum(x['bytes'] for x in entries.values()),'live_filtered_histories_checked':checked,'index_path':str(indexpath),'index_sha256':sha(indexpath.read_bytes()),'report_sha256':sha(report.read_bytes())},indent=2))
