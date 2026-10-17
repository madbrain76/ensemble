// Package tools implements the local operations owned by an Agent.
package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"ensemble/internal/common"
)

// Tool has one owner and one stable built-in definition; it holds no session state.
type Tool struct {
	// parent owns registration, policy, dispatch and result recording.
	parent common.Agent
	// definition is immutable model-facing metadata.
	definition common.ToolDefinition
}

// Parent returns the owning Agent.
func (t *Tool) Parent() common.Agent { return t.parent }

// Definition returns metadata with an independent schema buffer.
func (t *Tool) Definition() common.ToolDefinition {
	d := t.definition
	d.Parameters = append(json.RawMessage(nil), d.Parameters...)
	return d
}

// NewBuiltins constructs the CLI's complete registry; library Agents opt in.
func NewBuiltins(parent common.Agent) []common.Tool {
	limitProperties := `"ai_callback_delay":{"type":"number"},"ai_callback_pattern":{"type":"string"},"max_output_bytes":{"type":"integer"}`
	definitions := []struct {
		name, description, properties string
		required                      []string
	}{
		{"read_file", "Read text with one-based inclusive start_line/end_line (omitted or zero: first/EOF) and optional explicit max_bytes. Full uncapped result is retained by the job.", `"path":{"type":"string"},"start_line":{"type":"integer"},"end_line":{"type":"integer"},"max_bytes":{"type":"integer"}`, []string{"path"}},
		{"write_file", "Create a file and parent directories. Existing files require overwrite:true unless append:true; append preserves existing bytes.", `"path":{"type":"string"},"content":{"type":"string"},"overwrite":{"type":"boolean"},"append":{"type":"boolean"}`, []string{"path", "content"}},
		{"edit_file", "Replace exactly one literal old_text occurrence with new_text. Empty or ambiguous anchors fail without writing.", `"path":{"type":"string"},"old_text":{"type":"string"},"new_text":{"type":"string"}`, []string{"path", "old_text", "new_text"}},
		{"list_directory", "List immediate entries sorted by name; directories end in /.", `"path":{"type":"string"}`, []string{"path"}},
		{"search_files", "Recursively search text with a Go regular expression. context_lines defaults to 0; matching path:line:text and context path-line-text use merged grep-style windows.", `"path":{"type":"string"},"pattern":{"type":"string"},"context_lines":{"type":"integer"}`, []string{"path", "pattern"}},
		{"run_command", "Start a fresh /bin/sh in a PTY; merged output and exit_code persist on disk. cwd applies only to this call, relative to Workspace. Delay bounds observation, not execution; a running report is a handle, not cancellation.", `"command":{"type":"string"},"cwd":{"type":"string"},` + limitProperties, []string{"command"}},
		{"wait_for_job", "Wait for unseen output of a retained job, including a finished job. Never cancels work.", `"handle":{"type":"integer"},` + limitProperties, []string{"handle"}},
		{"send_input", "Submit a line to an interactive job (adds newline) and wait for unseen output.", `"handle":{"type":"integer"},"input":{"type":"string"},` + limitProperties, []string{"handle", "input"}},
		{"kill_job", "SIGKILL the process group, or mark a Go job killed (Go cannot stop its goroutine).", `"handle":{"type":"integer"}`, []string{"handle"}},
		{"tool_limits", "Set observation limits for the very next call, whichever tool it is; that call reports consumption.", limitProperties, []string{}},
	}
	out := make([]common.Tool, 0, len(definitions))
	for _, d := range definitions {
		required, _ := json.Marshal(d.required)
		schema := json.RawMessage(`{"type":"object","properties":{` + d.properties + `},"required":` + string(required) + `,"additionalProperties":false}`)
		out = append(out, &Tool{parent: parent, definition: common.ToolDefinition{Name: d.name, Description: d.description, Parameters: schema}})
	}
	return out
}

// decode requires a JSON object, all named required fields, and no unknown fields.
func decode(raw json.RawMessage, dst any, required ...string) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil || fields == nil {
		return fmt.Errorf("arguments must be a JSON object")
	}
	for _, key := range required {
		v, ok := fields[key]
		if !ok || bytes.Equal(bytes.TrimSpace(v), []byte("null")) {
			return fmt.Errorf("missing required argument %s", key)
		}
	}
	for key, value := range fields {
		if bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			return fmt.Errorf("argument %s must not be null", key)
		}
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}
	return nil
}

