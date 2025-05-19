from typing import Mapping

from aiexec_plugin.entities.model import (
    AIModelEntity,
    I18nObject
)

from aiexec_plugin.interfaces.model.openai_compatible.tts import OAICompatText2SpeechModel


class OpenAIText2SpeechModel(OAICompatText2SpeechModel):

    def get_customizable_model_schema(self, model: str, credentials: Mapping) -> AIModelEntity:
        entity = super().get_customizable_model_schema(model, credentials)

        if "display_name" in credentials and credentials["display_name"] != "":
            entity.label= I18nObject(
                en_US=credentials["display_name"],
                zh_Hans=credentials["display_name"]
            )

        return entity
