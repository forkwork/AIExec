from aiexec_app import AiexecApp


def init_app(app: AiexecApp):
    from events import event_handlers  # noqa: F401
