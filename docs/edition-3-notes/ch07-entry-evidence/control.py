from pathlib import Path
import sys,difflib,hashlib,json
w=Path(__file__).resolve().parent;kind,action=sys.argv[1:]
relative,old,new={
'replay':('gui/gui.go','message["replay"] = true','message["replay"] = true\n\t\tif message["type"] == "part_final" { message["part_id"] = fmt.Sprint("replay-", message["part_id"]) }'),
'pause':('internal/agent/dispatch.go','if resume := r.owner.pauseChannel(); resume != nil {','if resume := r.owner.pauseChannel(); false && resume != nil {')
}[kind]
p=w/'positive'/relative; original=w/'controls'/(kind+'.original')
if action=='apply':
 text=p.read_text();assert text.count(old)==1
 original.write_bytes(p.read_bytes());newtext=text.replace(old,new);p.write_text(newtext)
 (w/'controls'/(kind+'.diff')).write_text(''.join(difflib.unified_diff(text.splitlines(True),newtext.splitlines(True),fromfile='positive/'+relative,tofile='mutant/'+relative)))
 (w/'controls'/(kind+'.hashes.json')).write_text(json.dumps({'file':relative,'original_sha256':hashlib.sha256(original.read_bytes()).hexdigest(),'mutant_sha256':hashlib.sha256(p.read_bytes()).hexdigest()},indent=2)+'\n')
else:
 p.write_bytes(original.read_bytes())
 assert hashlib.sha256(p.read_bytes()).hexdigest()==json.loads((w/'controls'/(kind+'.hashes.json')).read_text())['original_sha256']
 print(kind+' restored exact original SHA256')
