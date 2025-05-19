from aiexec_plugin import Plugin, AiexecPluginEnv
import logging

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

plugin = Plugin(AiexecPluginEnv(MAX_REQUEST_TIMEOUT=120))

if __name__ == '__main__':
    plugin.run()
