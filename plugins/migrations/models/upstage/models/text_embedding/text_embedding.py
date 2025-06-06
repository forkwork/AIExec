import base64
import time
from collections.abc import Mapping
from typing import Union
import numpy as np
from aiexec_plugin.entities.model import EmbeddingInputType, PriceType
from aiexec_plugin.entities.model.text_embedding import EmbeddingUsage, TextEmbeddingResult
from aiexec_plugin.errors.model import CredentialsValidateFailedError
from aiexec_plugin.interfaces.model.text_embedding import TextEmbeddingModel
from openai import OpenAI
from tokenizers import Tokenizer
from core.model_runtime.model_providers.upstage._common import _CommonUpstage


class UpstageTextEmbeddingModel(_CommonUpstage, TextEmbeddingModel):
    """
    Model class for Upstage text embedding model.
    """

    def _get_tokenizer(self) -> Tokenizer:
        return Tokenizer.from_pretrained("upstage/solar-1-mini-tokenizer")

    def _invoke(
        self,
        model: str,
        credentials: dict,
        texts: list[str],
        user: str | None = None,
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
        credentials_kwargs = self._to_credential_kwargs(credentials)
        client = OpenAI(**credentials_kwargs)
        extra_model_kwargs = {}
        if user:
            extra_model_kwargs["user"] = user
        extra_model_kwargs["encoding_format"] = "base64"
        context_size = self._get_context_size(model, credentials)
        max_chunks = self._get_max_chunks(model, credentials)
        embeddings: list[list[float]] = [[] for _ in range(len(texts))]
        tokens = []
        indices = []
        used_tokens = 0
        tokenizer = self._get_tokenizer()
        for i, text in enumerate(texts):
            token = tokenizer.encode(text, add_special_tokens=False).tokens
            for j in range(0, len(token), context_size):
                tokens += [token[j : j + context_size]]
                indices += [i]
        batched_embeddings = []
        _iter = range(0, len(tokens), max_chunks)
        for i in _iter:
            (embeddings_batch, embedding_used_tokens) = self._embedding_invoke(
                model=model, client=client, texts=tokens[i : i + max_chunks], extra_model_kwargs=extra_model_kwargs
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
                    model=model, client=client, texts=[texts[i]], extra_model_kwargs=extra_model_kwargs
                )
                used_tokens += embedding_used_tokens
                average = embeddings_batch[0]
            else:
                average = np.average(_result, axis=0, weights=num_tokens_in_batch[i])
            embeddings[i] = (average / np.linalg.norm(average)).tolist()
        usage = self._calc_response_usage(model=model, credentials=credentials, tokens=used_tokens)
        return TextEmbeddingResult(embeddings=embeddings, usage=usage, model=model)

    def get_num_tokens(self, model: str, credentials: dict, texts: list[str]) -> int:
        tokenizer = self._get_tokenizer()
        "\n        Get number of tokens for given prompt messages\n\n        :param model: model name\n        :param credentials: model credentials\n        :param texts: texts to embed\n        :return:\n        "
        if len(texts) == 0:
            return 0
        tokenizer = self._get_tokenizer()
        total_num_tokens = 0
        for text in texts:
            tokenized_text = tokenizer.encode(text)
            total_num_tokens += len(tokenized_text)
        return total_num_tokens

    def validate_credentials(self, model: str, credentials: Mapping) -> None:
        """
        Validate model credentials

        :param model: model name
        :param credentials: model credentials
        :return:
        """
        try:
            credentials_kwargs = self._to_credential_kwargs(credentials)
            client = OpenAI(**credentials_kwargs)
            self._embedding_invoke(model=model, client=client, texts=["ping"], extra_model_kwargs={})
        except Exception as ex:
            raise CredentialsValidateFailedError(str(ex))

    def _embedding_invoke(
        self, model: str, client: OpenAI, texts: Union[list[str], str], extra_model_kwargs: dict
    ) -> tuple[list[list[float]], int]:
        """
        Invoke embedding model
        :param model: model name
        :param client: model client
        :param texts: texts to embed
        :param extra_model_kwargs: extra model kwargs
        :return: embeddings and used tokens
        """
        response = client.embeddings.create(model=model, input=texts, **extra_model_kwargs)
        if "encoding_format" in extra_model_kwargs and extra_model_kwargs["encoding_format"] == "base64":
            return (
                [
                    list(np.frombuffer(base64.b64decode(embedding.embedding), dtype=np.float32))
                    for embedding in response.data
                ],
                response.usage.total_tokens,
            )
        return ([data.embedding for data in response.data], response.usage.total_tokens)

    def _calc_response_usage(self, model: str, credentials: dict, tokens: int) -> EmbeddingUsage:
        """
        Calculate response usage

        :param model: model name
        :param credentials: model credentials
        :param tokens: input tokens
        :return: usage
        """
        input_price_info = self.get_price(
            model=model, credentials=credentials, tokens=tokens, price_type=PriceType.INPUT
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
