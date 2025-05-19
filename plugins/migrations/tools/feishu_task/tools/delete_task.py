from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from core.tools.utils.feishu_api_utils import FeishuRequest


class UpdateTaskTool(Tool):
    def _invoke(self, user_id: str, tool_parameters: dict[str, Any]) -> ToolInvokeMessage:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = FeishuRequest(app_id, app_secret)
        task_guid = tool_parameters.get("task_guid")
        res = client.delete_task(task_guid)
        return self.create_json_message(res)
