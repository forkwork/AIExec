from configs import aiexec_config
from aiexec_app import AiexecApp


def init_app(app: AiexecApp):
    if aiexec_config.RESPECT_XFORWARD_HEADERS_ENABLED:
        from werkzeug.middleware.proxy_fix import ProxyFix

        app.wsgi_app = ProxyFix(app.wsgi_app, x_port=1)  # type: ignore
