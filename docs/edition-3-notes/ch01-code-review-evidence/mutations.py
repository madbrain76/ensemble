"""Reviewer sensitivity probes; changes only disposable exact-candidate copies."""
from pathlib import Path
import shutil, subprocess, tempfile
E = Path(__file__).resolve().parent
source = Path((E / 'scratch-path.txt').read_text().strip())
mutations = [
 ('first-block-only', 'internal/engine/engine.go', [('text.WriteString(block.Text)', 'text.WriteString(block.Text)\n break')], '^TestConversation$'),
 ('overwrite-usage', 'internal/engine/engine.go', [('e.usage.Input +=', 'e.usage.Input ='), ('e.usage.Output +=', 'e.usage.Output =')], '^TestConversation$'),
 ('drop-history', 'internal/agent/agent.go', [('append(a.History(), common.Message', 'append([]common.Message(nil), common.Message')], '^TestConversation$'),
 ('drop-root-ownership', 'ensemble.go', [('e.agents = append(e.agents, a)', '// omitted by mutation')], '^TestConversation$'),
 ('lose-empty-answer-usage', 'internal/engine/engine.go', [('\te.usage.Input += *reply.Usage.Input\n\te.usage.Output += *reply.Usage.Output\n', ''), ('\te.parent.Parent().Logger().Printf(', '\te.usage.Input += *reply.Usage.Input\n\te.usage.Output += *reply.Usage.Output\n\te.parent.Parent().Logger().Printf(')], '^TestFailuresDoNotCommit$/empty$'),
 ('missing-usage-is-zero', 'internal/engine/engine.go', [('\tif reply.Usage.Input == nil || reply.Usage.Output == nil {', '\tzero := 0\n if reply.Usage.Input == nil { reply.Usage.Input = &zero }\n if reply.Usage.Output == nil { reply.Usage.Output = &zero }\n if reply.Usage.Input == nil || reply.Usage.Output == nil {')], '^TestFailuresDoNotCommit$/missing'),
 ('pollute-stdout', 'cmd/ch01/main.go', [('\tchat :=', '\tfmt.Fprintln(os.Stdout, "debug")\n\tchat :=')], '^TestCLI$/protocol$'),
 ('ignore-cancellation', 'internal/engine/engine.go', [('http.NewRequestWithContext(ctx,', 'http.NewRequestWithContext(context.Background(),')], '^TestCancellation$/headers$'),
 ('one-minute-timeout', 'internal/engine/engine.go', [('"strings"', '"strings"\n "time"'), ('client: &http.Client{}', 'client: &http.Client{Timeout: 60*time.Second}')], '^TestRequestLifetime$/headers$'),
]
with (E / 'mutations.txt').open('w') as log:
 for name, file, replacements, test in mutations:
  with tempfile.TemporaryDirectory(prefix='e3-ch01-mutant.') as directory:
   d = Path(directory)
   shutil.copytree(source, d, dirs_exist_ok=True, ignore=lambda path, names: ['ch01'] if Path(path) == source else [])
   p = d / file
   text = p.read_text()
   for old, new in replacements:
    assert text.count(old) == 1, (name, old)
    text = text.replace(old, new)
   p.write_text(text)
   command = ['go', 'test', '-count=1', '-v', '-run', test, '.']
   r = subprocess.run(command, cwd=d, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
   log.write(f'MUTATION {name}\nFILE {file}\nREPLACEMENTS {replacements!r}\nCOMMAND {command!r}\n{r.stdout}EXIT {r.returncode}\n\n')
   log.flush()
   assert r.returncode != 0 and '--- FAIL:' in r.stdout and '[build failed]' not in r.stdout, name
 log.write('All nine mutants rejected behaviorally; baseline archive and student repository unchanged.\n')
print('All nine reviewer mutations rejected behaviorally.')
