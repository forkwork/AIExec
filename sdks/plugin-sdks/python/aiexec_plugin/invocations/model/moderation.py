from aiexec_plugin.core.entities.invocation import InvokeType
from aiexec_plugin.core.runtime import BackwardsInvocation
from aiexec_plugin.entities.model.moderation import ModerationModelConfig, ModerationResult


class ModerationInvocation(BackwardsInvocation[ModerationResult]):
    def invoke(self, model_config: ModerationModelConfig, text: str) -> bool:
        """
        Invoke moderation
        """
        for data in self._backwards_invoke(
            InvokeType.Moderation,
            ModerationResult,
            {
                **model_config.model_dump(),
                "text": text,
            },
        ):
            return data.result

        raise Exception("No response from moderation")