type readArgs struct {
	// Path identifies a local file.
	Path string `json:"path"`
	// Start and End are optional one-based inclusive line bounds.
	Start *int `json:"start_line"`
	End   *int `json:"end_line"`
	// MaxBytes optionally caps selected text before job-output policy.
	MaxBytes int `json:"max_bytes"`
}
type writeArgs struct {
	// Path identifies the file to create, append or explicitly replace.
	Path string `json:"path"`
	// Content is the complete replacement, including an allowed empty string.
	Content string `json:"content"`
	// Overwrite explicitly authorizes replacing existing bytes.
	Overwrite bool `json:"overwrite"`
	// Append takes precedence over Overwrite and preserves existing bytes.
	Append bool `json:"append"`
}
type editArgs struct {
	// Path identifies the existing file.
	Path string `json:"path"`
	// Old is the exact nonempty unique anchor.
	Old string `json:"old_text"`
	// New replaces the anchor, including allowed deletion.
	New string `json:"new_text"`
}
type pathArgs struct {
	// Path identifies a directory to inspect.
	Path string `json:"path"`
}
type searchArgs struct {
	// Path identifies the tree root.
	Path string `json:"path"`
	// Pattern is a Go regular expression.
	Pattern string `json:"pattern"`
	// ContextLines is the nonnegative number of neighboring lines to include.
	ContextLines int `json:"context_lines"`
}
type commandArgs struct {
	// Command is interpreted by a fresh bash process.
	Command string `json:"command"`
}

