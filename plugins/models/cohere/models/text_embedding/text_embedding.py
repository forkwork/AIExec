import time
from typing import Optional

import cohere
import numpy as np
from cohere.core import RequestOptions
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
from aiexec_plugin.interfaces.model.text_embedding_model import TextEmbeddingModel


class CohereTextEmbeddingModel(TextEmbeddingModel):
    """
    Model class for Cohere text embedding model.
    """

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
        context_size = self._get_context_size(model, credentials)
        max_chunks = self._get_max_chunks(model, credentials)
        embeddings: list[list[float]] = [[] for _ in range(len(texts))]
        tokens = []
        indices = []
        used_tokens = 0
        for i, text in enumerate(texts):
            tokenize_response = self._tokenize(
                model=model, credentials=credentials, text=text
            )
            for j in range(0, len(tokenize_response), context_size):
                tokens += [tokenize_response[j : j + context_size]]
                indices += [i]
        batched_embeddings = []
        _iter = range(0, len(tokens), max_chunks)
        for i in _iter:
            (embeddings_batch, embedding_used_tokens) = self._embedding_invoke(
                model=model,
                credentials=credentials,
                texts=["".join(token) for token in tokens[i : i + max_chunks]],
            )
            used_tokens += embedding_used_tokens
            batched_embeddings += embeddings_batch
        results: list[list[list[float]]] = [[] for _ in range(len(texts))]
        num_tokens_in_batch: list[list[int]] = [[] for _ in range(len(texts))]
        for i in range(len(indices)):
            results[indices[i]].append(batched_embeddings[i])
            num_tokens_in_batch[indices[i]].append(len(tokens[i]))
        for i in range(len(texts)):
            _result = results[i]
            if len(_result) == 0:
                (embeddings_batch, embedding_used_tokens) = self._embedding_invoke(
                    model=model, credentials=credentials, texts=[" "]
                )
                used_tokens += embedding_used_tokens
                average = embeddings_batch[0]
            else:
                average = np.average(_result, axis=0, weights=num_tokens_in_batch[i])
            embedding = (average / np.linalg.norm(average)).tolist()
            if np.isnan(embedding).any():
                raise ValueError("Normalized embedding is nan please try again")
            embeddings[i] = embedding
        usage = self._calc_response_usage(
            model=model, credentials=credentials, tokens=used_tokens
        )
        return TextEmbeddingResult(embeddings=embeddings, usage=usage, model=model)

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
        if len(texts) == 0:
            return [0]
        tokens = []
        for text in texts:
            try:
                response = self._tokenize(
                    model=model, credentials=credentials, text=text
                )
            except Exception as e:
                raise self._transform_invoke_error(e)
            tokens.append(len(response))
        return tokens

    def _tokenize(self, model: str, credentials: dict, text: str) -> list[str]:
        """
        Tokenize text
        :param model: model name
        :param credentials: model credentials
        :param text: text to tokenize
        :return:
        """
        if not text:
            return []
        client = cohere.Client(
            credentials.get("api_key"), base_url=credentials.get("base_url")
        )
        response = client.tokenize(
            text=text,
            model=model,
            offline=False,
            request_options=RequestOptions(max_retries=0),
        )
        return response.token_strings

    def validate_credentials(self, model: str, credentials: dict) -> None:
        """
        Validate model credentials

        :param model: model name
        :param credentials: model credentials
        :return:
        """
        try:
            self._embedding_invoke(model=model, credentials=credentials, texts=["ping"])
        except Exception as ex:
            raise CredentialsValidateFailedError(str(ex))

    def _embedding_invoke(
        self, model: str, credentials: dict, texts: list[str]
    ) -> tuple[list[list[float]], int]:
        """
        Invoke embedding model

        :param model: model name
        :param credentials: model credentials
        :param texts: texts to embed
        :return: embeddings and used tokens
        """
        client = cohere.Client(
            credentials.get("api_key"), base_url=credentials.get("base_url")
        )
        response = client.embed(
            texts=texts,
            model=model,
            input_type="search_document" if len(texts) > 1 else "search_query",
            request_options=RequestOptions(max_retries=1),
        )
        return (response.embeddings, int(response.meta.billed_units.input_tokens))

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

    @property
    def _invoke_error_mapping(self) -> dict[type[InvokeError], list[type[Exception]]]:
        """
        Map model invoke error to unified error
        The key is the error type thrown to the caller
        The value is the error type thrown by the model,
        which needs to be converted into a unified error type for the caller.

        :return: Invoke error mapping
        """
        return {
            InvokeConnectionError: [
                cohere.errors.service_unavailable_error.ServiceUnavailableError
            ],
            InvokeServerUnavailableError: [
                cohere.errors.internal_server_error.InternalServerError
            ],
            InvokeRateLimitError: [
                cohere.errors.too_many_requests_error.TooManyRequestsError
            ],
            InvokeAuthorizationError: [
                cohere.errors.unauthorized_error.UnauthorizedError,
                cohere.errors.forbidden_error.ForbiddenError,
            ],
            InvokeBadRequestError: [
                cohere.core.api_error.ApiError,
                cohere.errors.bad_request_error.BadRequestError,
                cohere.errors.not_found_error.NotFoundError,
            ],
        }
