from aiexec_plugin import ToolProvider
from tools.feishu_api_utils import auth


class FeishuDocumentProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict) -> None:
        auth(credentials)
