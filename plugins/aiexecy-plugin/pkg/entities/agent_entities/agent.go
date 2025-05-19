package agent_entities

import "aiexec-plugin/pkg/entities/tool_entities"

type AgentStrategyResponseChunk struct {
	tool_entities.ToolResponseChunk `json:",inline"`
}
