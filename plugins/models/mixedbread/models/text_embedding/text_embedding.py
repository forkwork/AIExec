import time
from json import JSONDecodeError, dumps
from typing import Optional

import requests
from aiexec_plugin.entities import I18nObject
from aiexec_plugin.entities.model import (
    AIModelEntity,
    EmbeddingInputType,
    FetchFrom,
    ModelPropertyKey,
    ModelType,
    PriceType,
)
from aiexec_plugin.entities.model.text_embedding import (
    EmbeddingUsage,
    TextEmbeddingResult,
)
from aiexec_plugin.errors.model import (
    CredentialsValidateFailedError,
    InvokeAuthorizationError,
    InvokeBadRequestError,
    InvokeConnectionError,
    InvokeError,
    InvokeRateLimitError,
    InvokeServerUnavailableError,
)
from aiexec_plugin.interfaces.model.text_embedding_model import TextEmbeddingModel


class MixedBreadTextEmbeddingModel(TextEmbeddingModel):
    """
    Model class for MixedBread text embedding model.
    """

    api_base: str = "https://api.mixedbread.com/v1"

    def _invoke(
        self,
        model: str,
        credentials: dict,
        texts: list[str],
        user: Optional[str] = None,
        input_type: EmbeddingInputType = EmbeddingInputType.DOCUMENT,
    ) -> TextEmbeddingResult:
        """
        Invoke text embedding model

        :param model: model name
        :param credentials: model credentials
        :param texts: texts to embed
        :param user: unique user id
        :param input_type: input type
        :return: embeddings result
        """
        api_key = credentials["api_key"]
        if not api_key:
            raise CredentialsValidateFailedError("api_key is required")

        base_url = credentials.get("base_url", self.api_base)
        base_url = base_url.removesuffix("/")

        url = base_url + "/embeddings"
        headers = {
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json",
        }

        data = {"model": f"mixedbread-ai/{model}", "input": texts}

        try:
            response = requests.post(url, headers=headers, data=dumps(data))
        except Exception as e:
            raise InvokeConnectionError(str(e))

        if response.status_code != 200:
            try:
                resp = response.json()
                msg = resp.get("detail", resp["message"])
                if response.status_code == 401:
                    raise InvokeAuthorizationError(msg)
                elif response.status_code == 429:
                    raise InvokeRateLimitError(msg)
                elif response.status_code == 500:
                    raise InvokeServerUnavailableError(msg)
                else:
                    raise InvokeBadRequestError(msg)
            except JSONDecodeError as e:
                raise InvokeServerUnavailableError(
                    f"Failed to convert response to json: {e} with text: {response.text}"
                )

        try:
            resp = response.json()
            embeddings = resp["data"]
            usage = resp["usage"]
        except Exception as e:
            raise InvokeServerUnavailableError(
                f"Failed to convert response to json: {e} with text: {response.text}"
            )

        usage = self._calc_response_usage(
            model=model, credentials=credentials, tokens=usage["total_tokens"]
        )

        result = TextEmbeddingResult(
            model=model,
            embeddings=[[float(data) for data in x["embedding"]] for x in embeddings],
            usage=usage,
        )

        return result

    def get_num_tokens(
        self, model: str, credentials: dict, texts: list[str]
    ) -> list[int]:
        """
        Get number of tokens for given prompt messages

        :param model: model name
        :param credentials: model credentials
        :param texts: texts to embed
        :return:
        """
        tokens = []
        for text in texts:
            tokens.append(self._get_num_tokens_by_gpt2(text))
        return tokens

    def validate_credentials(self, model: str, credentials: dict) -> None:
        """
        Validate model credentials

        :param model: model name
        :param credentials: model credentials
        :return:
        """
        try:
            self._invoke(model=model, credentials=credentials, texts=["ping"])
        except Exception as e:
            raise CredentialsValidateFailedError(f"Credentials validation failed: {e}")

    @property
    def _invoke_error_mapping(self) -> dict[type[InvokeError], list[type[Exception]]]:
        return {
            InvokeConnectionError: [InvokeConnectionError],
            InvokeServerUnavailableError: [InvokeServerUnavailableError],
            InvokeRateLimitError: [InvokeRateLimitError],
            InvokeAuthorizationError: [InvokeAuthorizationError],
            InvokeBadRequestError: [KeyError, InvokeBadRequestError],
        }

    def _calc_response_usage(
        self, model: str, credentials: dict, tokens: int
    ) -> EmbeddingUsage:
        """
        Calculate response usage

        :param model: model name
        :param credentials: model credentials
        :param tokens: input tokens
        :return: usage
        """
        # get input price info
        input_price_info = self.get_price(
            model=model,
            credentials=credentials,
            price_type=PriceType.INPUT,
            tokens=tokens,
        )

        # transform usage
        usage = EmbeddingUsage(
            tokens=tokens,
            total_tokens=tokens,
            unit_price=input_price_info.unit_price,
            price_unit=input_price_info.unit,
            total_price=input_price_info.total_amount,
            currency=input_price_info.currency,
            latency=time.perf_counter() - self.started_at,
        )

        return usage

    def get_customizable_model_schema(
        self, model: str, credentials: dict
    ) -> Optional[AIModelEntity]:
        """
        Get customizable model schema

        :param model: model name
        :param credentials: model credentials
        :return: model schema
        """
        return AIModelEntity(
            model=model,
            label=I18nObject(en_US=model),
            model_type=ModelType.TEXT_EMBEDDING,
            fetch_from=FetchFrom.CUSTOMIZABLE_MODEL,
            model_properties={
                ModelPropertyKey.CONTEXT_SIZE: int(
                    credentials.get("context_size", "512")
                )
            },
            parameter_rules=[],
        )
