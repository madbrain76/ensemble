// Package diagnostics owns retention of complete API and GUI trace records.
package diagnostics

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Log is the file destination for existing complete line records.
// Its limit is read at each retention boundary, so settings apply without reopening it.
// This diagnostic policy is deliberately unrelated to recovery archives.
type Log struct {
	mu    sync.Mutex
	path  string
	limit func() int
}

func Open(path string, limit func() int) (*Log, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	return &Log{path: path, limit: limit}, nil
}
func (d *Log) Close() error { return nil }
func (d *Log) Write(data []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(data) == 0 {
		return 0, nil
	}
	if data[len(data)-1] != '\n' {
		return 0, errors.New("diagnostic record must end with newline")
	}
	f, err := os.OpenFile(d.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return 0, err
	}
	n, err := f.Write(data)
	return n, errors.Join(err, f.Close())
}

// Retain applies policy only at a caller-established safe recovery boundary.
// It shares Write's lock and replaces the whole file; later appends reopen the
// path, so they cannot retain an old offset into a replaced file.
func (d *Log) Retain() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	limit := 200
	if d.limit != nil {
		limit = d.limit()
	}
	if limit < 0 {
		return errors.New("negative diagnostic record limit")
	}
	if limit == 0 {
		return nil
	}
	old, err := os.ReadFile(d.path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	// Drop an incomplete old trailing record, not half of a new logical record.
	if len(old) > 0 && old[len(old)-1] != '\n' {
		old = old[:bytes.LastIndexByte(old, '\n')+1]
	}
	records := old
	start, count := len(records), 0
	for start > 0 {
		start--
		if records[start] == '\n' {
			count++
			if count > limit {
				start++
				break
			}
		}
	}
	f, err := os.CreateTemp(filepath.Dir(d.path), ".diagnostic-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	n, err := f.Write(records[start:])
	if err == nil && n != len(records[start:]) {
		err = io.ErrShortWrite
	}
	err = errors.Join(err, f.Close())
	if err == nil {
		err = os.Rename(f.Name(), d.path)
	}
	if err != nil {
		return err
	}
	return nil
}
