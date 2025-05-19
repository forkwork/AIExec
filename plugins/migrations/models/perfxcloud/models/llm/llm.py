from collections.abc import Generator
from typing import Optional, Union
from urllib.parse import urlparse
import tiktoken
from aiexec_plugin.entities.model.llm import LLMResult
from aiexec_plugin.entities.model.message import PromptMessage, PromptMessageTool
from core.model_runtime.model_providers.openai.llm.llm import OpenAILargeLanguageModel


class PerfXCloudLargeLanguageModel(OpenAILargeLanguageModel):
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

    def _num_tokens_from_string(self, model: str, text: str, tools: Optional[list[PromptMessageTool]] = None) -> int:
        """
        Calculate num tokens for text completion model with tiktoken package.

        :param model: model name
        :param text: prompt text
        :param tools: tools for tool calling
        :return: number of tokens
        """
        encoding = tiktoken.get_encoding("cl100k_base")
        num_tokens = len(encoding.encode(text))
        if tools:
            num_tokens += self._num_tokens_for_tools(encoding, tools)
        return num_tokens

    def _num_tokens_from_messages(
        self, model: str, messages: list[PromptMessage], tools: Optional[list[PromptMessageTool]] = None
    ) -> int:
        """Calculate num tokens for gpt-3.5-turbo and gpt-4 with tiktoken package.

        Official documentation: https://github.com/openai/openai-cookbook/blob/
        main/examples/How_to_format_inputs_to_ChatGPT_models.ipynb"""
        encoding = tiktoken.get_encoding("cl100k_base")
        tokens_per_message = 3
        tokens_per_name = 1
        num_tokens = 0
        messages_dict = [self._convert_prompt_message_to_dict(m) for m in messages]
        for message in messages_dict:
            num_tokens += tokens_per_message
            for key, value in message.items():
                if isinstance(value, list):
                    text = ""
                    for item in value:
                        if isinstance(item, dict) and item["type"] == "text":
                            text += item["text"]
                    value = text
                if key == "tool_calls":
                    for tool_call in value:
                        for t_key, t_value in tool_call.items():
                            num_tokens += len(encoding.encode(t_key))
                            if t_key == "function":
                                for f_key, f_value in t_value.items():
                                    num_tokens += len(encoding.encode(f_key))
                                    num_tokens += len(encoding.encode(f_value))
                            else:
                                num_tokens += len(encoding.encode(t_key))
                                num_tokens += len(encoding.encode(t_value))
                else:
                    num_tokens += len(encoding.encode(str(value)))
                if key == "name":
                    num_tokens += tokens_per_name
        num_tokens += 3
        if tools:
            num_tokens += self._num_tokens_for_tools(encoding, tools)
        return num_tokens

    @staticmethod
    def _add_custom_parameters(credentials: dict) -> None:
        credentials["mode"] = "chat"
        credentials["openai_api_key"] = credentials["api_key"]
        if "endpoint_url" not in credentials or credentials["endpoint_url"] == "":
            credentials["openai_api_base"] = "https://cloud.perfxlab.cn"
        else:
            parsed_url = urlparse(credentials["endpoint_url"])
            credentials["openai_api_base"] = f"{parsed_url.scheme}://{parsed_url.netloc}"
