package engine

import "ensemble/internal/common"

// LookupModel has no fallback; model capability is independent of wire support.
func LookupModel(model string) (common.ModelFeatures, bool) {
	rows := map[string]common.Media{"claude-sonnet-4-6": common.MediaImage | common.MediaDocument, "gpt-4.1-2025-04-14": common.MediaImage, "gemini-3.8-flash": common.MediaImage | common.MediaAudio | common.MediaVideo | common.MediaDocument}
	m, ok := rows[model]
	return common.ModelFeatures{Media: m}, ok
}

// ModelCapabilities requires an explicit matching provider and model row.
func ModelCapabilities(v common.Vendor, model string) common.MediaCapabilities {
	rows := map[common.Vendor]string{common.Anthropic: "claude-sonnet-4-6", common.OpenAI: "gpt-4.1-2025-04-14", common.Gemini: "gemini-3.8-flash"}
	if rows[v] != model {
		return common.MediaCapabilities{}
	}
	f, _ := LookupModel(model)
	return common.MediaCapabilities{Images: f.Media&common.MediaImage != 0, Media: f.Media}
}
