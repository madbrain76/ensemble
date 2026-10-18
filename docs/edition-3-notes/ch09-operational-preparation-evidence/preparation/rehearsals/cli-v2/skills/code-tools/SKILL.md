---
name: code-tools
type: loadable
description: Write and verify the Chapter 9 note
tools: write_file
depends: helper
---
CH09_CODE_TOOLS instruction: write note.txt with exact content SKILLS_OK, then read note.txt and verify its literal content. Request unload_skill for code-tools; unload is pending, not removal. If the user explicitly requests a subsequent write, it remains permitted. Answer with the verified content and explain pending unload. Never perform an unrelated action. CUSTOM=$CUSTOM_VAR UNKNOWN=$UNRESOLVED_CH09.
$TOOLS
$SKILLS
