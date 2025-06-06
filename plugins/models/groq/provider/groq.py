import logging
from aiexec_plugin import ModelProvider
from aiexec_plugin.entities.model import ModelType
from aiexec_plugin.errors.model import CredentialsValidateFailedError

logger = logging.getLogger(__name__)


class GroqProvider(ModelProvider):
    def validate_provider_credentials(self, credentials: dict) -> None:
        """
        Validate provider credentials
        if validate failed, raise exception

        :param credentials: provider credentials, credentials form defined in `provider_credential_schema`.
        """
        try:
            model_instance = self.get_model_instance(ModelType.LLM)
            model_instance.validate_credentials(model="llama3-8b-8192", credentials=credentials)
        except CredentialsValidateFailedError as ex:
            raise ex
        except Exception as ex:
            logger.exception(f"{self.get_provider_schema().provider} credentials validate failed")
            raise ex
