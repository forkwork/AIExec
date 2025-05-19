import binascii
from collections.abc import Generator

from aiexec_plugin.entities.model.tts import TTSModelConfig
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin.interfaces.tool import Tool


class TTS(Tool):
    def _invoke(self, tool_parameters: dict) -> Generator[ToolInvokeMessage, None, None]:
        response = self.session.model.tts.invoke(
            model_config=TTSModelConfig(
                provider="openai",
                model="tts-1",
                voice="alloy",
            ),
            content_text="Hello, world!",
        )

        for chunk in response:
            yield self.create_text_message(binascii.hexlify(chunk).decode())
