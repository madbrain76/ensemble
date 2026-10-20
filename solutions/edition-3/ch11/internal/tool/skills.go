package tool

import (
	"encoding/json"
	"ensemble/internal/common"
	"fmt"
	"io"
)

func (r *registry) skillTools() []entry {
	schema := json.RawMessage(`{"type":"object","properties":{"name":{"type":"string"}},"required":["name"],"additionalProperties":false}`)
	return []entry{
		{definition: common.ToolDefinition{Name: "load_skill", Description: "Enable a named skill and its dependencies atomically. Offers only advertise; tools and instructions become available after a successful load.", Parameters: schema}},
		{definition: common.ToolDefinition{Name: "unload_skill", Description: "Mark an enabled skill pending unload. Its instructions and tools remain available until a future supported context-compaction boundary.", Parameters: schema}},
	}
}
func (r *registry) skillCall(owner common.ToolCall, tool string, args arguments, seq common.Seq) error {
	if err := args.only("name"); err != nil {
		return err
	}
	name, err := args.text("name", true)
	if err != nil {
		return err
	}
	if tool == "unload_skill" {
		if err = r.parent.Skills().UnloadSkill(name); err != nil {
			return err
		}
		_, err = fmt.Fprintf(owner, "Skill %s is pending unload; its instructions and tools remain active.", name)
		return err
	}
	if seq == 0 {
		return fmt.Errorf("load_skill requires recorded tool admission")
	}
	err = r.parent.Skills().LoadDynamic(name, seq)
	if err != nil {
		return err
	}
	_, err = io.WriteString(owner, "Enabled skills: "+name+". Instructions are supplied through the provider system prompt; this result is status, not an instruction channel.")
	return err
}
