import logging

from aiexec_plugin.interfaces.model import ModelProvider

logger = logging.getLogger(__name__)


class SparkProvider(ModelProvider):
    def validate_provider_credentials(self, credentials: dict) -> None:
        """
        Validate provider credentials

        if validate failed, raise exception

        :param credentials: provider credentials, credentials form defined in `provider_credential_schema`.
        """
        # ignore credentials validation because every model has their own spark quota pool
        pass
