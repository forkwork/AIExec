from collections.abc import Generator
from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from tools.feishu_api_utils import FeishuRequest


class CreateDocumentTool(Tool):
    def _invoke(self, tool_parameters: dict[str, Any]) -> Generator[ToolInvokeMessage, None, None]:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = FeishuRequest(app_id, app_secret)
        title = tool_parameters.get("title")
        content = tool_parameters.get("content")
        folder_token = tool_parameters.get("folder_token")
        res = client.create_document(title, content, folder_token)
        yield self.create_json_message(res)
