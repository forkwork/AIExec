import logging
from typing import Any, Generator

import numexpr as ne
from aiexec_plugin import Tool
from aiexec_plugin.entities.tool import ToolInvokeMessage


class EvaluateExpressionTool(Tool):
    def _invoke(
        self, tool_parameters: dict[str, Any]
    ) -> Generator[ToolInvokeMessage, None, None]:
        """
        invoke tools
        """
        expression = tool_parameters.get("expression", "").strip()
        if not expression:
            yield self.create_text_message("Invalid expression")
        try:
            result = ne.evaluate(expression)
            result_str = str(result)
            yield self.create_text_message(f'The result of the expression "{expression}" is {result_str}')
        except Exception as e:
            logging.exception(f"Error evaluating expression: {expression}")
            yield self.create_text_message(f"Invalid expression: {expression}, error: {str(e)}")