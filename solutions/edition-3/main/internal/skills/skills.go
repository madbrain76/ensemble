package skills

import (
	"encoding/json"
	"ensemble/internal/common"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type skillState struct {
	info     common.SkillInfo
	rendered string
}
type registry struct {
	parent common.Agent
	// mu protects capability truth; no application callback runs under it.
	mu           sync.Mutex
	tools        []common.ToolSpec
	skills       map[string]*skillState
	active       []string
	vars         map[string]func() string
	listeners    map[uint64]func()
	nextListener uint64
}

func New(parent common.Agent, tools []common.ToolSpec) (common.SkillRegistry, error) {
	if parent == nil {
		return nil, fmt.Errorf("skills require Agent parent")
	}
	r := &registry{parent: parent, skills: map[string]*skillState{}, vars: map[string]func() string{}, listeners: map[uint64]func(){}}
	for _, t := range tools {
		if err := r.RegisterTool(t); err != nil {
			return nil, err
		}
	}
	if spec := parent.Config().Skills; spec != nil {
		if spec.Directory == "" || spec.Primary == "" {
			return nil, fmt.Errorf("skills require directory and primary")
		}
		if parent.Config().Tools != nil {
			return nil, fmt.Errorf("Tools and Skills cannot both select startup capabilities")
		}
		if err := r.discover(spec.Directory); err != nil {
			return nil, err
		}
		for name, f := range spec.Variables {
			r.vars[name] = f
		}
	}
	return r, nil
}
func (r *registry) Parent() common.Agent { return r.parent }
func (r *registry) RegisterTool(t common.ToolSpec) error {
	if strings.TrimSpace(t.Name) == "" || strings.TrimSpace(t.Description) == "" || t.Handler == nil {
		return fmt.Errorf("tool requires name, description and handler")
	}
	if len(t.Parameters) == 0 {
		t.Parameters = json.RawMessage(`{"type":"object","properties":{},"additionalProperties":false}`)
	}
	var schema map[string]json.RawMessage
	if json.Unmarshal(t.Parameters, &schema) != nil || schema == nil {
		return fmt.Errorf("tool parameters must be an object schema")
	}
	var schemaType string
	if json.Unmarshal(schema["type"], &schemaType) != nil || schemaType != "object" {
		return fmt.Errorf("tool parameters must have explicit object type")
	}
	t.Parameters = append(json.RawMessage(nil), t.Parameters...)
	if t.Source == "" {
		t.Source = "application"
	}
	if t.Enabled && t.Origin == "" {
		t.Origin = "manual"
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.toolIndex(t.Name) >= 0 {
		return fmt.Errorf("tool %q already registered", t.Name)
	}
	r.tools = append(r.tools, t)
	return nil
}
func (r *registry) toolIndex(name string) int {
	for i, t := range r.tools {
		if t.Name == name {
			return i
		}
	}
	return -1
}
func (r *registry) RegisterVar(name string, f func() string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if f == nil {
		delete(r.vars, name)
	} else {
		r.vars[name] = f
	}
}
func (r *registry) LoadInitial() error {
	cfg := r.parent.Config()
	if cfg.Skills != nil {
		return r.load(cfg.Skills.Primary, "initial", 0)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	names := cfg.Tools
	if names == nil {
		for _, t := range r.tools {
			if t.Name != "load_skill" && t.Name != "unload_skill" {
				names = append(names, t.Name)
			}
		}
	}
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			return fmt.Errorf("duplicate enabled tool %q", name)
		}
		seen[name] = true
		if r.toolIndex(name) < 0 {
			return fmt.Errorf("unknown enabled tool %q", name)
		}
	}
	for _, name := range names {
		t := &r.tools[r.toolIndex(name)]
		t.Enabled = true
		t.Origin = "initial"
	}
	return nil
}
func (r *registry) loadable(enabled map[string]bool) []common.SkillInfo {
	offered := map[string]bool{}
	for name, s := range r.skills {
		if enabled[name] {
			for _, n := range s.info.Offers {
				offered[n] = true
			}
		}
	}
	names := []string{}
	for name := range offered {
		if s := r.skills[name]; s != nil && s.info.Type == "loadable" && !enabled[name] {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	out := []common.SkillInfo{}
	for _, n := range names {
		out = append(out, cloneSkill(r.skills[n].info))
	}
	return out
}
func (r *registry) enabled() map[string]bool {
	out := map[string]bool{}
	for n, s := range r.skills {
		out[n] = s.info.Enabled
	}
	return out
}
func (r *registry) LoadableSkills() []common.SkillInfo {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.loadable(r.enabled())
}
func (r *registry) LoadSkill(name string) error { return r.load(name, "manual", 0) }
func (r *registry) LoadDynamic(name string, seq common.Seq) error {
	if seq == 0 {
		return fmt.Errorf("dynamic skill requires tool_called sequence")
	}
	return r.load(name, "tool", seq)
}
func (r *registry) load(name, origin string, seq common.Seq) error {
	callbacks, err := r.activate(name, origin, seq)
	for _, f := range callbacks {
		f()
	}
	return err
}
func (r *registry) activate(name, origin string, seq common.Seq) ([]func(), error) {
	r.mu.Lock()
	fail := func(err error) ([]func(), error) { r.mu.Unlock(); return nil, err }
	root := r.skills[name]
	if root == nil {
		return fail(fmt.Errorf("unknown skill %q", name))
	}
	if origin == "initial" {
		if root.info.Type != "primary" {
			return fail(fmt.Errorf("initial skill %q is not primary", name))
		}
	} else {
		if root.info.Type != "loadable" {
			return fail(fmt.Errorf("skill %q is not loadable", name))
		}
		if !root.info.Enabled {
			found := false
			for _, s := range r.loadable(r.enabled()) {
				if s.Name == name {
					found = true
				}
			}
			if !found {
				return fail(fmt.Errorf("skill %q is not offered", name))
			}
		}
	}
	order := []string{}
	visiting := map[string]bool{}
	visited := map[string]bool{}
	var visit func(string) error
	visit = func(name string) error {
		s := r.skills[name]
		if s == nil {
			return fmt.Errorf("missing dependency %q", name)
		}
		if visiting[name] {
			return fmt.Errorf("skill dependency cycle at %q", name)
		}
		if visited[name] {
			return nil
		}
		visiting[name] = true
		for _, dep := range s.info.Depends {
			if err := visit(dep); err != nil {
				return err
			}
		}
		for _, tool := range s.info.Tools {
			if r.toolIndex(tool) < 0 {
				return fmt.Errorf("skill %q requires missing tool %q", name, tool)
			}
		}
		visiting[name] = false
		visited[name] = true
		order = append(order, name)
		return nil
	}
	if err := visit(name); err != nil {
		return fail(err)
	}
	vars := map[string]func() string{}
	for n, f := range r.vars {
		vars[n] = f
	}
	r.mu.Unlock()
	// Graph identities are immutable and tool registration is append-only. Resolve
	// application values outside every lock, then publish against the current set.
	values := map[string]string{}
	for n, f := range vars {
		if n != "TOOLS" && n != "SKILLS" && f != nil {
			values[n] = f()
		}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	enabled := r.enabled()
	toolSet := map[string]bool{}
	for _, t := range r.tools {
		if t.Enabled {
			toolSet[t.Name] = true
		}
	}
	if origin == "initial" {
		toolSet["load_skill"] = true
		toolSet["unload_skill"] = true
	}
	for _, n := range order {
		enabled[n] = true
		for _, tool := range r.skills[n].info.Tools {
			toolSet[tool] = true
		}
	}
	tools := []string{}
	for n := range toolSet {
		tools = append(tools, n)
	}
	sort.Strings(tools)
	loadables := []string{}
	for _, s := range r.loadable(enabled) {
		loadables = append(loadables, s.Name+": "+s.Description)
	}
	values["TOOLS"] = strings.Join(tools, ", ")
	values["SKILLS"] = strings.Join(loadables, "\n")
	changed := false
	if origin == "initial" {
		for _, n := range []string{"load_skill", "unload_skill"} {
			t := &r.tools[r.toolIndex(n)]
			t.Enabled = true
			t.Origin = "initial"
		}
	}
	for _, n := range order {
		s := r.skills[n]
		if s.info.PendingUnload {
			s.info.PendingUnload = false
			changed = true
		}
		if s.info.Enabled {
			continue
		}
		s.info.Enabled = true
		s.info.Origin = origin
		s.info.EventSeq = seq
		s.rendered = expand(s.info.Instructions, values)
		r.active = append(r.active, n)
		changed = true
		for _, n := range s.info.Tools {
			t := &r.tools[r.toolIndex(n)]
			if t.Enabled {
				continue
			}
			t.Enabled = true
			t.Skill = s.info.Name
			t.Origin = origin
			t.EventSeq = seq
		}
	}
	return r.callbacks(changed), nil
}
func (r *registry) UnloadSkill(name string) error {
	r.mu.Lock()
	s := r.skills[name]
	if s == nil || !s.info.Enabled {
		r.mu.Unlock()
		return fmt.Errorf("skill %q is not enabled", name)
	}
	changed := !s.info.PendingUnload
	s.info.PendingUnload = true
	callbacks := r.callbacks(changed)
	r.mu.Unlock()
	for _, f := range callbacks {
		f()
	}
	return nil
}
func (r *registry) callbacks(changed bool) []func() {
	var out []func()
	if changed {
		for _, f := range r.listeners {
			out = append(out, f)
		}
	}
	return out
}
func (r *registry) SubscribeSkills(f func()) func() {
	if f == nil {
		return func() {}
	}
	r.mu.Lock()
	r.nextListener++
	id := r.nextListener
	r.listeners[id] = f
	r.mu.Unlock()
	return func() { r.mu.Lock(); delete(r.listeners, id); r.mu.Unlock() }
}
func cloneSkill(s common.SkillInfo) common.SkillInfo {
	s.Tools = append([]string{}, s.Tools...)
	s.Depends = append([]string{}, s.Depends...)
	s.Offers = append([]string{}, s.Offers...)
	return s
}
func (r *registry) ListSkills() []common.SkillInfo { return r.Snapshot().Skills }
func (r *registry) Snapshot() common.CapabilitySnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := common.CapabilitySnapshot{Definitions: []common.ToolDefinition{}, Tools: []common.ToolSpec{}, Skills: []common.SkillInfo{}}
	for _, t := range r.tools {
		t.Parameters = append(json.RawMessage(nil), t.Parameters...)
		out.Tools = append(out.Tools, t)
		if t.Enabled {
			out.Definitions = append(out.Definitions, t.ToolDefinition)
		}
	}
	names := []string{}
	for n := range r.skills {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		out.Skills = append(out.Skills, cloneSkill(r.skills[n].info))
	}
	for _, n := range r.active {
		out.Instructions += "\n\n" + r.skills[n].rendered
	}
	return out
}
func expand(text string, values map[string]string) string {
	var out strings.Builder
	for i := 0; i < len(text); {
		if text[i] != '$' {
			out.WriteByte(text[i])
			i++
			continue
		}
		end := i + 1
		for end < len(text) && (text[end] >= 'A' && text[end] <= 'Z' || text[end] >= 'a' && text[end] <= 'z' || text[end] >= '0' && text[end] <= '9' || text[end] == '_') {
			end++
		}
		name := text[i+1 : end]
		if value, ok := values[name]; ok {
			out.WriteString(value)
		} else {
			out.WriteString(text[i:end])
		}
		i = end
	}
	return out.String()
}
