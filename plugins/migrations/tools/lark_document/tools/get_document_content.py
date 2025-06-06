from typing import Any
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool
from core.tools.utils.lark_api_utils import LarkRequest


class GetDocumentRawContentTool(Tool):
    def _invoke(self, user_id: str, tool_parameters: dict[str, Any]) -> ToolInvokeMessage:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = LarkRequest(app_id, app_secret)
        document_id = tool_parameters.get("document_id")
        mode = tool_parameters.get("mode", "markdown")
        lang = tool_parameters.get("lang", "0")
        res = client.get_document_content(document_id, mode, lang)
        return self.create_json_message(res)
