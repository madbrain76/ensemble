package history

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"ensemble/internal/common"
)

// SaveFile captures one boundary under History's lock. Encoding and validation
// finish before touching either recovery candidate; filesystem work holds no
// append lock. Agent serializes saves to this path.
func (h *history) SaveFile(path string, config common.SavedConfig, mode common.LogMode) error {
	if mode == "" {
		mode = common.LogFull
	}
	h.mu.Lock()
	s := common.SaveFile{Config: config, LogMode: mode, AsOf: h.through, Context: &h.context}
	s.Log = []common.Event{}
	s.DialogueLog = append([]common.DialogueRecord(nil), h.dialogue...)
	if mode == common.LogFull {
		s.Log = h.events
	} else {
		s.DialogueLog = dialogueArchive(s.DialogueLog, h.events)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	h.mu.Unlock()
	if err != nil {
		return fmt.Errorf("encoding save: %w", err)
	}
	if _, _, err = decodeSave(data); err != nil {
		return fmt.Errorf("validating save: %w", err)
	}

	prior, primaryErr := readSave(path)
	if primaryErr == nil {
		// Never rotate an unvalidated primary over a good backup.
		if err = replaceFile(path+".bak", prior); err != nil {
			return fmt.Errorf("saving backup: %w", err)
		}
	} else {
		_, backupErr := readSave(path + ".bak")
		if !(errors.Is(primaryErr, os.ErrNotExist) && errors.Is(backupErr, os.ErrNotExist)) && backupErr != nil {
			return fmt.Errorf("refusing save: neither recovery candidate is valid (primary: %v; backup: %v)", primaryErr, backupErr)
		}
	}
	// The archive has been encoded/validated and recovery candidates checked.
	// Only now establish a caller-selected or isolated default directory.
	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	if err = replaceFile(path, append(data, '\n')); err != nil {
		return fmt.Errorf("saving primary: %w", err)
	}
	return nil
}

// replaceFile leaves the old destination untouched until a complete, closed
// sibling is ready. Sync improves ordinary durability; this is not a promise
// about every filesystem or power-loss ordering of directory entries.
func replaceFile(path string, data []byte) (err error) {
	f, err := os.CreateTemp(filepath.Dir(path), ".ensemble-save-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	_, err = f.Write(data)
	if err == nil {
		err = f.Sync()
	}
	err = errors.Join(err, f.Close())
	if err == nil {
		err = os.Rename(f.Name(), path)
	}
	return err
}

func readSave(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	_, _, err = decodeSave(data)
	return data, err
}

// LoadFile does not publish partial state. A backup is decoded independently,
// and its own boundary is reported, never the lost primary's later boundary.
func (h *history) LoadFile(path string) (common.Recovery, error) {
	var primaryErr error
	for i, candidate := range []string{path, path + ".bak"} {
		data, err := os.ReadFile(candidate)
		var s common.SaveFile
		var context common.Context
		if err == nil {
			s, context, err = decodeSave(data)
		}
		if err == nil {
			through := s.AsOf
			if len(s.Log) > 0 && s.Log[len(s.Log)-1].Seq > through {
				through = s.Log[len(s.Log)-1].Seq
			}
			h.mu.Lock()
			h.events, h.context, h.through, h.dialogue = s.Log, context, through, s.DialogueLog
			h.mu.Unlock()
			return common.Recovery{Path: candidate, AsOf: through, Backup: i == 1}, nil
		}
		if i == 0 {
			primaryErr = err
			continue
		}
		if errors.Is(primaryErr, os.ErrNotExist) && errors.Is(err, os.ErrNotExist) {
			return common.Recovery{}, nil
		}
		return common.Recovery{}, fmt.Errorf("cannot recover history: primary: %v; backup: %v", primaryErr, err)
	}
	panic("unreachable recovery")
}

// Strict decoding rejects silent data loss from unknown fields or malformed
// tagged payloads. A syntactically valid object alone is not a recovery point.
func decodeSave(data []byte) (common.SaveFile, common.Context, error) {
	var s common.SaveFile
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return s, common.Context{}, err
	}
	for _, key := range []string{"config", "as_of"} {
		if len(fields[key]) == 0 || bytes.Equal(fields[key], []byte("null")) {
			return s, common.Context{}, fmt.Errorf("missing %s", key)
		}
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&s); err != nil {
		return s, common.Context{}, err
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return s, common.Context{}, fmt.Errorf("trailing save data")
	}
	if s.LogMode == "" {
		s.LogMode = common.LogFull
	}
	if s.LogMode != common.LogFull && s.LogMode != common.LogDialogueOnly {
		return s, common.Context{}, fmt.Errorf("unknown log mode")
	}
	var configFields map[string]json.RawMessage
	if err := json.Unmarshal(fields["config"], &configFields); err != nil {
		return s, common.Context{}, err
	}
	for _, key := range []string{"vendor", "model", "system_prompt", "tools"} {
		if _, ok := configFields[key]; !ok {
			return s, common.Context{}, fmt.Errorf("missing config.%s", key)
		}
	}
	if err := savedConfig(s.Config); err != nil {
		return s, common.Context{}, err
	}
	if _, ok := fields["context"]; !ok {
		return s, common.Context{}, fmt.Errorf("missing context")
	}
	if _, ok := fields["log"]; !ok {
		return s, common.Context{}, fmt.Errorf("missing log")
	}
	if s.Context != nil {
		if err := contextValid(*s.Context, s.AsOf); err != nil {
			return s, common.Context{}, err
		}
	}
	if s.Context == nil && s.LogMode == common.LogDialogueOnly {
		return s, common.Context{}, fmt.Errorf("dialogue_only requires context")
	}
	var previous common.Seq
	for _, r := range s.DialogueLog {
		if r.SourceSeq == 0 || r.SourceSeq < previous || r.SourceSeq > s.AsOf || (r.Role != "user" && r.Role != "assistant") || r.Text == "" {
			return s, common.Context{}, fmt.Errorf("invalid dialogue archive")
		}
		previous = r.SourceSeq
	}

	var last common.Seq
	prefix := empty()
	complete := true
	context := empty()
	if s.Context != nil {
		context = clone(*s.Context)
	}
	for _, e := range s.Log {
		if e.Seq == 0 || e.Seq <= last || e.Time.IsZero() {
			return s, common.Context{}, fmt.Errorf("invalid event sequence/time")
		}
		if err := eventValid(e); err != nil {
			return s, common.Context{}, fmt.Errorf("event %d: %w", e.Seq, err)
		}
		if e.Redact != nil && (e.Redact.From == 0 || e.Redact.From > e.Redact.To || e.Redact.To >= e.Seq) {
			return s, common.Context{}, fmt.Errorf("invalid redaction boundary")
		}
		if s.Context == nil && e.Seq != last+1 {
			return s, common.Context{}, fmt.Errorf("incomplete canonical log cannot rebuild context")
		}
		if s.LogMode == common.LogDialogueOnly && e.Seq <= s.AsOf {
			return s, common.Context{}, fmt.Errorf("dialogue_only log must contain only recovery tail")
		}
		if e.Seq <= s.AsOf {
			if e.Seq != last+1 {
				complete = false
			}
			Apply(&prefix, e)
			if e.Seq == s.AsOf && complete && s.Context != nil && !sameContext(prefix, *s.Context) {
				return s, common.Context{}, fmt.Errorf("context disagrees with canonical boundary")
			}
		}
		last = e.Seq
		if s.Context == nil || e.Seq > s.AsOf {
			Apply(&context, e)
		}
	}
	if s.Context == nil && (last < s.AsOf || len(s.DialogueLog) > 0) {
		return s, common.Context{}, fmt.Errorf("incomplete canonical history cannot rebuild context")
	}
	if err := contextValid(context, max(last, s.AsOf)); err != nil {
		return s, common.Context{}, err
	}
	return s, context, nil
}

func savedConfig(c common.SavedConfig) error {
	if strings.TrimSpace(c.Model) == "" || c.MaxToolRounds < 0 {
		return fmt.Errorf("invalid saved configuration")
	}
	if c.Surface == 0 {
		switch c.Vendor {
		case common.Anthropic:
			c.Surface = common.Messages
		case common.OpenAI:
			c.Surface = common.ChatCompletions
		case common.Gemini:
			c.Surface = common.GenerateContent
		}
	}
	if !((c.Vendor == common.Anthropic && c.Surface == common.Messages) || (c.Vendor == common.OpenAI && c.Surface == common.ChatCompletions) || (c.Vendor == common.Gemini && c.Surface == common.GenerateContent)) {
		return fmt.Errorf("invalid saved provider surface")
	}
	switch c.StreamCapability {
	case common.StreamUnknown, common.StreamPlain, common.StreamSSE:
	default:
		return fmt.Errorf("invalid stream capability")
	}
	for _, t := range c.Tools {
		if t.Name == "" || !objectJSON(t.Schema) {
			return fmt.Errorf("invalid saved tool")
		}
	}
	return nil
}

func contextValid(c common.Context, through common.Seq) error {
	switch c.Turn {
	case common.Idle, common.Interrupted, common.InputPending, common.InFlight, common.ToolsPending:
	default:
		return fmt.Errorf("invalid context turn")
	}
	if c.Dialogue == nil {
		return fmt.Errorf("missing dialogue")
	}
	if c.Usage.Input < 0 || c.Usage.Output < 0 || c.Usage.CacheRead < 0 || c.Usage.CacheWrite < 0 {
		return fmt.Errorf("negative usage")
	}
	var last common.Seq
	for _, entry := range c.Dialogue {
		if entry.Seq == 0 || entry.Seq < last || entry.Seq > through {
			return fmt.Errorf("invalid dialogue boundary")
		}
		last = entry.Seq
		if entry.Actor != common.Human && entry.Actor != common.Assistant && entry.Actor != common.System && entry.Actor != common.Tool {
			return fmt.Errorf("invalid dialogue actor")
		}
		if err := parts(entry.Parts); err != nil {
			return err
		}
		if err := recoveryParts(entry.Parts); err != nil {
			return err
		}
	}
	if err := parts(c.Ephemera); err != nil {
		return err
	}
	return recoveryParts(c.Ephemera)
}

// Only visible human/assistant text belongs to the derived archive. Existing
// records survive future saves; raw tool/opaque data stays in recovery context.
func dialogueArchive(prior []common.DialogueRecord, events []common.Event) []common.DialogueRecord {
	out := append([]common.DialogueRecord(nil), prior...)
	var through common.Seq
	if len(prior) > 0 {
		through = prior[len(prior)-1].SourceSeq
	}
	for _, e := range events {
		if e.Seq <= through {
			continue
		}
		var role string
		var ps common.PartList
		if e.Message != nil && ((e.Type == common.MessageReceived && (e.Message.Actor == common.Human || e.Message.Actor == common.Assistant)) || (e.Type == common.HintReceived && e.Message.Actor == common.Human)) {
			role, ps = "user", e.Message.Parts
			if e.Message.Actor == common.Assistant {
				role = "assistant"
			}
		}
		if e.Type == common.ResponseEnded && !e.Stale && e.Response != nil {
			role, ps = "assistant", e.Response.Parts
		}
		var text strings.Builder
		for _, p := range ps {
			if p.Type == "text" {
				text.WriteString(p.Text)
			}
		}
		if role != "" && text.Len() > 0 {
			out = append(out, common.DialogueRecord{SourceSeq: e.Seq, Role: role, Text: text.String()})
		}
	}
	return out
}

func sameContext(a, b common.Context) bool {
	// Compare serialized values, not RawMessage whitespace/key order or slice
	// allocation. The two required collection fields do not use omitempty;
	// other empty collections normalize through their existing JSON tags.
	for _, c := range []*common.Context{&a, &b} {
		if len(c.Dialogue) == 0 {
			c.Dialogue = nil
		}
		if len(c.Ephemera) == 0 {
			c.Ephemera = nil
		}
	}
	value := func(c common.Context) (any, error) {
		raw, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}
		dec := json.NewDecoder(bytes.NewReader(raw))
		// Do not round provider arguments/opaque integers through float64.
		dec.UseNumber()
		var v any
		err = dec.Decode(&v)
		return v, err
	}
	x, err := value(a)
	if err != nil {
		return false
	}
	y, err := value(b)
	return err == nil && reflect.DeepEqual(x, y)
}

