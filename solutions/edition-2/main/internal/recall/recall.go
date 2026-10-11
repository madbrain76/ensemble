package recall

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"

	"ensemble/internal/common"
)

type recaller struct {
	parent  common.Agent
	mu      sync.Mutex
	digest  [32]byte
	indexes map[string]index
}

// New creates one Agent-owned corpus. Failure to read a source degrades recall,
// not construction of the working Agent; each later query retries the scan.
func New(parent common.Agent) common.Recaller {
	if parent == nil {
		panic("Recall requires Agent")
	}
	r := &recaller{parent: parent}
	if _, err := r.search(nil, 0); err != nil {
		parent.Ensemble().Logf("recall index: %v", err)
	}
	return r
}

// Agent exposes the required owner of configuration, files and model services.
func (r *recaller) Agent() common.Agent { return r.parent }

// Validate rejects negative limits; zero retains the documented default.
func Validate(s common.RecallSettings) error {
	if s.MinScore < 0 || s.MaxSnippets < 0 || s.MinQueryLength < 0 || s.Candidates < 0 || s.ContextMsgs < 0 || s.MaxRecallBytes < 0 {
		return fmt.Errorf("recall limits must be nonnegative")
	}
	return nil
}
func settings(s common.RecallSettings) common.RecallSettings {
	if s.MinScore == 0 {
		s.MinScore = 3
	}
	if s.MaxSnippets == 0 {
		s.MaxSnippets = 3
	}
	if s.MinQueryLength == 0 {
		s.MinQueryLength = 80
	}
	if s.Candidates == 0 {
		s.Candidates = 20
	}
	if s.ContextMsgs == 0 {
		s.ContextMsgs = 3
	}
	if s.MaxRecallBytes == 0 {
		s.MaxRecallBytes = 6144
	}
	return s
}

