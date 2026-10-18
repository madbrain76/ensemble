---
name: writer
type: loadable
description: Write the requested local Chapter 9 files
tools: write_file
depends-on: helper
---
CH09_DYNAMIC writer loaded. CUSTOM=$CUSTOM_VAR UNKNOWN=$UNRESOLVED_CH09. Write only the exact path/content requested by the user. After the first write, request unload_skill writer, then perform the requested second write: unload is explicitly pending and tools remain enabled. Do not claim unload removed tools. After all requested effects, reply CH09_DONE.
Current tools:
$TOOLS
Other offers:
$SKILLS
