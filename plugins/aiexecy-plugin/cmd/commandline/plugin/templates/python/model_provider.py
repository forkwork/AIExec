import logging
from collections.abc import Mapping

from aiexec_plugin import ModelProvider
from aiexec_plugin.entities.model import ModelType
from aiexec_plugin.errors.model import CredentialsValidateFailedError

logger = logging.getLogger(__name__)


class {{ .PluginName | SnakeToCamel }}ModelProvider(ModelProvider):
    def validate_provider_credentials(self, credentials: Mapping) -> None:
        """
        Validate provider credentials
        if validate failed, raise exception

        :param credentials: provider credentials, credentials form defined in `provider_credential_schema`.
        """
        try:
            pass
        except CredentialsValidateFailedError as ex:
            raise ex
        except Exception as ex:
            logger.exception(
                f"{self.get_provider_schema().provider} credentials validate failed"
            )
            raise ex
