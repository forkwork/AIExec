from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from core.tools.utils.lark_api_utils import LarkRequest


class CreateSpreadsheetTool(Tool):
    def _invoke(self, user_id: str, tool_parameters: dict[str, Any]) -> ToolInvokeMessage:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = LarkRequest(app_id, app_secret)
        title = tool_parameters.get("title")
        folder_token = tool_parameters.get("folder_token")
        res = client.create_spreadsheet(title, folder_token)
        return self.create_json_message(res)
