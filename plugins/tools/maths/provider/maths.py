from typing import Any
from aiexec_plugin.errors.tool import ToolProviderCredentialValidationError
from tools.eval_expression import EvaluateExpressionTool
from aiexec_plugin import ToolProvider


class MathsProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict[str, Any]) -> None:
        try:
            EvaluateExpressionTool().invoke(tool_parameters={"expression": "1+(2+3)*4"})
        except Exception as e:
            raise ToolProviderCredentialValidationError(str(e))
