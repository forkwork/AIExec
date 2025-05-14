from aiexec_app import AiexecApp


def init_app(app: AiexecApp):
    import flask_migrate  # type: ignore

    from extensions.ext_database import db

    flask_migrate.Migrate(app, db)
