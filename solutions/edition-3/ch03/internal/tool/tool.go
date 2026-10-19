// Package tool owns one Agent's six synchronous built-ins. Paths resolve from
// the launch directory; this is not a sandbox or a filesystem confinement claim.
package tool

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"ensemble/internal/common"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type registry struct{ privateState }
type privateState struct {
	parent    common.Agent
	directory string
	entries   []entry
}

// Each registered entry joins the declaration and the function it advertises.
// The registry owns these closures; common exposes no handler implementation.
type entry struct {
	definition common.ToolDefinition
	execute    func(context.Context, arguments) (string, error)
}

func New(parent common.Agent) (common.ToolRegistry, error) {
	if parent == nil {
		return nil, fmt.Errorf("Tools requires an Agent parent")
	}
	directory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("reading launch directory: %w", err)
	}
	r := &registry{privateState: privateState{parent: parent, directory: directory}}
	available := r.builtins()
	selected := parent.Config().Tools
	if selected == nil {
		r.entries = available
		return r, nil
	}
	for _, name := range selected {
		for _, registered := range r.entries {
			if registered.definition.Name == name {
				return nil, fmt.Errorf("duplicate enabled tool %q", name)
			}
		}
		found := false
		for _, candidate := range available {
			if candidate.definition.Name == name {
				r.entries = append(r.entries, candidate)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("unknown enabled tool %q", name)
		}
	}
	return r, nil
}
func (r *registry) Parent() common.Agent { return r.parent }
func (r *registry) Definitions() []common.ToolDefinition {
	out := make([]common.ToolDefinition, len(r.entries))
	for i, registered := range r.entries {
		out[i] = registered.definition
		out[i].Parameters = append(json.RawMessage(nil), registered.definition.Parameters...)
	}
	return out
}
func (r *registry) path(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(r.directory, path)
}

