from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from core.tools.utils.feishu_api_utils import FeishuRequest


class ListSpreadsheetSheetsTool(Tool):
    def _invoke(self, user_id: str, tool_parameters: dict[str, Any]) -> ToolInvokeMessage:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = FeishuRequest(app_id, app_secret)
        spreadsheet_token = tool_parameters.get("spreadsheet_token")
        res = client.list_spreadsheet_sheets(spreadsheet_token)
        return self.create_json_message(res)
