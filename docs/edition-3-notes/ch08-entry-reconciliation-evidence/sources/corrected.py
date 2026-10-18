from run import *
if run('build-driver-corrected',['go','build','-o',str(T/'bin/driver'),'.'],T/'driver',120):sys.exit(2)
argv=json.loads((E/'historical-full-first-invalid.receipt.json').read_text())['argv']
sys.exit(run('historical-full-corrected',argv,T/'work',420))
