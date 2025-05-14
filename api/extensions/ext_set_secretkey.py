from configs import aiexec_config
from aiexec_app import AiexecApp


def init_app(app: AiexecApp):
    app.secret_key = aiexec_config.SECRET_KEY
