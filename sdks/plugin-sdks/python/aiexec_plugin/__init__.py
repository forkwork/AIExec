from gevent import monkey

# patch all the blocking calls
monkey.patch_all(sys=True)

from aiexec_plugin.config.config import AiexecPluginEnv
from aiexec_plugin.interfaces.agent import AgentProvider, AgentStrategy
from aiexec_plugin.interfaces.endpoint import Endpoint
from aiexec_plugin.interfaces.model import ModelProvider
from aiexec_plugin.interfaces.model.large_language_model import LargeLanguageModel
from aiexec_plugin.interfaces.model.moderation_model import ModerationModel
from aiexec_plugin.interfaces.model.openai_compatible.llm import OAICompatLargeLanguageModel
from aiexec_plugin.interfaces.model.openai_compatible.provider import OAICompatProvider
from aiexec_plugin.interfaces.model.openai_compatible.rerank import OAICompatRerankModel
from aiexec_plugin.interfaces.model.openai_compatible.speech2text import OAICompatSpeech2TextModel
from aiexec_plugin.interfaces.model.openai_compatible.text_embedding import OAICompatEmbeddingModel
from aiexec_plugin.interfaces.model.openai_compatible.tts import OAICompatText2SpeechModel
from aiexec_plugin.interfaces.model.rerank_model import RerankModel
from aiexec_plugin.interfaces.model.speech2text_model import Speech2TextModel
from aiexec_plugin.interfaces.model.text_embedding_model import TextEmbeddingModel
from aiexec_plugin.interfaces.model.tts_model import TTSModel
from aiexec_plugin.interfaces.tool import Tool, ToolProvider
from aiexec_plugin.invocations.file import File
from aiexec_plugin.plugin import Plugin

__all__ = [
    "AgentProvider",
    "AgentStrategy",
    "AiexecPluginEnv",
    "Endpoint",
    "File",
    "LargeLanguageModel",
    "ModelProvider",
    "ModerationModel",
    "OAICompatEmbeddingModel",
    "OAICompatLargeLanguageModel",
    "OAICompatProvider",
    "OAICompatRerankModel",
    "OAICompatSpeech2TextModel",
    "OAICompatText2SpeechModel",
    "Plugin",
    "RerankModel",
    "Speech2TextModel",
    "TTSModel",
    "TextEmbeddingModel",
    "Tool",
    "ToolProvider",
]
