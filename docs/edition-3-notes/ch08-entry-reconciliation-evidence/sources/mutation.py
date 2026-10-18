from run import *
import hashlib
G=T/'real-gui';p=G/'settings.go';original=p.read_bytes();needle=b'err = os.Rename(path+".tmp", path)';assert original.count(needle)==1
negative=original.replace(needle,b'err = os.Remove(path+".tmp") // OWNED AUDIT MUTANT: discard durable update')
(E/'mutation-sources').mkdir()
(E/'mutation-sources/settings-old.go.txt').write_bytes(original)
(E/'mutation-sources/settings-negative.go.txt').write_bytes(negative)
(E/'mutation-sources/go.mod').write_bytes((G/'go.mod').read_bytes())
records=[]
try:
 for stage,source in [('old',original),('negative',negative),('restored',original)]:
  p.write_bytes(source)
  rec={'stage':stage,'source_path':str(p),'source_sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'source_bytes':len(source)}
  argv=[str(T/'bin/driver'),'-mode','settings','-gui',str(G),'-gui-model','claude-sonnet-4-6']
  rc=run('settings-'+stage,argv,T/'work',90)
  rec['driver_exit']=rc
  try:
   report=json.loads((E/('settings-'+stage+'.stdout')).read_text());rec.update(score=report['score'],max_score=report['max_score'],passed=report['passed'],failed_check_ids=[c['id'] for c in report['checks'] if not c['passed']])
  except Exception as ex:rec['report_error']=str(ex)
  records.append(rec);(E/'mutation-results.json').write_text(json.dumps(records,indent=2)+'\n')
  expected=1 if stage=='negative' else 0
  if rc!=expected:raise RuntimeError('unexpected stage exit; stopping')
finally:
 p.write_bytes(original)
 (E/'mutation-sources/settings-restored.go.txt').write_bytes(p.read_bytes())
 (E/'mutation-restore.json').write_text(json.dumps({'scratch_source_restored':p.read_bytes()==original,'restored_sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'repository_source_sha256':hashlib.sha256((R/'solutions/ch08/gui/settings.go').read_bytes()).hexdigest()},indent=2)+'\n')
