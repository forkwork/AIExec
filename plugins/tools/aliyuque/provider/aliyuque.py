from aiexec_plugin.errors.tool import ToolProviderCredentialValidationError
from tools.base import AliYuqueTool
from aiexec_plugin import ToolProvider


class AliYuqueProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict) -> None:
        token = credentials.get("token")
        if not token:
            raise ToolProviderCredentialValidationError("token is required")
        try:
            resp = AliYuqueTool.auth(token)
            if resp and resp.get("data", {}).get("id"):
                return
            raise ToolProviderCredentialValidationError(resp)
        except Exception as e:
            raise ToolProviderCredentialValidationError(str(e))
