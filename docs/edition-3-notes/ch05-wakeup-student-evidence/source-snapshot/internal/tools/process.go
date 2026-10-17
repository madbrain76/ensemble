package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"ensemble/internal/common"
	"github.com/creack/pty"
)

func (t *Tool) path(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if strings.HasPrefix(filepath.ToSlash(filepath.Clean(path)), "cr/io/") {
		return filepath.Join(t.parent.Config().DataDir, path)
	}
	return filepath.Join(t.parent.Config().Workspace, path)
}

// Start returns immediately after attaching a process; its reader owns Finish.
func (t *Tool) Start(ctx context.Context, j common.Job, raw json.RawMessage) (string, bool, error) {
	if t.definition.Name != "run_command" {
		out, err := t.execute(ctx, raw, true)
		return out, false, err
	}
	var a struct {
		Command string   `json:"command"`
		CWD     string   `json:"cwd"`
		Delay   *float64 `json:"ai_callback_delay"`
		Pattern *string  `json:"ai_callback_pattern"`
		Cap     *int     `json:"max_output_bytes"`
	}
	if err := decode(raw, &a, "command"); err != nil {
		return "", false, err
	}
	if strings.TrimSpace(a.Command) == "" {
		return "", false, fmt.Errorf("command must not be empty")
	}
	if err := ctx.Err(); err != nil {
		return "", false, err
	}
	cwd := t.parent.Config().Workspace
	if a.CWD != "" {
		cwd = a.CWD
		if !filepath.IsAbs(cwd) {
			cwd = filepath.Join(t.parent.Config().Workspace, cwd)
		}
		j.SetCWD(cwd)
	}
	cmd := exec.Command("/bin/sh", "-c", a.Command)
	cmd.Dir = cwd
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		switch key {
		case "LLM_API_KEY", "ANTHROPIC_API_KEY", "OPENAI_API_KEY", "GEMINI_API_KEY":
			continue
		}
		if !strings.HasPrefix(entry, "TERM=") {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	cmd.Env = append(cmd.Env, "TERM=dumb")
	err := j.LaunchProcess(func() (io.WriteCloser, func() error, error) {
		terminal, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 50, Cols: 200})
		if err != nil {
			return nil, nil, fmt.Errorf("start PTY: %w", err)
		}
		done := make(chan struct{})
		stop := func() error {
			err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			if errors.Is(err, syscall.ESRCH) {
				err = nil
			}
			_ = terminal.Close()
			<-done
			return err
		}
		go func() {
			defer close(done)
			// Keep a trailing CR across reads so split CRLF pairs normalize correctly.
			buf := make([]byte, 8192)
			pending := ""
			var readFailure error
			for {
				n, err := terminal.Read(buf)
				if n > 0 {
					chunk := pending + string(buf[:n])
					pending = ""
					if strings.HasSuffix(chunk, "\r") {
						pending = "\r"
						chunk = strings.TrimSuffix(chunk, "\r")
					}
					if _, writeErr := io.WriteString(j, strings.ReplaceAll(chunk, "\r\n", "\n")); writeErr != nil {
						break
					}
				}
				if err != nil {
					if !errors.Is(err, io.EOF) && !errors.Is(err, syscall.EIO) && !errors.Is(err, os.ErrClosed) {
						readFailure = err
					}
					break
				}
			}
			if pending != "" {
				_, _ = io.WriteString(j, pending)
			}
			waitErr := cmd.Wait()
			_ = terminal.Close()
			code := 0
			if waitErr != nil {
				var exit *exec.ExitError
				if errors.As(waitErr, &exit) {
					code = exit.ExitCode()
				} else {
					readFailure = waitErr
					code = -1
				}
			}
			j.Finish(fmt.Sprintf("\nexit_code: %d\n", code), readFailure, &code)
		}()
		return terminal, stop, nil
	})
	if err != nil {
		return "", false, err
	}
	return "", true, nil
}
