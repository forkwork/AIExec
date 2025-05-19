from typing import Any
from aiexec_plugin.errors.tool import ToolProviderCredentialValidationError
from tools.dalle2 import DallE2Tool
from aiexec_plugin import ToolProvider


class OpenAIProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict[str, Any]) -> None:
        try:
            for _ in DallE2Tool.from_credentials(credentials, user_id="").invoke(
                tool_parameters={"prompt": "cute girl, blue eyes, white hair, anime style", "size": "small", "n": 1}
            ):
                pass
        except Exception as e:
            raise ToolProviderCredentialValidationError(str(e))
