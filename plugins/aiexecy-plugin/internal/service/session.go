package service

import (
	"errors"

	"aiexec-plugin/internal/core/plugin_daemon/access_types"
	"aiexec-plugin/internal/core/plugin_manager"
	"aiexec-plugin/internal/core/session_manager"
	"aiexec-plugin/pkg/entities/plugin_entities"
)

func createSession[T any](
	r *plugin_entities.InvokePluginRequest[T],
	access_type access_types.PluginAccessType,
	access_action access_types.PluginAccessAction,
	cluster_id string,
) (*session_manager.Session, error) {
	manager := plugin_manager.Manager()
	if manager == nil {
		return nil, errors.New("failed to get plugin manager")
	}

	// try fetch plugin identifier from plugin id

	runtime, err := manager.Get(r.UniqueIdentifier)
	if err != nil {
		return nil, errors.New("failed to get plugin runtime")
	}

	session := session_manager.NewSession(
		session_manager.NewSessionPayload{
			TenantID:               r.TenantId,
			UserID:                 r.UserId,
			PluginUniqueIdentifier: r.UniqueIdentifier,
			ClusterID:              cluster_id,
			InvokeFrom:             access_type,
			Action:                 access_action,
			Declaration:            runtime.Configuration(),
			BackwardsInvocation:    manager.BackwardsInvocation(),
			IgnoreCache:            false,
			ConversationID:         r.ConversationID,
			MessageID:              r.MessageID,
			AppID:                  r.AppID,
			EndpointID:             r.EndpointID,
		},
	)

	session.BindRuntime(runtime)
	return session, nil
}
