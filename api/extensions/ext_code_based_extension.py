from core.extension.extension import Extension
from aiexec_app import AiexecApp


def init_app(app: AiexecApp):
    code_based_extension.init()


code_based_extension = Extension()
