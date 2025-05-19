from collections.abc import Generator
from typing import Optional, Union
from aiexec_plugin.entities.model.llm import LLMMode, LLMResult
from aiexec_plugin.entities.model.message import PromptMessage, PromptMessageTool
from yarl import URL
from aiexec_plugin import OAICompatLargeLanguageModel


class XAILargeLanguageModel(OAICompatLargeLanguageModel):
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
        self._add_custom_parameters(credentials)
        return super()._invoke(model, credentials, prompt_messages, model_parameters, tools, stop, stream)

    def validate_credentials(self, model: str, credentials: dict) -> None:
        self._add_custom_parameters(credentials)
        super().validate_credentials(model, credentials)

    @staticmethod
    def _add_custom_parameters(credentials) -> None:
        credentials["endpoint_url"] = str(URL(credentials["endpoint_url"])) or "https://api.x.ai/v1"
        credentials["mode"] = LLMMode.CHAT.value
        credentials["function_calling_type"] = "tool_call"
        credentials["stream_function_calling"] = "support"
        credentials["vision_support"] = "support"
