from core.hosting_configuration import HostingConfiguration

hosting_configuration = HostingConfiguration()


from aiexec_app import AiexecApp


def init_app(app: AiexecApp):
    hosting_configuration.init_app(app)
