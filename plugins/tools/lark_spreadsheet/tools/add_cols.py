from typing import Any, Generator
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from tools.lark_api_utils import LarkRequest


class AddColsTool(Tool):
    def _invoke(self, tool_parameters: dict[str, Any]) -> Generator[ToolInvokeMessage, None, None]:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = LarkRequest(app_id, app_secret)
        spreadsheet_token = tool_parameters.get("spreadsheet_token")
        sheet_id = tool_parameters.get("sheet_id")
        sheet_name = tool_parameters.get("sheet_name")
        length = tool_parameters.get("length")
        values = tool_parameters.get("values")
        res = client.add_cols(spreadsheet_token, sheet_id, sheet_name, length, values)
        yield self.create_json_message(res)
