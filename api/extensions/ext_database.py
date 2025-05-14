from aiexec_app import AiexecApp
from models import db


def init_app(app: AiexecApp):
    db.init_app(app)
