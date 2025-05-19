import logging

from aiexec_plugin.entities.model import ModelType
from aiexec_plugin.errors.model import CredentialsValidateFailedError
from aiexec_plugin.interfaces.model import ModelProvider

logger = logging.getLogger(__name__)


class VoyageProvider(ModelProvider):
    def validate_provider_credentials(self, credentials: dict) -> None:
        """
        Validate provider credentials
        if validate failed, raise exception

        :param credentials: provider credentials, credentials form defined in `provider_credential_schema`.
        """
        try:
            model_instance = self.get_model_instance(ModelType.TEXT_EMBEDDING)

            # Use `voyage-3` model for validate,
            # no matter what model you pass in, text completion model or chat model
            model_instance.validate_credentials(model="voyage-3", credentials=credentials)
        except CredentialsValidateFailedError as ex:
            raise ex
        except Exception as ex:
            logger.exception(f"{self.get_provider_schema().provider} credentials validate failed")
            raise ex
