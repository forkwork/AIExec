from typing import Any
from aiexec_plugin import ToolProvider


class HapProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict[str, Any]) -> None:
        pass
