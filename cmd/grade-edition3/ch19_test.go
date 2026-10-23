package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestChapter19GuardsAndRoute(t *testing.T) {
	dir := t.TempDir()
	core := filepath.Join(dir, "core")
	if err := os.Mkdir(core, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(core, "go.mod"), []byte("module course.invalid/sandbox\n\ngo 1.25\n"), 0600); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(dir, "launched")
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("#!/bin/sh\nprintf launched >> "+strconv.Quote(marker)+"\nexit 71\n"), 0700); err != nil {
		t.Fatal(err)
	}
	bad := [][]string{{"-ch", "19", target}, {"-ch", "19", "-core", "", target}, {"-ch", "19", "-core", filepath.Join(dir, "absent"), target}, {"-ch", "19", "-core", core}, {"-ch", "20", target}}
	for _, flag := range []string{"gui", "consumer", "consumer-case", "workflow", "media-models", "stream-exercise", "stream-model", "gui-model", "cache-prices", "cache-current", "cache-prior", "status-prices", "status-second-prices"} {
		bad = append(bad, []string{"-ch", "19", "-core", core, "-" + flag, "", target})
	}
	for _, args := range bad {
		if code := run(args); code != 2 {
			t.Fatalf("%v exit%d", args, code)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("invalid route launched target: %v", args)
		}
	}
	out, err := os.CreateTemp(dir, "report-")
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = out
	code := run([]string{"-ch", "19", "-core", core, "-json", target})
	os.Stdout = old
	if err = out.Close(); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(out.Name())
	if err != nil {
		t.Fatal(err)
	}
	var report struct {
		Title    string
		MaxScore int `json:"max_score"`
		Checks   []struct{ ID string }
	}
	if err = json.Unmarshal(body, &report); err != nil {
		t.Fatalf("%v %s", err, body)
	}
	if code != 1 || !strings.Contains(report.Title, "Chapter 19") || report.MaxScore != 100 || len(report.Checks) != 7 {
		t.Fatalf("wrong route: exit%d %s", code, body)
	}
	if runtime.GOOS == "darwin" {
		if _, err = os.Stat(marker); err != nil {
			t.Fatal(err)
		}
	}
}
