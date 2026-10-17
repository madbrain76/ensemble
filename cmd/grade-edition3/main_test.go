package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestInvalidInputsFailBeforeGrading(t *testing.T) {
	root := t.TempDir()
	module := filepath.Join(root, "module")
	if err := os.Mkdir(module, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.org/fixture\n\ngo 1.25\n"), 0600); err != nil {
		t.Fatal(err)
	}
	badCase := filepath.Join(root, "bad-case.json")
	if err := os.WriteFile(badCase, []byte("not JSON"), 0600); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(root, "missing")
	marker := filepath.Join(root, "invoked")
	target := filepath.Join(root, "sentinel")
	if err := os.WriteFile(target, []byte("#!/bin/sh\nprintf invoked >> "+strconv.Quote(marker)+"\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	cases := [][]string{
		{}, {""}, {"-ch", "2", target}, {"-ch", "4"},
		{"-ch", "1", "-core", module, target},
		{"-ch", "1", "-consumer", "", target},
		{"-ch", "1", "-consumer-case", "", target},
		{"-ch", "4", target},
		{"-ch", "4", "-core", module, target},
		{"-ch", "4", "-consumer", module, target},
		{"-ch", "4", "-core", missing, "-consumer", module, target},
		{"-ch", "4", "-core", module, "-consumer", root, target},
		{"-ch", "4", "-core", module, "-consumer", module, "-consumer-case", badCase, target},
		{"-ch", "4", "-core", module, "-consumer", module, "-consumer-case", missing, target},
		{"-ch", "4", "-core", module, "-consumer", module, missing},
		{"-ch", "1", missing}, {"-ch", "1", target, target},
		{"-unknown", target},
	}
	for _, args := range cases {
		t.Run("args="+strings.Join(args, " "), func(t *testing.T) {
			got := run(args)
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatalf("invalid input launched executable: marker error %v", err)
			}
			if got != 2 {
				t.Fatalf("run(%q) = %d, want 2", args, got)
			}
		})
	}
}

func TestHelpSucceedsWithoutLaunchingTarget(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "invoked")
	target := filepath.Join(root, "sentinel")
	if err := os.WriteFile(target, []byte("#!/bin/sh\nprintf invoked >> "+strconv.Quote(marker)+"\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"-h"}, {"-help"}, {"-ch", "1", "-h", target}, {"-ch", "4", "-h", target}, {"-ch", "5", "-h", target}} {
		if got := run(args); got != 0 {
			t.Errorf("run(%q) = %d, want 0", args, got)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("help launched executable: marker error %v", err)
		}
	}
}

func TestChapter5ValidationBeforeLaunch(t *testing.T) {
	dir := t.TempDir()
	module := filepath.Join(dir, "module")
	if err := os.Mkdir(module, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.test/core\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	invalidJSON := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(invalidJSON, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	ran := filepath.Join(dir, "target-ran")
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("#!/bin/sh\nprintf ran > '"+ran+"'\nexit 99\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		flags []string
	}{
		{"missing core", []string{"-ch", "5", "-consumer", module, "-workflow", module}},
		{"missing consumer", []string{"-ch", "5", "-core", module, "-workflow", module}},
		{"missing workflow", []string{"-ch", "5", "-core", module, "-consumer", module}},
		{"workflow not module", []string{"-ch", "5", "-core", module, "-consumer", module, "-workflow", dir}},
		{"workflow missing", []string{"-ch", "5", "-core", module, "-consumer", module, "-workflow", filepath.Join(dir, "absent")}},
		{"bad media JSON", []string{"-ch", "5", "-core", module, "-consumer", module, "-workflow", module, "-media-models", invalidJSON}},
		{"missing media JSON", []string{"-ch", "5", "-core", module, "-consumer", module, "-workflow", module, "-media-models", filepath.Join(dir, "absent.json")}},
		{"bad consumer JSON", []string{"-ch", "5", "-core", module, "-consumer", module, "-workflow", module, "-consumer-case", invalidJSON}},
		{"ch1 workflow", []string{"-ch", "1", "-workflow="}},
		{"ch1 media", []string{"-ch", "1", "-media-models="}},
		{"ch4 workflow", []string{"-ch", "4", "-core", module, "-consumer", module, "-workflow="}},
		{"ch4 media", []string{"-ch", "4", "-core", module, "-consumer", module, "-media-models="}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if code := run(append(tc.flags, target)); code != 2 {
				t.Fatalf("exit=%d, want 2", code)
			}
			if _, err := os.Stat(ran); !os.IsNotExist(err) {
				t.Fatalf("invalid input launched target: %v", err)
			}
		})
	}
}
