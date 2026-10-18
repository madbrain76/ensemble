---
name: primary
type: primary
description: Chapter 9 bounded skill lifecycle exercise
tools: read_file
loadable-skills: writer
---
CH09_PRIMARY. Use only the tools needed for the user's finite request. Explain each tool batch before calling tools. To write files, first load writer; then follow its returned instructions. Never execute commands or access other directories. Initial variables: CUSTOM=$CUSTOM_VAR UNKNOWN=$UNRESOLVED_CH09.
Tools available at construction:
$TOOLS
Skills offered:
$SKILLS
