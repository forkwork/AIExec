package integration_test

import (
	_ "embed"
	"testing"

	"aiexec-plugin/internal/core/plugin_daemon/access_types"
	"aiexec-plugin/internal/core/plugin_manager/test_utils"
	"aiexec-plugin/internal/utils/parser"
	"aiexec-plugin/internal/utils/routine"
	"aiexec-plugin/pkg/entities/agent_entities"
	"aiexec-plugin/pkg/entities/requests"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/khulnasoft-agent_0.0.14.aiexecpkg
var aiexecOfficialAgent []byte

//go:embed testdata/invoke_agent_strategy_json.json
var invokeAgentStrategyJson []byte

const (
	_testingPath = "./integration_test_cwd"
)

func TestAiexecOfficialAgentIntegration(t *testing.T) {
	routine.InitPool(10000)

	defer test_utils.ClearTestingPath(_testingPath)

	runtime, err := test_utils.GetRuntime(aiexecOfficialAgent, _testingPath)
	assert.NoError(t, err)

	invokePayload, err := parser.UnmarshalJsonBytes2Map(invokeAgentStrategyJson)
	assert.NoError(t, err)
	response, err := test_utils.RunOnce[requests.RequestInvokeAgentStrategy, agent_entities.AgentStrategyResponseChunk](
		runtime,
		access_types.PLUGIN_ACCESS_TYPE_AGENT_STRATEGY,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_AGENT_STRATEGY,
		requests.RequestInvokeAgentStrategy{
			InvokeAgentStrategySchema: requests.InvokeAgentStrategySchema{
				AgentStrategyProvider: "agent",
				AgentStrategy:         "function_calling",
				AgentStrategyParams:   invokePayload,
			},
		},
	)

	assert.NoError(t, err)

	for response.Next() {
		_, err := response.Read()
		assert.NoError(t, err)
	}
}
