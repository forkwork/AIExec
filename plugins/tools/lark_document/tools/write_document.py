from typing import Any, Generator
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from tools.lark_api_utils import LarkRequest


class CreateDocumentTool(Tool):
    def _invoke(self, tool_parameters: dict[str, Any]) -> Generator[ToolInvokeMessage, None, None]:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = LarkRequest(app_id, app_secret)
        document_id = tool_parameters.get("document_id")
        content = tool_parameters.get("content")
        position = tool_parameters.get("position", "end")
        res = client.write_document(document_id, content, position)
        yield self.create_json_message(res)
