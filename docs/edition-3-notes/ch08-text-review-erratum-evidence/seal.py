from pathlib import Path
import hashlib,json,datetime
E=Path('docs/edition-3-notes/ch08-text-review-erratum-evidence')
R=Path('docs/edition-3-notes/ch08-text-review-erratum.md')
I=E/'inventory.json'
assert not I.exists()
def ident(p):
 b=p.read_bytes();return {'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest()}
def save(name,obj):
 p=E/name;assert not p.exists();p.write_text(json.dumps(obj,indent=2)+'\n')
expected={
'docs/edition-3-notes/ch08-code-review.md':(16914,'fef65092b5bb0226386240693777b39a69dabf3d9b16aeadc02d3097421b9d1e'),
'docs/edition-3-notes/ch08-code-review-r1.md':(9543,'49b863d171f28086c0cd8ec99a40d9e5423fa72d9f90cf7f65029a091abdfda0'),
'docs/edition-3-notes/ch08-text-review.md':(11066,'10e7ea7b6b49ad163f4d84aa437909a72448aa75ad9c41c7a3ec75225edc4436'),
'docs/edition-3-notes/ch08-text-review-evidence/inventory.json':(4525,'eb6a699208d2fd82626d8f8a448d8e9539f86149648903809d3b10558f020fca'),
'docs/edition-3-notes/ch08-text-review-evidence/final-submission.json':(1142,'6af468942ff41126fc8dfc2cd52857d331eafe03aeec3fe4a7e976b3b2d843c6')}
pins=[]
for path,(size,digest) in expected.items():
 got=ident(Path(path));assert got=={'bytes':size,'sha256':digest},path
 pins.append({'path':path,**got,'matches_pre_erratum_read':True})
save('pins.json',{'scope':'Two cited source reports, original sealed text report/inventory/submission. No original evidence member or manuscript was modified.','files':pins})
requests=[('initial-review-excerpts.txt','docs/edition-3-notes/ch08-code-review.md',[(3,5),(33,51)]),('r1-review-excerpts.txt','docs/edition-3-notes/ch08-code-review-r1.md',[(1,27)]),('original-report-excerpt.txt','docs/edition-3-notes/ch08-text-review.md',[(37,37)])]
excerpts=[]
for name,path,ranges in requests:
 lines=Path(path).read_text().splitlines(keepends=True)
 text=''.join(f'--- {path}, lines {a}-{z} ---\n'+''.join(f'{n}: {lines[n-1]}' for n in range(a,z+1)) for a,z in ranges)
 p=E/name;assert not p.exists();p.write_text(text)
 excerpts.append({'source':path,'ranges':ranges,'output':name,**ident(p)})
save('reconciliation-receipt.json',{'request':'Root detected false attribution in sealed report; narrow erratum only; formal review count remains one.','source_excerpts':excerpts,'reads_performed':['Initial independent report lines 1-55, narrowed to cited 3-5 and 33-51.','Initial attempt emitted full follow-up report within a combined tool call; display truncated. No full-review conclusion is derived from that attempted output.','Recovery read initial report 45-55 and follow-up 1-43, completely displayed. Relevant follow-up evidence is 1-27.','Original report disputed bullet and immediate context; original report/inventory/submission identities.'],'setup_failure':'First combined output was 21,503 bytes against a 13,500-byte cap. Missing relevant passages were recovered in a separate bounded call. No omitted output is counted as a completed read.','causal_correction':'Test-fixture initialization/compatibility blocker and fixture-only repair, not demonstrated production cross-tab TTS defect or repair.','pass_effect':'Existing manuscript PASS unchanged for this correction only; not a new PASS or independent revalidation of the report.','root_detection_acknowledged':True,'formal_reviews_total':1,'new_formal_reviews':0,'chapter_reads':0,'lint_runs':0,'grade_runs':0,'browser_provider_runs':0,'children_created':0,'author_code_changes':0})
save('final-submission.json',{'type':'bounded factual erratum; not a formal text review','erratum':{'path':str(R),**ident(R)},'original_report':pins[2],'existing_manuscript_verdict':'PASS remains unchanged by this correction; no new formal verdict','required_author_revisions_from_erratum':[],'correction':'Sole initial R1 was inherited production-page fixture initialization/compatibility failure; R1 changed only that fixture, with production unchanged.','original_report_inventory_preserved':True,'formal_reviews_total':1,'new_formal_reviews':0,'inventory':str(I),'root_verification':'Pending root verification and join; no closure beyond this bounded submission claimed.'})
for row in pins:assert ident(Path(row['path']))=={'bytes':row['bytes'],'sha256':row['sha256']}
save('preseal-checks.json',{'source_report_pins_unchanged':True,'original_report_inventory_unchanged':True,'new_output_scope':[str(R),str(E)+'/'],'utf8_no_nul':True,'method':'Only metadata identity checks, exact report excerpt capture and new erratum serialization. No lint/grade/browser/provider or full chapter review.','seal_policy':'Inventory written last, self-excluded, then validated read-only. No writes afterward.'})
files=[]
for p in sorted(E.iterdir()):
 if p.is_file():
  b=p.read_bytes();b.decode('utf-8');assert b'\0' not in b
  files.append({'path':p.name,**ident(p)})
b=R.read_bytes();b.decode('utf-8');assert b'\0' not in b
manifest={'schema':'ch08-text-review-erratum-evidence-v1','created_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'purpose':'Narrow factual correction to sealed review; not a second formal review','self_excluded':'inventory.json; avoids circular self-hash','files':files,'external_outputs':[{'path':str(R),**ident(R)}]}
# Final write: all subsequent operations are read-only.
I.write_text(json.dumps(manifest,indent=2)+'\n')
assert {p.name for p in E.iterdir() if p.is_file() and p!=I}=={x['path'] for x in files}
for x in files:assert ident(E/x['path'])=={'bytes':x['bytes'],'sha256':x['sha256']}
for x in manifest['external_outputs']:assert ident(Path(x['path']))=={'bytes':x['bytes'],'sha256':x['sha256']}
for row in pins:assert ident(Path(row['path']))=={'bytes':row['bytes'],'sha256':row['sha256']}
print(json.dumps({'read_only_postseal_validation':'PASS','erratum':{'path':str(R),**ident(R)},'inventory':{'path':str(I),**ident(I)},'saved_submission':{'path':str(E/'final-submission.json'),**ident(E/'final-submission.json')},'sealed_evidence_members':len(files),'original_pins_unchanged':len(pins),'writes_after_seal':0},indent=2))
