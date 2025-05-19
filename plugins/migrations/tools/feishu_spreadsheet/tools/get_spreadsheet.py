from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from core.tools.utils.feishu_api_utils import FeishuRequest


class GetSpreadsheetTool(Tool):
    def _invoke(self, user_id: str, tool_parameters: dict[str, Any]) -> ToolInvokeMessage:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = FeishuRequest(app_id, app_secret)
        spreadsheet_token = tool_parameters.get("spreadsheet_token")
        user_id_type = tool_parameters.get("user_id_type", "open_id")
        res = client.get_spreadsheet(spreadsheet_token, user_id_type)
        return self.create_json_message(res)
