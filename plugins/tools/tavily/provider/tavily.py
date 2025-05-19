from typing import Any
from aiexec_plugin.errors.tool import ToolProviderCredentialValidationError
from tools.tavily_search import TavilySearchTool
from aiexec_plugin import ToolProvider


class TavilyProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict[str, Any]) -> None:
        try:
            TavilySearchTool.from_credentials(credentials, user_id="").invoke(
                tool_parameters={
                    "query": "Aiexec AI",
                }
            )
        except Exception as e:
            raise ToolProviderCredentialValidationError(str(e))
