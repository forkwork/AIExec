from typing import Any
from aiexec_plugin.errors.tool import ToolProviderCredentialValidationError
from tools.qrcode_generator import QRCodeGeneratorTool
from aiexec_plugin import ToolProvider


class QRCodeProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict[str, Any]) -> None:
        try:
            QRCodeGeneratorTool().invoke(tool_parameters={"content": "Aiexec 123 😊"})
        except Exception as e:
            raise ToolProviderCredentialValidationError(str(e))
