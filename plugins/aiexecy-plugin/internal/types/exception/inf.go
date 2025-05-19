package exception

import "aiexec-plugin/pkg/entities"

type PluginDaemonError interface {
	error

	ToResponse() *entities.Response
}
