package gui

import (
	"encoding/json"
	"ensemble"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

// Settings is a complete GUI-local snapshot. False and zero stay on the wire.
// Read Aloud controls playback, not capture; Agent Interrupt remains separate.
type Settings struct {
	Theme                string           `json:"theme"`
	PauseWhileSpeaking   bool             `json:"pause_while_speaking"`
	TTSEnabled           bool             `json:"tts_enabled"`
	TTSSpeed             float64          `json:"tts_speed"`
	FontSize             int              `json:"font_size"`
	LogMode              ensemble.LogMode `json:"log_mode"`
	DiagnosticMaxRecords int              `json:"diagnostic_max_records"`
	MaxToolRounds        int              `json:"max_tool_rounds"`
}

func DefaultSettings() Settings {
	return Settings{Theme: "dark", TTSEnabled: true, PauseWhileSpeaking: true, TTSSpeed: 1, FontSize: 16, LogMode: ensemble.LogFull, DiagnosticMaxRecords: 200}
}

// A sparse patch is decoded onto a copy, then validated as a whole. The same
// path clamps wire and disk input; a type error cannot publish half an update.
func patchSettings(current Settings, raw json.RawMessage) (Settings, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil || fields == nil {
		return current, errors.New("settings must be a JSON object")
	}
	next := current
	for key, value := range fields {
		if string(value) == "null" {
			return current, fmt.Errorf("%s cannot be null", key)
		}
		var target any
		switch key {
		case "theme":
			target = &next.Theme
		case "pause_while_speaking":
			target = &next.PauseWhileSpeaking
		case "tts_enabled":
			target = &next.TTSEnabled
		case "tts_speed":
			target = &next.TTSSpeed
		case "font_size":
			target = &next.FontSize
		case "log_mode":
			target = &next.LogMode
		case "diagnostic_max_records":
			target = &next.DiagnosticMaxRecords
		case "max_tool_rounds":
			target = &next.MaxToolRounds
		default:
			return current, fmt.Errorf("unknown settings field %q", key)
		}
		if err := json.Unmarshal(value, target); err != nil {
			return current, fmt.Errorf("invalid value for %s", key)
		}
	}
	switch next.Theme {
	case "dark", "light", "system":
	default:
		next.Theme = "dark"
	}
	if math.IsNaN(next.TTSSpeed) || math.IsInf(next.TTSSpeed, 0) {
		return current, errors.New("tts_speed must be finite")
	}
	next.TTSSpeed = max(0.5, min(5, next.TTSSpeed))
	next.FontSize = max(12, min(24, next.FontSize))
	next.MaxToolRounds = max(0, next.MaxToolRounds)
	if next.DiagnosticMaxRecords < 0 {
		return current, errors.New("diagnostic_max_records must be nonnegative")
	}
	if next.LogMode != ensemble.LogFull && next.LogMode != ensemble.LogDialogueOnly {
		return current, errors.New("invalid log_mode")
	}
	return next, nil
}

// The hub mutex serializes save, runtime application and publication. This is
// atomic replacement of one file, not a claim of power-loss durability.
type settingsStore struct{ path string }

func (s settingsStore) load() (Settings, error) {
	defaults := DefaultSettings()
	if s.path == "" {
		return defaults, nil
	}
	raw, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return defaults, nil
	}
	if err != nil {
		return defaults, err
	}
	return patchSettings(defaults, raw)
}
func (s settingsStore) save(value Settings) error {
	if s.path == "" {
		return nil
	} // Standalone in-memory hubs have no persistence path.
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(s.path)
	if err = os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, ".settings-*")
	if err != nil {
		return err
	}
	name := file.Name()
	defer os.Remove(name)
	if _, err = file.Write(append(data, '\n')); err != nil {
		file.Close()
		return err
	}
	if err = file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return os.Rename(name, s.path)
}

// Read budget policy back from its owner, not a second GUI runtime budget.
// All Agents in this workbench receive the same saved next-turn policy.
func (h *Hub) currentSettings() Settings {
	value := h.settings
	if agents := h.en.Agents(); len(agents) > 0 {
		cfg := agents[0].Config()
		value.MaxToolRounds = cfg.MaxToolRounds
		value.LogMode = cfg.LogMode
		value.DiagnosticMaxRecords = *cfg.DiagnosticMaxRecords
	}
	return value
}
func (h *Hub) updateSettings(raw json.RawMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	next, err := patchSettings(h.currentSettings(), raw)
	if err != nil {
		return err
	}
	if err = h.store.save(next); err != nil {
		return fmt.Errorf("could not save settings: %w", err)
	}
	// Validation guarantees a nonnegative policy; the synchronous setter publishes
	// it before the success snapshot, while active turns keep their captured limit.
	for _, a := range h.en.Agents() {
		_ = a.SetMaxToolRounds(next.MaxToolRounds)
		_ = a.SetPersistencePolicy(next.LogMode, next.DiagnosticMaxRecords)
	}
	h.settings = next
	h.settingsApplied = true
	settings := h.currentSettings()
	h.broadcast(Wire{Type: "settings_changed", Settings: &settings})
	return nil
}
