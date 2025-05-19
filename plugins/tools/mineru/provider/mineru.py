from typing import Any

from aiexec_plugin import ToolProvider
from aiexec_plugin.errors.tool import ToolProviderCredentialValidationError

from tools.parse import MineruTool


class MineruProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict[str, Any]) -> None:
        try:
            instance = MineruTool.from_credentials(credentials)
            assert isinstance(instance, MineruTool)
            instance.validate_token()
            pass
        except Exception as e:
            raise ToolProviderCredentialValidationError(str(e))
