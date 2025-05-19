from collections.abc import Generator

from aiexec_plugin.entities.model.moderation import ModerationModelConfig
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin.interfaces.tool import Tool


class Moderation(Tool):
    def _invoke(self, tool_parameters: dict) -> Generator[ToolInvokeMessage, None, None]:
        response = self.session.model.moderation.invoke(
            model_config=ModerationModelConfig(
                provider="openai",
                model="text-moderation-stable",
            ),
            text=tool_parameters["text"],
        )

        yield self.create_json_message(
            {
                "data": response,
            }
        )
