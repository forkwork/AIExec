from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from core.tools.utils.feishu_api_utils import FeishuRequest


class CreateTaskTool(Tool):
    def _invoke(self, user_id: str, tool_parameters: dict[str, Any]) -> ToolInvokeMessage:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = FeishuRequest(app_id, app_secret)
        summary = tool_parameters.get("summary")
        start_time = tool_parameters.get("start_time")
        end_time = tool_parameters.get("end_time")
        completed_time = tool_parameters.get("completed_time")
        description = tool_parameters.get("description")
        res = client.create_task(summary, start_time, end_time, completed_time, description)
        return self.create_json_message(res)
