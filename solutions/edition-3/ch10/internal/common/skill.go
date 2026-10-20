package common

// Skill is a parsed capability bundle. Source is the winning Markdown path.
// Tools grants concrete tools; Depends activates skills; Offers only advertises.
type Skill struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	Instructions string   `json:"instructions"`
	Tools        []string `json:"tools"`
	Depends      []string `json:"depends"`
	Offers       []string `json:"offers"`
	Source       string   `json:"source"`
}
type SkillInfo struct {
	Skill
	Enabled       bool   `json:"enabled"`
	PendingUnload bool   `json:"pending_unload"`
	Origin        string `json:"origin,omitempty"` // initial, tool, or manual
	EventSeq      Seq    `json:"event_seq"`        // zero for initial/manual, never a guessed history position
}

// Tool joins a declaration with execution and immutable first-enable provenance.
type ToolSpec struct {
	ToolDefinition
	Handler  ToolHandler `json:"-"`
	Skill    string      `json:"skill,omitempty"`
	Source   string      `json:"source"`
	Enabled  bool        `json:"enabled"`
	Origin   string      `json:"origin,omitempty"`
	EventSeq Seq         `json:"event_seq"`
}

// Snapshot is a single registry checkpoint, not independently sampled tables.
type CapabilitySnapshot struct {
	Definitions  []ToolDefinition
	Tools        []ToolSpec
	Skills       []SkillInfo
	Instructions string
}
type SkillSpec struct {
	Directory string
	Primary   string
	Variables map[string]func() string `json:"-"`
}

type SkillRegistry interface {
	Parent() Agent
	Snapshot() CapabilitySnapshot
	RegisterTool(ToolSpec) error
	RegisterVar(string, func() string)
	LoadInitial() error
	LoadDynamic(string, Seq) error
	LoadableSkills() []SkillInfo
	ListSkills() []SkillInfo
	LoadSkill(string) error   // explicit application/manual activation, not a fabricated tool event
	UnloadSkill(string) error // pending only: context and callable tools are retained
	SubscribeSkills(func()) func()
}
