from binascii import hexlify
from typing import IO

from aiexec_plugin.core.entities.invocation import InvokeType
from aiexec_plugin.core.runtime import BackwardsInvocation
from aiexec_plugin.entities.model.speech2text import Speech2TextModelConfig, Speech2TextResult


class Speech2TextInvocation(BackwardsInvocation[Speech2TextResult]):
    def invoke(self, model_config: Speech2TextModelConfig, file: IO[bytes]) -> str:
        """
        Invoke speech2text
        """
        for data in self._backwards_invoke(
            InvokeType.Speech2Text,
            Speech2TextResult,
            {
                **model_config.model_dump(),
                "file": hexlify(file.read()),
            },
        ):
            return data.result

        raise Exception("No response from speech2text")