func objectJSON(raw json.RawMessage) bool {
	var v map[string]json.RawMessage
	return json.Unmarshal(raw, &v) == nil && v != nil
}
func recoveryParts(ps common.PartList) error {
	for _, p := range ps {
		if len(p.Args) > 0 && !json.Valid(p.Args) || len(p.Data) > 0 && !json.Valid(p.Data) {
			return fmt.Errorf("invalid part JSON")
		}
		if p.Type == "tool_call" && !objectJSON(p.Args) {
			return fmt.Errorf("tool call requires argument object")
		}
		if p.Type == "opaque" && (p.From == nil || p.From.Vendor == 0 || p.From.Surface == 0 || p.From.Model == "") {
			return fmt.Errorf("opaque part requires provenance")
		}
		if err := recoveryParts(p.Parts); err != nil {
			return err
		}
	}
	return nil
}
func eventValid(e common.Event) error {
	if err := validate(e); err != nil {
		return err
	}
	var ps common.PartList
	if e.Message != nil {
		ps = e.Message.Parts
	}
	if e.Response != nil {
		ps = e.Response.Parts
	}
	if e.Tool != nil {
		ps = e.Tool.Parts
		// Direct public Dispatch operations belong to no model turn and have
		// no native call ID. Preserve that fact; model-turn calls still require
		// their supplied correlation ID, and all other validation remains.
		if e.Tool.CallID == "" && e.TurnID != 0 {
			return fmt.Errorf("missing tool call ID")
		}
		if len(e.Tool.Args) > 0 && !objectJSON(e.Tool.Args) {
			return fmt.Errorf("invalid tool args")
		}
	}
	if e.Redact != nil {
		ps = e.Redact.Replacement
	}
	return recoveryParts(ps)
}
