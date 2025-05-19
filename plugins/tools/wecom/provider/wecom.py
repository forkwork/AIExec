from tools.wecom_group_bot import WecomGroupBotTool
from aiexec_plugin import ToolProvider


class WecomProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict) -> None:
        WecomGroupBotTool()
