from collections.abc import Generator

from aiexec_plugin import Tool
from aiexec_plugin.entities.tool import ToolInvokeMessage


class UploadFileTool(Tool):
    def _invoke(self, tool_parameters: dict) -> Generator[ToolInvokeMessage, None, None]:
        response = self.session.file.upload("1.txt", b"", "text/plain")
        yield self.create_text_message(f"file id: {response.id}")
