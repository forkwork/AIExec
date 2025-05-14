import os
import time

from aiexec_app import AiexecApp


def init_app(app: AiexecApp):
    os.environ["TZ"] = "UTC"
    # windows platform not support tzset
    if hasattr(time, "tzset"):
        time.tzset()
