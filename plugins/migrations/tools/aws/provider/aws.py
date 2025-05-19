from aiexec_plugin.errors.tool import ToolProviderCredentialValidationError
from tools.sagemaker_text_rerank import SageMakerReRankTool
from aiexec_plugin import ToolProvider


class SageMakerProvider(ToolProvider):
    def _validate_credentials(self, credentials: dict) -> None:
        try:
            SageMakerReRankTool.from_credentials(credentials, user_id="").invoke(
                tool_parameters={
                    "sagemaker_endpoint": "",
                    "query": "misaka mikoto",
                    "candidate_texts": "hello$$$hello world",
                    "topk": 5,
                    "aws_region": "",
                }
            )
        except Exception as e:
            raise ToolProviderCredentialValidationError(str(e))