// Execute performs the selected synchronous operation under caller cancellation.
func (t *Tool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	switch t.definition.Name {
	case "wait_for_job", "send_input", "kill_job", "tool_limits":
		return t.parent.CallTool(ctx, t.definition.Name, raw)
	}
	j, err := t.parent.Parent().Jobs().Create(t.parent)
	if err != nil {
		return "", err
	}
	j.Start(ctx, t, raw)
	return j.Wait(ctx, common.WaitOptions{Delay: 3 * time.Second, MaxBytes: 16384})
}
func (t *Tool) execute(ctx context.Context, raw json.RawMessage, job bool) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	switch t.definition.Name {
	case "read_file":
		var a readArgs
		if err := decode(raw, &a, "path"); err != nil {
			return "", err
		}
		if a.MaxBytes < 0 || (a.Start != nil && *a.Start < 0) || (a.End != nil && *a.End < 0) {
			return "", fmt.Errorf("line bounds and max_bytes must not be negative")
		}
		if job {
			a.Path = t.path(a.Path)
		}
		data, err := os.ReadFile(a.Path)
		if err != nil {
			return "", err
		}
		lines := strings.SplitAfter(string(data), "\n")
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		start, end := 1, len(lines)
		if a.Start != nil && *a.Start != 0 {
			start = *a.Start
		}
		if a.End != nil && *a.End != 0 {
			end = *a.End
		}
		if len(lines) == 0 && start == 1 && end == 0 {
			return "", nil
		}
		if end < start || start > len(lines) || end > len(lines) {
			return "", fmt.Errorf("line range must be within 1..%d and start <= end", len(lines))
		}
		selected := strings.Join(lines[start-1:end], "")
		cap := a.MaxBytes
		if cap == 0 {
			if job {
				cap = len(selected)
			} else {
				cap = 32768
			}
		}
		if len(selected) > cap {
			return selected[:cap] + fmt.Sprintf("\n[truncated at %d bytes]", cap), nil
		}
		return selected, nil
	case "write_file":
		var a writeArgs
		if err := decode(raw, &a, "path", "content"); err != nil {
			return "", err
		}
		if a.Path == "" {
			return "", fmt.Errorf("path must not be empty")
		}
		if job {
			a.Path = t.path(a.Path)
		}
		if err := os.MkdirAll(filepath.Dir(a.Path), 0755); err != nil {
			return "", err
		}
		flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
		if a.Append {
			flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
		} else if a.Overwrite {
			flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		}
		file, err := os.OpenFile(a.Path, flags, 0644)
		if err != nil {
			if os.IsExist(err) {
				info, statErr := os.Stat(a.Path)
				if statErr != nil {
					return "", statErr
				}
				return "", fmt.Errorf("file already exists (%d bytes); unchanged: read it, then set overwrite:true to replace or append:true to add", info.Size())
			}
			return "", err
		}
		_, writeErr := io.WriteString(file, a.Content)
		closeErr := file.Close()
		if err := errors.Join(writeErr, closeErr); err != nil {
			return "", err
		}
		return "file written", nil
	case "edit_file":
		var a editArgs
		if err := decode(raw, &a, "path", "old_text", "new_text"); err != nil {
			return "", err
		}
		if a.Old == "" {
			return "", fmt.Errorf("old_text must not be empty")
		}
		if job {
			a.Path = t.path(a.Path)
		}
		data, err := os.ReadFile(a.Path)
		if err != nil {
			return "", err
		}
		text := string(data)
		first := strings.Index(text, a.Old)
		if first < 0 {
			return "", fmt.Errorf("old_text was not found in %q; file unchanged: read the file and retry with an exact unique anchor", a.Path)
		}
		// Search from the next byte, not the end of the match: overlaps are ambiguous too.
		if strings.Contains(text[first+1:], a.Old) {
			return "", fmt.Errorf("old_text matches more than once in %q (possibly overlapping); file unchanged: include more surrounding text for a unique anchor", a.Path)
		}
		text = text[:first] + a.New + text[first+len(a.Old):]
		info, err := os.Stat(a.Path)
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(a.Path, []byte(text), info.Mode().Perm()); err != nil {
			return "", err
		}
		return "file edited", nil
	case "list_directory":
		var a pathArgs
		if err := decode(raw, &a, "path"); err != nil {
			return "", err
		}
		if job {
			a.Path = t.path(a.Path)
		}
		entries, err := os.ReadDir(a.Path)
		if err != nil {
			return "", err
		}
		var out strings.Builder
		for _, entry := range entries {
			out.WriteString(entry.Name())
			if entry.IsDir() {
				out.WriteString("/")
			}
			out.WriteByte('\n')
		}
		return out.String(), nil
	case "search_files":
		var a searchArgs
		if err := decode(raw, &a, "path", "pattern"); err != nil {
			return "", err
		}
		if a.ContextLines < 0 {
			return "", fmt.Errorf("context_lines must not be negative")
		}
		pattern, err := regexp.Compile(a.Pattern)
		if err != nil {
			return "", err
		}
		var out strings.Builder
		rootPath := a.Path
		relative := !filepath.IsAbs(a.Path)
		if job {
			rootPath = t.path(a.Path)
		}
		err = filepath.WalkDir(rootPath, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if err := ctx.Err(); err != nil {
				return err
			}
			if !entry.Type().IsRegular() {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if bytes.Contains(data, []byte{0}) {
				return nil
			}
			lines := strings.Split(string(data), "\n")
			if len(lines) > 0 && lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
			matches := make([]bool, len(lines))
			include := make([]bool, len(lines))
			for i, line := range lines {
				if !pattern.MatchString(line) {
					continue
				}
				matches[i] = true
				lo := 0
				if a.ContextLines < i {
					lo = i - a.ContextLines
				}
				hi := len(lines) - 1
				if a.ContextLines < hi-i {
					hi = i + a.ContextLines
				}
				for j := lo; j <= hi; j++ {
					include[j] = true
				}
			}
			if job && relative {
				var err error
				path, err = filepath.Rel(t.parent.Config().Workspace, path)
				if err != nil {
					return err
				}
			}
			previous := -2
			for i, line := range lines {
				if !include[i] {
					continue
				}
				if out.Len() > 0 && previous != i-1 && a.ContextLines > 0 {
					out.WriteString("--\n")
				}
				separator := "-"
				if matches[i] {
					separator = ":"
				}
				fmt.Fprintf(&out, "%s%s%d%s%s\n", path, separator, i+1, separator, line)
				previous = i
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		return out.String(), nil

	}
	return "", fmt.Errorf("unknown tool %q", t.definition.Name)
}
