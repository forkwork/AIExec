from collections.abc import Generator
from typing import Optional, Union
from aiexec_plugin import OAICompatLargeLanguageModel
from aiexec_plugin.entities.model.llm import LLMResult
from aiexec_plugin.entities.model.message import PromptMessage, PromptMessageTool


class LeptonAILargeLanguageModel(OAICompatLargeLanguageModel):
    MODEL_PREFIX_MAP = {
        "llama2-7b": "llama2-7b",
        "gemma-7b": "gemma-7b",
        "mistral-7b": "mistral-7b",
        "mixtral-8x7b": "mixtral-8x7b",
        "llama3-70b": "llama3-70b",
        "llama2-13b": "llama2-13b",
    }

    def _invoke(
        self,
        model: str,
        credentials: dict,
        prompt_messages: list[PromptMessage],
        model_parameters: dict,
        tools: Optional[list[PromptMessageTool]] = None,
        stop: Optional[list[str]] = None,
        stream: bool = True,
        user: Optional[str] = None,
    ) -> Union[LLMResult, Generator]:
        self._add_custom_parameters(credentials, model)
        return super()._invoke(
            model, credentials, prompt_messages, model_parameters, tools, stop, stream
        )

    def validate_credentials(self, model: str, credentials: dict) -> None:
        self._add_custom_parameters(credentials, model)
        super().validate_credentials(model, credentials)

    @classmethod
    def _add_custom_parameters(cls, credentials: dict, model: str) -> None:
        credentials["mode"] = "chat"
        credentials["endpoint_url"] = (
            f"https://{cls.MODEL_PREFIX_MAP[model]}.lepton.run/api/v1"
        )
