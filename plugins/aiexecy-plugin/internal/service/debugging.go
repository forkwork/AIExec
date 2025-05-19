package service

import (
	"aiexec-plugin/internal/core/plugin_manager/debugging_runtime"
	"aiexec-plugin/internal/types/exception"
	"aiexec-plugin/pkg/entities"
)

func GetRemoteDebuggingKey(tenant_id string) *entities.Response {
	type response struct {
		Key string `json:"key"`
	}

	key, err := debugging_runtime.GetConnectionKey(debugging_runtime.ConnectionInfo{
		TenantId: tenant_id,
	})

	if err != nil {
		return exception.InternalServerError(err).ToResponse()
	}

	return entities.NewSuccessResponse(response{
		Key: key,
	})
}
