package local_runtime

import (
	"aiexec-plugin/internal/core/plugin_daemon/access_types"
	"aiexec-plugin/internal/utils/log"
	"aiexec-plugin/internal/utils/parser"
	"aiexec-plugin/pkg/entities"
	"aiexec-plugin/pkg/entities/plugin_entities"
)

func (r *LocalPluginRuntime) Listen(session_id string) *entities.Broadcast[plugin_entities.SessionMessage] {
	listener := entities.NewBroadcast[plugin_entities.SessionMessage]()
	listener.OnClose(func() {
		r.stdioHolder.removeStdioHandlerListener(session_id)
	})
	r.stdioHolder.setupStdioEventListener(session_id, func(b []byte) {
		// unmarshal the session message
		data, err := parser.UnmarshalJsonBytes[plugin_entities.SessionMessage](b)
		if err != nil {
			log.Error("unmarshal json failed: %s, failed to parse session message", err.Error())
			return
		}

		listener.Send(data)
	})
	return listener
}

func (r *LocalPluginRuntime) Write(session_id string, action access_types.PluginAccessAction, data []byte) {
	r.stdioHolder.write(append(data, '\n'))
}
