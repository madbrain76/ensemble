from pathlib import Path
p=Path('book/edition-3/chapter-10.md')
s=p.read_text()
a=s.index('**Exercise.** Extend your Chapter 9 agent')
b=s.index('| Check | Points | Proves |',a)
s=s[:a]+'''**Exercise.** Extend your Chapter 9 agent with this save/load lifecycle in
both CLI and GUI. Keep one public Agent implementation behind them.

Build the executable package, not necessarily the module root. A library root
can be imported by the grader without being a runnable command. The fresh
student implementation keeps its CLI in `cmd/ch01`, for example:

```sh
cd /absolute/path/to/your/core-module
go build -o /tmp/ensemble-ch10 ./cmd/ch01
```

Use your own CLI package path. Use `go build -o /tmp/ensemble-ch10 .` only if
that directory actually contains your command's `main` package.

From the **book repository root**, run the current Edition 3 command with all
five module roots and both declared model IDs. Replace the placeholders with
your implementation's paths and capability rows:

```sh
go run ./cmd/grade-edition3 -ch 10 -json \\
  -core /absolute/path/to/core-module \\
  -consumer /absolute/path/to/consumer-module \\
  -consumer-case /absolute/path/to/consumer-case.json \\
  -workflow /absolute/path/to/workflow-module \\
  -media-models /absolute/path/to/media-map.json \\
  -stream-exercise /absolute/path/to/stream-exercise-module \\
  -stream-model YOUR_DECLARED_STREAMING_MODEL \\
  -gui /absolute/path/to/gui-module \\
  -gui-model YOUR_DECLARED_GUI_MODEL \\
  /tmp/ensemble-ch10

go run ./cmd/grade-edition3 -ch 1 -json /tmp/ensemble-ch10
```

The consumer-case and media-model map paths are optional inputs inherited from
the earlier chapters; remove those two flags if you do not supply them. The
media map selects your declared `anthropic`, `openai`, and `gemini` model rows.
It does not establish live-provider capability. All flags precede the explicit
CLI executable; the core module and executable are separate inputs.

Both reports must have `passed: true` and exit zero on the **same frozen
binary**. Chapter 1 still runs the real 62-second delayed-header and delayed-body
phases. Allow both to finish: a 90-second outer timeout cannot cover them, and
a 100/100 score can still fail a required zero-point check.

The GUI stays a separate module; the grader uses the appropriate command for
each phase rather than applying headless checks to a GUI server.

''' + s[b:]
s=s.replace("The student's first save code copied `Config.SystemPrompt`.","The earlier Edition 2 implementation's first save code copied `Config.SystemPrompt`.")
s=s.replace('The correction asks Engine for the effective instruction through its existing','That correction asks Engine for the effective instruction through its existing')
needle='## §10.6 Default Load, and Refusing a Bad File\n'
assert needle in s
s=s.replace(needle,'''The browser has a related continuity problem. A provider's tool-call ID and a
transcript card's address serve different owners. Reusing a parser's turn-local
number as a global card address can make an old result select a new card after
restart. The fresh student implementation had to preserve old card addresses
while allocating new ones beyond the restored conversation.

That includes loading a legacy archive whose parser-local IDs already collide.
Repairing the display must leave provider payloads and call/result correlations
intact. Renumbering those facts to make the DOM convenient would change the
conversation being replayed. I want the old card to keep its address and a new
card to get an unused one, without prescribing a new public identity API.

## §10.6 Default Load, and Refusing a Bad File
''')
a=s.index('Run the agent in an empty directory, ask it to remember a word, and')
b=s.index('Then break copies on purpose.',a)
s=s[:a]+'''Use a fresh scratch directory so you cannot accidentally load an earlier
conversation. The following launch matches the fresh implementation's
`cmd/ch01` command built above. It chooses the model through `LLM_MODEL` and
puts `--save` **before** the `chat` subcommand. Other implementations may expose
a different launcher; use their documented model selection and save-path options.

Configure credentials through your normal authorized provider setup, never by
putting a key in a command or prompt. This example uses the declared
`claude-sonnet-4-6` row; choose a supported model for your implementation before
beginning this bounded exercise.

```sh
work=$(mktemp -d "${TMPDIR:-/tmp}/ensemble-ch10.XXXXXX")
cd "$work"
printf 'hello from disk\\n' > note.txt
LLM_VENDOR=anthropic LLM_MODEL=claude-sonnet-4-6 \\
  /tmp/ensemble-ch10 --save "$work/save.json" chat
```

Ask the agent to read `note.txt` and summarize it. Verify an actual `read_file`
call and result, not just an answer that could have guessed the file contents.
Wait for the completed answer, then send EOF with Ctrl-D on an empty line.
Wait for the process to exit successfully and check that `save.json` exists.

Restart from that same directory with the same model and save path:

```sh
LLM_VENDOR=anthropic LLM_MODEL=claude-sonnet-4-6 \\
  /tmp/ensemble-ch10 --save "$work/save.json" chat
```

Before sending a prompt, verify that loading the archive alone has not made a
provider request. Ask, "What did I ask you to read?" Check the restored answer,
then send EOF again and confirm another successful exit. If driving the
accepted CLI over JSONL instead, omit `chat` and send a live-protocol record
such as `{"kind":"prompt","text":"What did I ask you to read?"}`; close stdin
and wait for completion. Do not mix that envelope with the human chat input.

The explicit path above selects the ordinary filename. In a second fresh
directory, use `--save chosen.json` on both starts to check the override: it
should resume `chosen.json` without also creating `save.json`. Keep the complete
save from the first run for the independent variants below.

Repeat the tool-bearing turn through an owned GUI instance and browser profile,
using a separate scratch directory and an unused loopback port. Keep read-aloud
and microphone inactive. Send the first prompt through the GUI's JSONL stdin,
let its tool turn finish, and record the browser's tool/result cards and IDs.
Close stdin, wait for the server to exit, and confirm that `save.json` appears.

Restart from that same GUI working directory and reopen the browser view. The
restored conversation and old card addresses should be present before a new
prompt, with no provider request caused merely by loading the archive or
requesting state. Make a new tool-bearing turn, then check that its new card IDs
do not collide with the old ones and that old controls still select old results.
Close the owned server and browser after the check. A backend replay test alone
cannot establish these browser properties.

''' + s[b:]
s += '''
The fresh Edition 3 run encountered a different problem. Its first frozen
Chapter 10 grade was **80/100**: both default-load checks failed, while snapshot,
replay-tail and refusal checks passed. The student preserved the failed report
and stopped rather than changing persistence code to fit an uncertain diagnosis.

The subsequent investigation corrected the test path. The default-load probe
needed a valid minimal running configuration, and inherited skill checks needed
to honor the implementation's declared models. The source stayed unchanged.
With those corrections, Chapter 10 passed **100/100**, and the same frozen CLI
passed Chapter 1, including its delayed-header and delayed-body phases. Restoring
conversation and replacing runtime configuration are separate operations; a
test that confuses them can reward the wrong restart behavior.

The operational exercise then took a real tool turn through EOF and relaunch:
read a small file, save, and answer from the restored conversation. A first paid
attempt failed at the transport boundary; a bounded TLS recovery completed both
launches. The original failure's cause remained unproven. Separately, an owned
browser exercise used a local fake provider to check old and new card IDs across
a GUI restart. Its DOM assertions passed; captured screenshots were not a human
visual or audio review. These observations support specific restart claims,
not a new multi-provider campaign or a claim that every browser behavior was
inspected.
'''
p.write_text(s)
print('chapter updated:',len(s.encode()),'bytes;',len(s.splitlines()),'lines')
