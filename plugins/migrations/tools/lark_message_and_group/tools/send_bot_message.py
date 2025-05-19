from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from core.tools.utils.lark_api_utils import LarkRequest


class SendBotMessageTool(Tool):
    def _invoke(self, user_id: str, tool_parameters: dict[str, Any]) -> ToolInvokeMessage:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = LarkRequest(app_id, app_secret)
        receive_id_type = tool_parameters.get("receive_id_type")
        receive_id = tool_parameters.get("receive_id")
        msg_type = tool_parameters.get("msg_type")
        content = tool_parameters.get("content")
        res = client.send_bot_message(receive_id_type, receive_id, msg_type, content)
        return self.create_json_message(res)
