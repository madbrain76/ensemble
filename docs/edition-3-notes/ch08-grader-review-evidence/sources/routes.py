import sys;sys.dont_write_bytecode=True
from common import *
pins=json.loads((E/'pins-before.json').read_text());cli=pins['reused_artifacts'][0]['path'];binary=pins['candidate']['binary']['path']
args=['-json','-core',str(R/'solutions/ch08'),'-gui',str(R/'solutions/ch08/gui'),'-consumer',str(R/'solutions/ch08/ch05'),'-workflow',str(R/'solutions/ch08/ch06'),'-stream-exercise',str(R/'solutions/ch08/ch07'),'-stream-model','claude-sonnet-4-6','-gui-model','claude-sonnet-4-6','-consumer-case',str(T/'consumer-case.json'),'-media-models',str(T/'media-models.json'),cli]
for ch in [8,7]:
 name='historical-ch'+str(ch);assert run(name,[binary,'-ch',str(ch)]+args,T/'work',420)==0
 report=json.loads((E/(name+'.stdout')).read_text());assert report['passed'] and report['score']==100 and report['max_score']==100,report
