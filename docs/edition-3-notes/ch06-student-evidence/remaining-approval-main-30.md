# Separate remaining-only authorization

Parent message main-30 approved ONCE the exact command:
`python3 /Users/bill/projects/ensemble/docs/edition-3-notes/ch06-student-evidence/live_harness_remaining.py --remaining --approved-live`

Launcher SHA256 ad6d5603e637004567d06b60a392c6ca98e5e2398d7a0fba8993a9cd1503f0c7.
Same clean source, source-manifest, archive/core/shape pins as original final freeze.
Fixed order OpenAI human, Gemini human, OpenAI shape streaming, OpenAI shape plain.
At most 3 forwards/case, 12 total, 1024 requested output tokens/request. Read only
named directOpenAIAPIKey/directGeminiAPIKey in process after preflight; no child
key/argv/file leakage. No Anthropic, discovery, substitution, retry or additional
case. Stop first failure/unavailability; no implicit permission to continue.
Original live-v2 failed check stays unchanged. New verifier explicitly checks
next-response-before-tools behavior, not persistence through a later tool cycle.
This is exercise authorization, not implementation acceptance.
