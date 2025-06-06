import time
from json import dumps
from typing import Optional
from aiexec_plugin import TextEmbeddingModel
from aiexec_plugin.entities.model import EmbeddingInputType, PriceType
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
from requests import post
from models.llm.baichuan_tokenizer import BaichuanTokenizer
from models.llm.baichuan_turbo_errors import (
    BadRequestError,
    InsufficientAccountBalanceError,
    InternalServerError,
    InvalidAPIKeyError,
    InvalidAuthenticationError,
    RateLimitReachedError,
)


class BaichuanTextEmbeddingModel(TextEmbeddingModel):
    """
    Model class for BaiChuan text embedding model.
    """

    api_base: str = "http://api.baichuan-ai.com/v1/embeddings"

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
        if model != "baichuan-text-embedding":
            raise ValueError("Invalid model name")
        if not api_key:
            raise CredentialsValidateFailedError("api_key is required")
        chunks = []
        for i in range(0, len(texts), 16):
            chunks.append(texts[i : i + 16])
        embeddings = []
        token_usage = 0
        for chunk in chunks:
            (chunk_embeddings, chunk_usage) = self.embedding(
                model=model, api_key=api_key, texts=chunk, user=user
            )
            embeddings.extend(chunk_embeddings)
            token_usage += chunk_usage
        result = TextEmbeddingResult(
            model=model,
            embeddings=embeddings,
            usage=self._calc_response_usage(
                model=model, credentials=credentials, tokens=token_usage
            ),
        )
        return result

    def embedding(
        self, model: str, api_key, texts: list[str], user: Optional[str] = None
    ) -> tuple[list[list[float]], int]:
        """
        Embed given texts

        :param model: model name
        :param credentials: model credentials
        :param texts: texts to embed
        :param user: unique user id
        :return: embeddings result
        """
        url = self.api_base
        headers = {
            "Authorization": "Bearer " + api_key,
            "Content-Type": "application/json",
        }
        data = {"model": "Baichuan-Text-Embedding", "input": texts}
        try:
            response = post(url, headers=headers, data=dumps(data))
        except Exception as e:
            raise InvokeConnectionError(str(e))
        if response.status_code != 200:
            try:
                resp = response.json()
                err = resp["error"]["code"]
                msg = resp["error"]["message"]
            except Exception as e:
                raise InternalServerError(
                    f"Failed to convert response to json: {e} with text: {response.text}"
                )
            if err == "invalid_api_key":
                raise InvalidAPIKeyError(msg)
            elif err == "insufficient_quota":
                raise InsufficientAccountBalanceError(msg)
            elif err == "invalid_authentication":
                raise InvalidAuthenticationError(msg)
            elif err and "rate" in err:
                raise RateLimitReachedError(msg)
            elif err and "internal" in err:
                raise InternalServerError(msg)
            elif err == "api_key_empty":
                raise InvalidAPIKeyError(msg)
            else:
                raise InternalServerError(f"Unknown error: {err} with message: {msg}")
        try:
            resp = response.json()
            embeddings = resp["data"]
            usage = resp["usage"]
        except Exception as e:
            raise InternalServerError(
                f"Failed to convert response to json: {e} with text: {response.text}"
            )
        return ([data["embedding"] for data in embeddings], usage["total_tokens"])

    def get_num_tokens(
        self, model: str, credentials: dict, texts: list[str]
    ) -> list[int]:
        """
        Get number of tokens for given prompt messages

        :param model: model name
        :param credentials: model credentials
        :param texts: texts to embed
        :return: number of tokens for each text
        """
        num_tokens = []
        for text in texts:
            num_tokens.append(BaichuanTokenizer._get_num_tokens(text))
        return num_tokens

    def validate_credentials(self, model: str, credentials: dict) -> None:
        """
        Validate model credentials

        :param model: model name
        :param credentials: model credentials
        :return:
        """
        try:
            self._invoke(model=model, credentials=credentials, texts=["ping"])
        except InvalidAPIKeyError:
            raise CredentialsValidateFailedError("Invalid api key")

    @property
    def _invoke_error_mapping(self) -> dict[type[InvokeError], list[type[Exception]]]:
        return {
            InvokeConnectionError: [],
            InvokeServerUnavailableError: [InternalServerError],
            InvokeRateLimitError: [RateLimitReachedError],
            InvokeAuthorizationError: [
                InvalidAuthenticationError,
                InsufficientAccountBalanceError,
                InvalidAPIKeyError,
            ],
            InvokeBadRequestError: [BadRequestError, KeyError],
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
        input_price_info = self.get_price(
            model=model,
            credentials=credentials,
            price_type=PriceType.INPUT,
            tokens=tokens,
        )
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
