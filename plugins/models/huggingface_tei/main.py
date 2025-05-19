from aiexec_plugin import Plugin, AiexecPluginEnv

plugin = Plugin(AiexecPluginEnv(MAX_REQUEST_TIMEOUT=120))

if __name__ == "__main__":
    plugin.run()