// Corpus files belong to the Agent, not to cwd or another application's home.
// Hash the bytes already read so human edits (even equal-size edits) invalidate
// indexes. Only changed corpora incur tokenization; there is no watcher lifetime.
func (r *recaller) search(query []string, floor float64) ([]hit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	root := r.parent.Config().DataDir
	sources := []struct{ name, path string }{
		{"memory", filepath.Join(root, "memory")},
		{"memory", filepath.Join(root, "bucket-0")},
		{"memory", filepath.Join(root, "bucket-1")},
		{"memory", filepath.Join(root, "MEMORY.md")},
		{"memory", filepath.Join(root, "SOUL.md")},
		{"docs", filepath.Join(root, "docs")},
		{"handoffs", filepath.Join(root, "handoffs")},
		{"learnings", filepath.Join(root, "learnings")},
	}
	if skills := r.parent.Skills(); skills != nil {
		sources = append(sources, struct{ name, path string }{"skills", skills.Directory()})
	}
	corpus := map[string][]chunk{}
	hash := sha256.New()
	// Keep the file texts until we know whether rebuilding is necessary.
	type document struct{ source, file, text string }
	var documents []document
	for _, source := range sources {
		err := filepath.WalkDir(source.path, func(path string, entry fs.DirEntry, err error) error {
			if os.IsNotExist(err) {
				return nil
			}
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".md" && !(source.name == "memory" && strings.HasSuffix(path, ".md.retired")) || source.name == "skills" && entry.Name() != "SKILL.md" {
				return nil
			}
			body, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			name, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			file := name
			if source.name == "skills" {
				name, _ = filepath.Rel(source.path, path)
				file = filepath.Join("skills", name)
			}
			fmt.Fprintf(hash, "%s\x00%d\x00%s\x00", file, len(body), body)
			documents = append(documents, document{source.name, file, string(body)})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	var digest [32]byte
	copy(digest[:], hash.Sum(nil))
	if r.indexes == nil || digest != r.digest {
		for _, doc := range documents {
			corpus[doc.source] = append(corpus[doc.source], chunks(doc.file, doc.text)...)
		}
		r.indexes = map[string]index{}
		for source, docs := range corpus {
			r.indexes[source] = build(docs)
		}
		r.digest = digest
	}
	var hits []hit
	// Fixed order makes equal scores deterministic without filesystem/map lottery.
	for _, source := range []string{"memory", "docs", "handoffs", "learnings", "skills"} {
		hits = append(hits, r.indexes[source].search(source, query, floor)...)
	}
	rank(hits)
	return hits, nil
}

// Recall returns immutable attachment bytes. The caller passes a snapshot so
// the stateless judge can run off the actor without borrowing live conversation.
func (r *recaller) Recall(ctx context.Context, query string, conversation *common.Context) (string, error) {
	s := settings(r.parent.Config().Recall)
	if s.Disabled || utf8.RuneCountInString(query) < s.MinQueryLength {
		return "", nil
	}
	floor := s.MinScore
	if !s.DisableJudge {
		floor = 0.5
	}
	hits, err := r.search(queryTokens(query), floor)
	if err != nil {
		return "", err
	}
	// Include text nested in tool results and blockquoted prior attachments.
	// Normalizing quote prefixes permits comparison with the original snippet.
	var contextText strings.Builder
	contextText.WriteString(r.parent.Engine().SystemPrompt())
	var collect func([]common.Part)
	collect = func(parts []common.Part) {
		for _, p := range parts {
			contextText.WriteByte('\n')
			contextText.WriteString(p.Text)
			collect(p.Parts)
		}
	}
	for _, e := range conversation.Dialogue {
		collect(e.Parts)
	}
	collect(conversation.Ephemera)
	present := unquote(contextText.String())
	available := hits[:0]
	for _, h := range hits {
		if !strings.Contains(present, unquote(h.content)) {
			available = append(available, h)
		}
	}
	hits = available
	if len(hits) == 0 {
		return "", nil
	}
	var chosen []hit
	judged := false
	if !s.DisableJudge {
		candidates := hits[:min(len(hits), s.Candidates)]
		chosen, err = r.judge(ctx, query, conversation, candidates, s)
		judged = err == nil
		if err != nil {
			r.parent.Ensemble().Logf("recall judge failed; using source quotas: %v", err)
		}
	}
	if !judged {
		chosen = quotas(hits, s.MaxSnippets)
	}
	return format(chosen, s.MaxRecallBytes), nil
}
func unquote(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, "> ")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// Quotas apply to the FINAL mechanical selection. A successful judge never
// comes through here, including a valid empty selection. Spare slots flow to
// the best remaining scores, so a sparse source does not waste the budget.
func quotas(hits []hit, limit int) []hit {
	quota := map[string]int{"memory": (limit + 1) / 2}
	var sources []string
	for _, name := range []string{"docs", "handoffs", "learnings", "skills"} {
		for _, h := range hits {
			if h.source == name {
				sources = append(sources, name)
				break
			}
		}
	}
	remaining := limit - quota["memory"]
	for i, name := range sources {
		quota[name] = remaining / len(sources)
		if i < remaining%len(sources) {
			quota[name]++
		}
	}
	selected := make([]bool, len(hits))
	var out []hit
	for i, h := range hits {
		if quota[h.source] > 0 {
			out = append(out, h)
			selected[i] = true
			quota[h.source]--
		}
	}
	for i, h := range hits {
		if len(out) < limit && !selected[i] {
			out = append(out, h)
		}
	}
	return out
}
func format(hits []hit, limit int) string {
	const title = "[Auto-recalled memories]\n"
	out := title
	for _, h := range hits {
		prefix := fmt.Sprintf("\nFrom %s — %s:\n> ", h.file, h.header)
		room := limit - len(out) - len(prefix)
		if room <= 0 {
			break
		}
		body := strings.ReplaceAll(h.content, "\n", "\n> ")
		if len(body) > room {
			body = body[:room]
			for !utf8.ValidString(body) {
				body = body[:len(body)-1]
			}
		}
		out += prefix + body
	}
	if out == title {
		return ""
	}
	return out
}
