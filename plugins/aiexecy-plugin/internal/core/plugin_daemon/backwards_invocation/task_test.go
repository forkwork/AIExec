package backwards_invocation

import (
	"testing"

	"aiexec-plugin/internal/core/aiexec_invocation"
	"aiexec-plugin/internal/core/aiexec_invocation/tester"
	"aiexec-plugin/internal/core/plugin_daemon/access_types"
	"aiexec-plugin/internal/core/session_manager"
	"aiexec-plugin/pkg/entities/plugin_entities"
)

func getTestSession() *session_manager.Session {
	return session_manager.NewSession(
		session_manager.NewSessionPayload{
			UserID:                 "test",
			TenantID:               "test",
			PluginUniqueIdentifier: plugin_entities.PluginUniqueIdentifier(""),
			ClusterID:              "test",
			InvokeFrom:             access_types.PLUGIN_ACCESS_TYPE_ENDPOINT,
			Action:                 access_types.PLUGIN_ACCESS_ACTION_GET_AI_MODEL_SCHEMAS,
			Declaration:            nil,
			BackwardsInvocation:    tester.NewMockedAiexecInvocation(),
			IgnoreCache:            true,
		},
	)
}

func TestBackwardsInvocationAllPermittedPermission(t *testing.T) {
	allPermittedRuntime := plugin_entities.PluginDeclaration{
		PluginDeclarationWithoutAdvancedFields: plugin_entities.PluginDeclarationWithoutAdvancedFields{
			Resource: plugin_entities.PluginResourceRequirement{
				Permission: &plugin_entities.PluginPermissionRequirement{
					Tool: &plugin_entities.PluginPermissionToolRequirement{
						Enabled: true,
					},
					Model: &plugin_entities.PluginPermissionModelRequirement{
						Enabled:       true,
						LLM:           true,
						TextEmbedding: true,
						Rerank:        true,
						Moderation:    true,
						TTS:           true,
						Speech2text:   true,
					},
					Node: &plugin_entities.PluginPermissionNodeRequirement{
						Enabled: true,
					},
					App: &plugin_entities.PluginPermissionAppRequirement{
						Enabled: true,
					},
				},
			},
		},
	}

	invokeLlmRequest := NewBackwardsInvocation(
		aiexec_invocation.INVOKE_TYPE_LLM,
		"test",
		getTestSession(),
		nil,
		nil,
	)
	if err := checkPermission(&allPermittedRuntime, invokeLlmRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}

	invokeTextEmbeddingRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_TEXT_EMBEDDING, "", getTestSession(), nil, nil)
	if err := checkPermission(&allPermittedRuntime, invokeTextEmbeddingRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}

	invokeRerankRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_RERANK, "", getTestSession(), nil, nil)
	if err := checkPermission(&allPermittedRuntime, invokeRerankRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}

	invokeTtsRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_TTS, "", getTestSession(), nil, nil)
	if err := checkPermission(&allPermittedRuntime, invokeTtsRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}

	invokeSpeech2textRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_SPEECH2TEXT, "", getTestSession(), nil, nil)
	if err := checkPermission(&allPermittedRuntime, invokeSpeech2textRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}

	invokeModerationRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_MODERATION, "", getTestSession(), nil, nil)
	if err := checkPermission(&allPermittedRuntime, invokeModerationRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}

	invokeToolRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_TOOL, "", getTestSession(), nil, nil)
	if err := checkPermission(&allPermittedRuntime, invokeToolRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}

	invokeNodeParameterExtractorRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_NODE_PARAMETER_EXTRACTOR, "", getTestSession(), nil, nil)
	if err := checkPermission(&allPermittedRuntime, invokeNodeParameterExtractorRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}

	invokeNodeQuestionClassifierRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_NODE_QUESTION_CLASSIFIER, "", getTestSession(), nil, nil)
	if err := checkPermission(&allPermittedRuntime, invokeNodeQuestionClassifierRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}

	invokeAppRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_APP, "", getTestSession(), nil, nil)
	if err := checkPermission(&allPermittedRuntime, invokeAppRequest); err != nil {
		t.Errorf("checkPermission failed: %s", err.Error())
	}
}

func TestBackwardsInvocationAllDeniedPermission(t *testing.T) {
	allDeniedRuntime := plugin_entities.PluginDeclaration{
		PluginDeclarationWithoutAdvancedFields: plugin_entities.PluginDeclarationWithoutAdvancedFields{
			Resource: plugin_entities.PluginResourceRequirement{},
		},
	}

	invokeLlmRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_LLM, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeLlmRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}

	invokeTextEmbeddingRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_TEXT_EMBEDDING, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeTextEmbeddingRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}

	invokeRerankRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_RERANK, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeRerankRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}

	invokeTtsRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_TTS, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeTtsRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}

	invokeSpeech2textRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_SPEECH2TEXT, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeSpeech2textRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}

	invokeModerationRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_MODERATION, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeModerationRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}

	invokeToolRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_TOOL, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeToolRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}

	invokeNodeRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_NODE_PARAMETER_EXTRACTOR, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeNodeRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}

	invokeNodeQuestionClassifierRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_NODE_QUESTION_CLASSIFIER, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeNodeQuestionClassifierRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}

	invokeAppRequest := NewBackwardsInvocation(aiexec_invocation.INVOKE_TYPE_APP, "", getTestSession(), nil, nil)
	if err := checkPermission(&allDeniedRuntime, invokeAppRequest); err == nil {
		t.Errorf("checkPermission failed: expected error, got nil")
	}
}