func (r *registry) Dispatch(ctx context.Context, name string, input json.RawMessage) (common.PartList, bool) {
	fail := func(err error) (common.PartList, bool) { return common.PartList{common.Text(err.Error())}, true }
	var selected *entry
	for i := range r.entries {
		if r.entries[i].definition.Name == name {
			selected = &r.entries[i]
			break
		}
	}
	if selected == nil {
		return fail(fmt.Errorf("unknown or disabled tool %q", name))
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	args, err := decode(input)
	if err != nil {
		return fail(err)
	}
	output, err := selected.execute(ctx, args)
	if err != nil {
		return fail(err)
	}
	return common.PartList{common.Text(output)}, false
}

type arguments map[string]json.RawMessage

func decode(raw json.RawMessage) (arguments, error) {
	var a arguments
	if err := json.Unmarshal(raw, &a); err != nil || a == nil {
		return nil, fmt.Errorf("arguments must be a JSON object")
	}
	return a, nil
}
func (a arguments) only(keys ...string) error {
	for key := range a {
		found := false
		for _, allowed := range keys {
			if key == allowed {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown argument %q", key)
		}
	}
	return nil
}
func (a arguments) text(key string, required bool) (string, error) {
	raw, ok := a[key]
	if !ok {
		if required {
			return "", fmt.Errorf("%s is required", key)
		}
		return "", nil
	}
	var value string
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || json.Unmarshal(raw, &value) != nil {
		return "", fmt.Errorf("%s must be a string", key)
	}
	return value, nil
}
func (a arguments) boolean(key string) (bool, error) {
	raw, ok := a[key]
	if !ok {
		return false, nil
	}
	var value bool
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || json.Unmarshal(raw, &value) != nil {
		return false, fmt.Errorf("%s must be a boolean", key)
	}
	return value, nil
}
func (a arguments) integer(key string, fallback int) (int, error) {
	raw, ok := a[key]
	if !ok {
		return fallback, nil
	}
	var value int
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || json.Unmarshal(raw, &value) != nil {
		return 0, fmt.Errorf("%s must be an integer", key)
	}
	return value, nil
}
func (a arguments) path() (string, error) {
	path, err := a.text("path", true)
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", fmt.Errorf("path must not be empty")
	}
	return path, nil
}

// lines preserves final newlines but does not invent an empty line after EOF.
func lines(b []byte) []string {
	if len(b) == 0 {
		return nil
	}
	out := strings.SplitAfter(string(b), "\n")
	if out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return out
}
func (r *registry) read(a arguments) (string, error) {
	if err := a.only("path", "start_line", "end_line", "max_bytes"); err != nil {
		return "", err
	}
	path, err := a.path()
	if err != nil {
		return "", err
	}
	start, err := a.integer("start_line", 1)
	if err != nil {
		return "", err
	}
	end, err := a.integer("end_line", 0)
	if err != nil {
		return "", err
	}
	maxBytes, err := a.integer("max_bytes", 32768)
	if err != nil {
		return "", err
	}
	if start < 0 || end < 0 || maxBytes < 0 {
		return "", fmt.Errorf("line bounds and max_bytes must be nonnegative")
	}
	if start == 0 {
		start = 1
	}
	if maxBytes == 0 {
		maxBytes = 32768
	}
	if end != 0 && end < start {
		return "", fmt.Errorf("end_line must be zero (EOF) or >= start_line")
	}
	file, err := os.Open(r.path(path))
	if err != nil {
		return "", err
	}
	defer file.Close()
	// Read bounded fragments, not Scanner-sized lines or the whole file. First
	// select the inclusive line range, then cap only the bytes in that range.
	reader := bufio.NewReader(file)
	var out strings.Builder
	for line := 1; end == 0 || line <= end; {
		fragment, readErr := reader.ReadSlice('\n')
		if line >= start {
			remaining := maxBytes - out.Len()
			if len(fragment) > remaining {
				out.Write(fragment[:remaining])
				fmt.Fprintf(&out, "\n[truncated at max_bytes=%d]\n", maxBytes)
				return out.String(), nil
			}
			out.Write(fragment)
		}
		if len(fragment) > 0 && fragment[len(fragment)-1] == '\n' {
			line++
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil && readErr != bufio.ErrBufferFull {
			return "", readErr
		}
	}
	return out.String(), nil
}
func (r *registry) list(a arguments) (string, error) {
	if err := a.only("path"); err != nil {
		return "", err
	}
	path, err := a.path()
	if err != nil {
		return "", err
	}
	entries, err := os.ReadDir(r.path(path))
	if err != nil {
		return "", err
	}
	var out strings.Builder
	for _, entry := range entries {
		out.WriteString(entry.Name())
		if entry.IsDir() {
			out.WriteByte('/')
		}
		out.WriteByte('\n')
	}
	return out.String(), nil
}
func (r *registry) search(a arguments) (string, error) {
	if err := a.only("path", "pattern", "context_lines"); err != nil {
		return "", err
	}
	path, err := a.path()
	if err != nil {
		return "", err
	}
	pattern, err := a.text("pattern", true)
	if err != nil {
		return "", err
	}
	contextLines, err := a.integer("context_lines", 0)
	if err != nil {
		return "", err
	}
	if contextLines < 0 {
		return "", fmt.Errorf("context_lines must be nonnegative")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regular expression: %w", err)
	}
	var out strings.Builder
	root := r.path(path)
	err = filepath.WalkDir(root, func(file string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		// Do not follow symlink directories or read devices/FIFOs during a search.
		if !d.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if bytes.IndexByte(data, 0) >= 0 {
			return nil
		}
		text := lines(data)
		hits := make([]bool, len(text))
		type span struct{ start, end int }
		var spans []span
		for i, line := range text {
			if !re.MatchString(strings.TrimSuffix(line, "\n")) {
				continue
			}
			hits[i] = true
			start, end := 0, len(text)-1
			if contextLines < i {
				start = i - contextLines
			}
			if contextLines < len(text)-1-i {
				end = i + contextLines
			}
			if len(spans) > 0 && start <= spans[len(spans)-1].end+1 {
				spans[len(spans)-1].end = end
			} else {
				spans = append(spans, span{start, end})
			}
		}
		display := file
		if !filepath.IsAbs(path) {
			rel, err := filepath.Rel(r.directory, file)
			if err != nil {
				return err
			}
			display = rel
		}
		for _, span := range spans {
			if contextLines > 0 && out.Len() > 0 {
				out.WriteString("--\n")
			}
			for i := span.start; i <= span.end; i++ {
				separator := "-"
				if hits[i] {
					separator = ":"
				}
				fmt.Fprintf(&out, "%s%s%d%s%s\n", display, separator, i+1, separator, strings.TrimSuffix(text[i], "\n"))
			}
		}
		return nil
	})
	return out.String(), err
}
func (r *registry) write(a arguments) (string, error) {
	if err := a.only("path", "content", "overwrite", "append"); err != nil {
		return "", err
	}
	path, err := a.path()
	if err != nil {
		return "", err
	}
	content, err := a.text("content", true)
	if err != nil {
		return "", err
	}
	overwrite, err := a.boolean("overwrite")
	if err != nil {
		return "", err
	}
	appendMode, err := a.boolean("append")
	if err != nil {
		return "", err
	}
	if appendMode && overwrite {
		return "", fmt.Errorf("append and overwrite cannot both be true")
	}
	target := r.path(path)
	// Default creation is exclusive, so a racing creator cannot be truncated.
	flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	if overwrite {
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	if appendMode {
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return "", err
	}
	file, err := os.OpenFile(target, flags, 0644)
	if err != nil {
		if !overwrite && !appendMode && os.IsExist(err) {
			info, statErr := os.Stat(target)
			if statErr != nil {
				return "", fmt.Errorf("%s exists; size unavailable (%v); file unchanged; use overwrite:true to replace or append:true to append", path, statErr)
			}
			return "", fmt.Errorf("%s exists (%d bytes); file unchanged; use overwrite:true to replace or append:true to append", path, info.Size())
		}
		return "", err
	}
	_, writeErr := file.WriteString(content)
	closeErr := file.Close()
	if writeErr != nil {
		return "", writeErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	return fmt.Sprintf("wrote %d bytes to %s", len(content), path), nil
}
func (r *registry) edit(a arguments) (string, error) {
	if err := a.only("path", "old_text", "new_text"); err != nil {
		return "", err
	}
	path, err := a.path()
	if err != nil {
		return "", err
	}
	old, err := a.text("old_text", true)
	if err != nil {
		return "", err
	}
	replacement, err := a.text("new_text", true)
	if err != nil {
		return "", err
	}
	if old == "" {
		return "", fmt.Errorf("old_text must not be empty")
	}
	target := r.path(path)
	data, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}
	// Count every match, including overlaps, but retain only bounded locations.
	// Advancing one byte (not the anchor length) preserves overlapping edit sites.
	const locationLimit = 8
	needle := []byte(old)
	index, count := -1, 0
	locations := make([]int, 0, locationLimit)
	line, lineOffset := 1, 0
	for offset := 0; offset <= len(data)-len(needle); {
		relative := bytes.Index(data[offset:], needle)
		if relative < 0 {
			break
		}
		match := offset + relative
		if count == 0 {
			index = match
		}
		count++
		if len(locations) < locationLimit {
			line += bytes.Count(data[lineOffset:match], []byte{'\n'})
			locations = append(locations, line)
			lineOffset = match
		}
		offset = match + 1
	}
	if count == 0 {
		return "", fmt.Errorf("file %q: old_text has 0 matches; file unchanged; reread the file and provide an exact anchor", path)
	}
	if count > 1 {
		omitted := ""
		if count > len(locations) {
			omitted = fmt.Sprintf("; location list truncated: %d more match locations omitted", count-len(locations))
		}
		return "", fmt.Errorf("file %q: old_text has %d matches (including overlaps); match-start lines: %v%s; file unchanged; reread the file and provide a larger exact anchor", path, count, locations, omitted)
	}
	next := make([]byte, 0, len(data)-len(old)+len(replacement))
	next = append(next, data[:index]...)
	next = append(next, replacement...)
	next = append(next, data[index+len(old):]...)
	if err := os.WriteFile(target, next, 0644); err != nil {
		return "", err
	}
	return "replaced one exact occurrence in " + path, nil
}
func (r *registry) run(ctx context.Context, a arguments) (string, error) {
	if err := a.only("command", "cwd"); err != nil {
		return "", err
	}
	command, err := a.text("command", true)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(command) == "" {
		return "", fmt.Errorf("command must not be empty")
	}
	cwd, err := a.text("cwd", false)
	if err != nil {
		return "", err
	}
	if _, present := a["cwd"]; present && cwd == "" {
		return "", fmt.Errorf("cwd must not be empty")
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
	cmd.Dir = r.path(cwd)
	// Shells inherit useful build settings, but not the four provider credentials.
	for _, variable := range os.Environ() {
		name, _, _ := strings.Cut(variable, "=")
		switch name {
		case "LLM_API_KEY", "ANTHROPIC_API_KEY", "OPENAI_API_KEY", "GEMINI_API_KEY":
			continue
		}
		cmd.Env = append(cmd.Env, variable)
	}
	// An explicitly empty environment must not fall back to process inheritance.
	if cmd.Env == nil {
		cmd.Env = []string{}
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err = cmd.Run()
	code := 0
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("command canceled: %w", ctx.Err())
		}
		if exit, ok := err.(*exec.ExitError); ok {
			code = exit.ExitCode()
		} else {
			return "", fmt.Errorf("starting command: %w", err)
		}
	}
	// A nonzero process exit is a completed command result, not a tool transport failure.
	output, err := json.Marshal(struct {
		Stdout   string `json:"stdout"`
		Stderr   string `json:"stderr"`
		ExitCode int    `json:"exit_code"`
	}{stdout.String(), stderr.String(), code})
	return string(output), err
}

func (r *registry) builtins() []entry {
	// Keep definitions local to each registry. They are not process-global state.
	entries := []entry{
		{definition: common.ToolDefinition{Name: "read_file", Description: "Read text with one-based inclusive start_line/end_line, selecting the range before max_bytes caps returned text. Omitted/zero start_line means 1, end_line means EOF, and max_bytes means 32768. Truncated output ends with an explicit [truncated at max_bytes=N] marker; the marker is outside the cap.", Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"},"start_line":{"type":"integer","minimum":0},"end_line":{"type":"integer","minimum":0},"max_bytes":{"type":"integer","minimum":0}},"required":["path"],"additionalProperties":false}`)}, execute: func(_ context.Context, a arguments) (string, error) { return r.read(a) }},
		{definition: common.ToolDefinition{Name: "list_directory", Description: "List immediate entries sorted by name; directories end in /.", Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"}},"required":["path"],"additionalProperties":false}`)}, execute: func(_ context.Context, a arguments) (string, error) { return r.list(a) }},
		{definition: common.ToolDefinition{Name: "search_files", Description: "Recursively search regular text files with a Go regular expression. context_lines defaults to zero. Matching lines use path:line:text; context uses path-line-text; overlapping or adjacent context ranges are merged; positive-context groups are separated by --, including across files.", Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"},"pattern":{"type":"string"},"context_lines":{"type":"integer","minimum":0}},"required":["path","pattern"],"additionalProperties":false}`)}, execute: func(_ context.Context, a arguments) (string, error) { return r.search(a) }},
		{definition: common.ToolDefinition{Name: "write_file", Description: "Write text and create parent directories. Existing files require overwrite:true or append:true; a default refusal reports current size and how to opt in. Append preserves existing bytes. append and overwrite are mutually exclusive.", Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"},"content":{"type":"string"},"overwrite":{"type":"boolean"},"append":{"type":"boolean"}},"required":["path","content"],"additionalProperties":false}`)}, execute: func(_ context.Context, a arguments) (string, error) { return r.write(a) }},
		{definition: common.ToolDefinition{Name: "edit_file", Description: "Replace exactly one nonempty exact old_text anchor with new_text. Zero matches report the count and reread guidance. Multiple matches (including overlaps) report the total and up to eight one-based match-start lines, noting omitted locations. Refusals leave the file unchanged; reread and provide a larger exact anchor.", Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"},"old_text":{"type":"string"},"new_text":{"type":"string"}},"required":["path","old_text","new_text"],"additionalProperties":false}`)}, execute: func(_ context.Context, a arguments) (string, error) { return r.edit(a) }},
		{definition: common.ToolDefinition{Name: "run_command", Description: "Run /bin/sh -c synchronously, returning JSON text with separate stdout, stderr, and exit_code fields. Provider API-key environment variables are removed; ordinary build variables are retained. Optional cwd applies only to this command. Relative paths start at the launch directory; no sandbox is provided. No tool timeout is imposed; caller cancellation applies.", Parameters: json.RawMessage(`{"type":"object","properties":{"command":{"type":"string"},"cwd":{"type":"string"}},"required":["command"],"additionalProperties":false}`)}, execute: r.run},
	}
	// Declaration order is stable and independent of map iteration.
	sort.Slice(entries, func(i, j int) bool { return entries[i].definition.Name < entries[j].definition.Name })
	return entries
}
