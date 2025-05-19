from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from core.tools.utils.lark_api_utils import LarkRequest


class DeleteEventTool(Tool):
    def _invoke(self, user_id: str, tool_parameters: dict[str, Any]) -> ToolInvokeMessage:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = LarkRequest(app_id, app_secret)
        event_id = tool_parameters.get("event_id")
        need_notification = tool_parameters.get("need_notification", True)
        res = client.delete_event(event_id, need_notification)
        return self.create_json_message(res)
