package recall

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

// TestBM25 preserves query-only filtering, positive small-corpus scores and ranking.
func TestBM25(t *testing.T) {
	if got := queryTokens("THE to OAuth/oauth x AB-12"); !reflect.DeepEqual(got, []string{"oauth", "ab", "12"}) {
		t.Fatalf("query tokens: %v", got)
	}
	x := build([]chunk{{content: "the oauth oauth migration"}, {content: "the ordinary gardening"}})
	if x.chunks[0].length != 4 || x.chunks[0].tf["the"] != 1 {
		t.Fatal("document stop words removed")
	}
	hits := x.search("memory", []string{"oauth"}, 0)
	if len(hits) != 1 || hits[0].content != "the oauth oauth migration" || hits[0].score <= 0 {
		t.Fatalf("search: %v", hits)
	}
	if len(x.search("memory", []string{"oauth"}, 100)) != 0 {
		t.Fatal("score floor ignored")
	}
	x = build([]chunk{{content: "oauth migration"}})
	if len(x.search("memory", []string{"oauth"}, 0)) != 1 {
		t.Fatal("small archive lost positive IDF")
	}
}

// TestChunking checks distinct sections, breadcrumbs and paragraph-size boundaries.
func TestChunking(t *testing.T) {
	body := "## Parent\nintro\n### Child\n" + strings.Repeat(strings.Repeat("details ", 80)+"\n\n", 15) + "## Other\nseparate"
	got := chunks("docs/test.md", body)
	if len(got) < 4 || got[len(got)-1].header != "Other" {
		t.Fatalf("sections: %v", got)
	}
	found := false
	for _, c := range got {
		if c.file != "docs/test.md" || len(c.content) > 4096 {
			t.Fatal("lost attribution or paragraph split")
		}
		if c.header == "Parent > Child" {
			found = true
		}
	}
	if !found {
		t.Fatal("breadcrumb missing")
	}
}

// TestFinalQuotas protects memories from higher-scoring documents and redistributes spares.
func TestFinalQuotas(t *testing.T) {
	hits := []hit{{source: "docs", score: 30}, {source: "docs", score: 29}, {source: "docs", score: 28}, {source: "memory", score: 5}, {source: "memory", score: 4}, {source: "memory", score: 3}}
	got := quotas(hits, 3)
	memory := 0
	for _, h := range got {
		if h.source == "memory" {
			memory++
		}
	}
	if len(got) != 3 || memory != 2 {
		t.Fatalf("final quotas: %v", got)
	}
	if got := quotas(hits[:3], 3); len(got) != 3 {
		t.Fatalf("unused memory slots lost: %v", got)
	}
}

// TestJudgeParsing accepts common model formatting without inventing indices.
func TestJudgeParsing(t *testing.T) {
	for _, input := range []string{`[0, 3]`, `Here: ["0", "3"] done`, `[0,"3",null,{},true,"bad"]`} {
		got, err := parseIndices(input)
		if err != nil || !reflect.DeepEqual(got, []int{0, 3}) {
			t.Fatalf("%s: %v %v", input, got, err)
		}
	}
	if got, err := parseIndices("[]"); err != nil || len(got) != 0 {
		t.Fatalf("empty selection: %v %v", got, err)
	}
	if _, err := parseIndices("not json"); err == nil {
		t.Fatal("invalid reply accepted")
	}
}

// TestAttachmentCap counts the complete formatted bytes and preserves UTF-8.
func TestAttachmentCap(t *testing.T) {
	got := format([]hit{{chunk: chunk{file: "memory/a.md", header: "facts", content: strings.Repeat("éclair\n", 2000)}}}, 6144)
	if len(got) != 6144 && len(got) != 6143 {
		t.Fatalf("cap not reached: %d", len(got))
	}
	if !utf8.ValidString(got) || !strings.HasPrefix(got, "[Auto-recalled memories]\n\nFrom memory/a.md — facts:\n> ") || !strings.Contains(got, "\n> éclair") {
		t.Fatal("broken attribution, quoting or Unicode")
	}
	if format([]hit{{chunk: chunk{content: "x"}}}, 5) != "" {
		t.Fatal("header exceeds cap")
	}
}
