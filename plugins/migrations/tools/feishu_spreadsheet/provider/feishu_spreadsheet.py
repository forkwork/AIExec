from aiexec_plugin import ToolProvider
from core.tools.utils.feishu_api_utils import auth


class FeishuMessageProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict) -> None:
        auth(credentials)
