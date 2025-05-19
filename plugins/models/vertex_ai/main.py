import grpc.experimental.gevent
from aiexec_plugin import Plugin, AiexecPluginEnv

grpc.experimental.gevent.init_gevent()

plugin = Plugin(AiexecPluginEnv(MAX_REQUEST_TIMEOUT=120))

if __name__ == '__main__':
    plugin.run()
