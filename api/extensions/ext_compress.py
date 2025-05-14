from configs import aiexec_config
from aiexec_app import AiexecApp


def is_enabled() -> bool:
    return aiexec_config.API_COMPRESSION_ENABLED


def init_app(app: AiexecApp):
    from flask_compress import Compress  # type: ignore

    compress = Compress()
    compress.init_app(app)
