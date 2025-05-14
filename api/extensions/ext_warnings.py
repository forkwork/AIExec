from aiexec_app import AiexecApp


def init_app(app: AiexecApp):
    import warnings

    warnings.simplefilter("ignore", ResourceWarning)
