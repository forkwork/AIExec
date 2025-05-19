from collections.abc import Generator

from aiexec_plugin import Tool
from aiexec_plugin.entities.tool import ToolInvokeMessage


class QuestionClassifierTool(Tool):
    def _invoke(self, tool_parameters: dict) -> Generator[ToolInvokeMessage, None, None]:
        return super()._invoke(tool_parameters)
