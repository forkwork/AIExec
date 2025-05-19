from collections.abc import Generator

from aiexec_plugin import Tool
from aiexec_plugin.entities.tool import ToolInvokeMessage


class ToolTool(Tool):
    def _invoke(self, tool_parameters: dict) -> Generator[ToolInvokeMessage, None, None]:
        yield self.create_image_message("https://assets.aiexec.ai/images/aiexec_logo_dark_s.png")
