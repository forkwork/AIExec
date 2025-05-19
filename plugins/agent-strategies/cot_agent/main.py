import sys

sys.path.append("../..")

from aiexec_plugin import AiexecPluginEnv, Plugin

plugin = Plugin(AiexecPluginEnv(MAX_REQUEST_TIMEOUT=240))

if __name__ == "__main__":
    plugin.run()
