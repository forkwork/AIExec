import logging
from collections.abc import Mapping

from aiexec_plugin import ModelProvider
from aiexec_plugin.entities.model import ModelType
from aiexec_plugin.errors.model import CredentialsValidateFailedError

logger = logging.getLogger(__name__)


class AmazonBedrockModelProvider(ModelProvider):
    def validate_provider_credentials(self, credentials: Mapping) -> None:
        """
        Validate provider credentials
        if validate failed, raise exception

        :param credentials: provider credentials, credentials form defined in `provider_credential_schema`.
        """
        try:
            model_instance = self.get_model_instance(ModelType.LLM)

            # Use `amazon.titan-text-lite-v1` model by default for validating credentials
            model_for_validation = credentials.get("model_for_validation", "amazon.titan-text-lite-v1")
            model_instance.validate_credentials(model=model_for_validation, credentials=credentials)
        except CredentialsValidateFailedError as ex:
            raise ex
        except Exception as ex:
            logger.exception(f"{self.get_provider_schema().provider} credentials validate failed")
            raise ex
