#!/usr/bin/env python3
"""Finite LOCAL-only controls; imports real harness, never enters its paid branch."""
import argparse, concurrent.futures, hashlib, http.client, json, subprocess, time
from pathlib import Path
import launcher as h

def request(port,body,path='/v1/messages'):
    raw=json.dumps(body,separators=(',',':')).encode();c=http.client.HTTPConnection('127.0.0.1',port,timeout=90)
    try:c.request('POST',path,body=raw,headers={'Content-Type':'application/json'});r=c.getresponse();return r.status,r.read()
    finally:c.close()

def main():
    ap=argparse.ArgumentParser();ap.add_argument('--policy',required=True);ap.add_argument('--policy-sha256',required=True);ap.add_argument('--run-dir',required=True);a=ap.parse_args()
    p=h.verify(a.policy,a.policy_sha256);p.update(_policy_path=a.policy,_policy_hash=a.policy_sha256)
    directory=Path(a.run_dir).resolve();h.setup(directory)
    receipt={'scope':'LOCAL ONLY; no browser/native/key/external-provider path','argv':__import__('sys').argv,'policy_sha256':a.policy_sha256,'launcher_sha256':h.sha(h.__file__),'rehearsal_sha256':h.sha(__file__),'checks':[],'unexpected_failure':None}
    start=time.monotonic();all_forwards=[];gui=None
    def record(name,**data):receipt['checks'].append({'name':name,**data});h.save(directory/'receipt.json',receipt)
    try:
        bad=[__import__('sys').executable,str(Path(h.__file__).resolve()),'serve-local','--policy',a.policy,'--policy-sha256','0'*64,'--approval',str(directory/'NO-APPROVAL'),'--run-dir',str(directory/'must-not-exist')]
        got=subprocess.run(bad,env=h.clean_env(directory),capture_output=True,text=True)
        (directory/'wrongpin.stdout').write_text(got.stdout);(directory/'wrongpin.stderr').write_text(got.stderr)
        assert got.returncode!=0 and 'policy pin mismatch' in got.stderr and not (directory/'must-not-exist').exists()
        record('wrong pin rejects actual launcher before setup/process/upstream',argv=bad,exit=got.returncode)
        peer_class=h.fake_peer(p,directory)
        with h.server(peer_class) as peer:
            upstream=('127.0.0.1',peer.server_port)
            base={'model':'claude-sonnet-4-6','max_tokens':4096,'stream':True,'messages':[{'role':'user','content':'local guard receipt'}]}
            for name,body,path in [('wrong model',{**base,'model':'NOT-A-MODEL'},'/v1/messages'),('wrong native max',{**base,'max_tokens':64},'/v1/messages'),('wrong path',base,'/elsewhere')]:
                f=h.Forward(p,a.policy,a.policy_sha256,upstream);all_forwards.append(f);before=len(peer_class.bodies)
                with h.server(f.handler()) as proxy:status,_=request(proxy.server_port,body,path)
                assert status==502 and f.slots==0 and len(peer_class.bodies)==before
                record(name,status=status,upstream_requests=0,events=f.events)
            f=h.Forward(p,a.policy,a.policy_sha256,upstream);all_forwards.append(f)
            with h.server(f.handler()) as proxy:
                gui=h.GUI(p,proxy.server_port,directory,1);began=time.monotonic()
                gui.process.stdin.write('tools slow\n');gui.process.stdin.flush()
                # Controlled fixture only: allow the full 90-second observation ceiling, no short model-work timeout.
                while time.monotonic()-began<90:
                    if f.failed:raise AssertionError('positive forwarding failed')
                    if sum(e['outcome']=='forwarded' for e in f.events)==2 and 'Tools finished.' in gui.path.read_text():break
                    time.sleep(.05)
                else:raise AssertionError('local positive did not complete within observation ceiling')
                assert f.events[0]['elapsed_seconds']>3
                assert (directory/'workspace/survivor.txt').read_text()=='survived'
                assert (directory/'workspace/second.txt').read_text()=='second'
                assert [e['body_sha256'] for e in f.events]==peer_class.bodies
                record('actual frozen GUI / delayed streaming / real owned file tools',gui_pid=gui.process.pid,gui_sha256=p['gui_sha256'],elapsed_seconds=round(time.monotonic()-began,3),forwarded=f.events.copy(),served_assets_verified=True,body_bytes_preserved=True)
                gui.stop();assert gui.process.returncode==0;gui=None
                # Two slots left: three simultaneous valid sends must forward exactly two.
                with concurrent.futures.ThreadPoolExecutor(max_workers=3) as pool:statuses=list(pool.map(lambda _:request(proxy.server_port,base)[0],range(3)))
                assert sorted(statuses)==[200,200,502],statuses
                assert f.slots==4 and len(peer_class.bodies)==4
                assert sorted(e['body_sha256'] for e in f.events if e['outcome']=='forwarded')==sorted(peer_class.bodies)
                record('atomic over-cap rejects fifth request',statuses=statuses,reserved_slots=f.slots,actual_upstream_requests=len(peer_class.bodies),events=f.events)
            peer_class.redirect=True;f=h.Forward(p,a.policy,a.policy_sha256,upstream);all_forwards.append(f);before=len(peer_class.bodies)
            with h.server(f.handler()) as proxy:status,_=request(proxy.server_port,base)
            assert status==502 and len(peer_class.bodies)==before+1 and f.events[0]['status']==302
            record('redirect rejected, never followed',status=status,events=f.events)
        assert all(f.key_reads==0 and f.external_requests==0 for f in all_forwards)
        record('no key/external/browser/native execution',key_reads=0,external_requests=0,browser_processes=0,native_calls=0,qualification='actual exercised proxy branches are fixed 127.0.0.1 destinations; counters plus source inspection, not a general OS network audit')
        receipt['passed']=True
    except Exception as e:
        receipt['passed']=False;receipt['unexpected_failure']={'type':type(e).__name__,'message':str(e)};raise
    finally:
        for f in all_forwards:f.abort()
        if gui:gui.stop()
        receipt['elapsed_seconds']=round(time.monotonic()-start,3);receipt['all_owned_processes_stopped']=True;h.save(directory/'receipt.json',receipt)
if __name__=='__main__':main()
