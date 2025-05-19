import logging

from aiexec_plugin.interfaces.model import ModelProvider

logger = logging.getLogger(__name__)


class OAICompatProvider(ModelProvider):
    def validate_provider_credentials(self, credentials: dict) -> None:
        pass
