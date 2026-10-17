# Three-agent manuscript workflow

This independent module imports only `ensemble`. It owns author/editor/reviewer
configuration, prompts, role tools, synchronized manuscript storage, observations,
and sequential handoffs; the framework has no role or workflow knowledge.

```
GOWORK=off go build -o /tmp/ch06 .
# Use the standard LLM_VENDOR, LLM_MODEL, LLM_API_KEY environment.
printf '%s\n' '{"kind":"prompt","text":"Write an 80-word account of a rainy library."}' | /tmp/ch06
```

While the workflow runs, stdin remains open for `kind:hint` or `kind:interrupt`
with `agent:author`, `agent:editor`, or `agent:reviewer`. A prompt is accepted once
per process. Hints are classified by the targeted Agent's current state. Interruption
of a submitted stage ends that reliable request with an error and prevents later
stage handoffs. EOF waits for admitted stages; it does not abort the last request.

Author owns `write_draft`, `word_count`, `think`; editor owns `read_draft`,
`edit_draft`; reviewer owns `review`. Tools validate ordinary JSON arguments and
have real file effects in `manuscript.txt` and `review.json`. The reviewer receives
the stored edited manuscript, never the author's final chat summary.

`CH06_LOG` selects observation JSONL (default `workflow.jsonl`); per-role event logs
are adjacent `.author.jsonl`, `.editor.jsonl`, `.reviewer.jsonl` files. Observations
are bounded/lossy and callbacks never write stdout. Application-owned `turn_ended`
and pipeline markers are emitted after reliable Ask replies, not inferred from
observation timing. Thinking/opaque parts are excluded from executable output.

`GOWORK=off go test -race ./...` exercises real public Agents over local fake HTTP,
real role tool side effects and final reliable handoffs. It is not live-provider
evidence. Credentials must not be placed in arguments, source or receipts.
