package generic_invoke

import (
	"aiexec-plugin/internal/core/plugin_daemon/access_types"
	"aiexec-plugin/internal/core/session_manager"
	"aiexec-plugin/internal/utils/parser"
)

func getBasicPluginAccessMap(
	user_id string,
	access_type access_types.PluginAccessType,
	action access_types.PluginAccessAction,
) map[string]any {
	return map[string]any{
		"user_id": user_id,
		"type":    access_type,
		"action":  action,
	}
}

func GetInvokePluginMap(
	session *session_manager.Session,
	request any,
) map[string]any {
	req := getBasicPluginAccessMap(
		session.UserID,
		session.InvokeFrom,
		session.Action,
	)
	for k, v := range parser.StructToMap(request) {
		req[k] = v
	}
	return req
}
