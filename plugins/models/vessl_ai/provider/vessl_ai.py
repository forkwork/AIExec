import logging
from aiexec_plugin import ModelProvider

logger = logging.getLogger(__name__)


class VesslAIProvider(ModelProvider):
    def validate_provider_credentials(self, credentials: dict) -> None:
        pass
