from pathlib import Path
import subprocess,sys,json
E=Path(__file__).resolve().parent;R=E.parents[2]
def run(label,args,expected):
 n=subprocess.call([sys.executable,str(E/'run.py'),label,str(R),*args]);assert n==expected,(label,n,expected)
run('actual-gemini-shape',[sys.executable,str(E/'shape_fixture.py'),'/tmp/e3-ch06-correction-r1/shape'],0)
run('actual-cli-gated-and-stricter-assertion',[sys.executable,str(E/'cli-gated-control.py')],1)
d=json.loads((E/'gated-cli.json').read_text());assert d['candidate_exit']==0 and d['controller_error'] is None;assert all(v for k,v in d['checks'].items() if k!='hint_in_answer');assert d['checks']['hint_in_answer'] is False
(E/'integration-results.json').write_text(json.dumps({'actual_gemini_shape_exit':0,'gated_cli_candidate_exit':0,'gated_cli_public_checks_pass':True,'stricter_hint_answer_check':False,'stricter_check_command_exit':1,'no_live_model_calls':True},indent=2)+'\n');print('All planned local integrations completed; stricter assertion failure retained, not credited.')
