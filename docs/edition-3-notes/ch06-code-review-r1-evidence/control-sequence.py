from pathlib import Path
import hashlib,json,subprocess,sys,difflib
E=Path(__file__).resolve().parent;S=Path((E/'scratch-path.txt').read_text().strip());M=S/'public-checks';D=S/'control-source';P=D/'internal/engine/stream.go';corrected=(S/'source/internal/engine/stream.go').read_bytes();original=(S/'original/internal/engine/stream.go').read_bytes();expected=json.loads((E/'source-hashes.json').read_text());results={'source_paths':{'corrected':str(S/'source'),'original':str(S/'original'),'control':str(D)},'correction_sha256':hashlib.sha256(corrected).hexdigest(),'original_sha256':hashlib.sha256(original).hexdigest(),'checks':[]}
def sha(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def run(label,argv,want=0):
 n=subprocess.call([sys.executable,str(E/'run.py'),label,str(M),*map(str,argv)]);results['checks'].append({'label':label,'exit':n,'expected_exit':want});(E/'control-results.json').write_text(json.dumps(results,indent=2)+'\n');assert n==want,(label,n,want)
def source_match():
 for p,h in expected.items():assert sha(D/p)==h,p
source_match();assert corrected!=original
results['test_sources']={p.name:sha(p) for p in M.glob('*_test.go')};results['module_sources']={p.name:sha(p) for p in [M/'go.mod',M/'go.sum']}
run('independent-corrected-compile',['go','test','-race','-c','-o',S/'corrected-public.test','.']);results['corrected_binary_sha256']=sha(S/'corrected-public.test')
run('independent-corrected-all',[S/'corrected-public.test','-test.v','-test.count=1']);run('independent-vet',['go','vet','./...'])
P.write_bytes(original);assert sha(P)==results['original_sha256'];changed=[p for p,h in expected.items() if sha(D/p)!=h];assert changed==['internal/engine/stream.go'];results['changed_for_control']=changed;(E/'original-restoration-control.diff').write_text(''.join(difflib.unified_diff(corrected.decode().splitlines(True),original.decode().splitlines(True),fromfile='corrected/stream.go',tofile='original-control/stream.go')))
run('independent-original-compile',['go','test','-race','-c','-o',S/'original-public.test','.']);results['original_binary_sha256']=sha(S/'original-public.test');assert results['original_binary_sha256']!=results['corrected_binary_sha256']
run('independent-original-boundaries',[S/'original-public.test','-test.v','-test.count=1','-test.run','^(TestGeminiSignedAndOpaqueBoundaries|TestReviewerGeminiExplicitBoundariesAndReplay)$'],1)
P.write_bytes(corrected);source_match();results['restored_99_source_files_equal']=True
run('independent-restored-compile',['go','test','-race','-c','-o',S/'restored-public.test','.']);results['restored_binary_sha256']=sha(S/'restored-public.test');results['restored_binary_equals_corrected']=results['restored_binary_sha256']==results['corrected_binary_sha256'];assert results['restored_binary_equals_corrected']
run('independent-restored-all',[S/'restored-public.test','-test.v','-test.count=1']);assert results['test_sources']=={p.name:sha(p) for p in M.glob('*_test.go')};(E/'control-results.json').write_text(json.dumps(results,indent=2)+'\n');print(json.dumps(results,indent=2))
