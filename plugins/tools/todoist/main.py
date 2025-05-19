from aiexec_plugin import AiexecPluginEnv, Plugin

plugin = Plugin(AiexecPluginEnv(MAX_REQUEST_TIMEOUT=120))

if __name__ == "__main__":
    plugin.run()
