from aiexec_plugin.errors.tool import ToolProviderCredentialValidationError
from tools.arxiv_search import ArxivSearchTool
from aiexec_plugin import ToolProvider


class ArxivProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict) -> None:
        try:
            ArxivSearchTool.from_credentials(credentials, user_id="").invoke(tool_parameters={"query": "John Doe"})
        except Exception as e:
            raise ToolProviderCredentialValidationError(str(e))
