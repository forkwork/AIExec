from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from core.tools.utils.lark_api_utils import LarkRequest


class CreateDocumentTool(Tool):
    def _invoke(self, user_id: str, tool_parameters: dict[str, Any]) -> ToolInvokeMessage:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = LarkRequest(app_id, app_secret)
        document_id = tool_parameters.get("document_id")
        content = tool_parameters.get("content")
        position = tool_parameters.get("position", "end")
        res = client.write_document(document_id, content, position)
        return self.create_json_message(res)
