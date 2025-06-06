import json
import logging
from typing import Any, Union
import boto3
from botocore.exceptions import BotoCoreError
from pydantic import BaseModel, Field
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class GuardrailParameters(BaseModel):
    guardrail_id: str = Field(..., description="The identifier of the guardrail")
    guardrail_version: str = Field(..., description="The version of the guardrail")
    source: str = Field(..., description="The source of the content")
    text: str = Field(..., description="The text to apply the guardrail to")
    aws_region: str = Field(..., description="AWS region for the Bedrock client")


class ApplyGuardrailTool(Tool):
    def _invoke(
        self, user_id: str, tool_parameters: dict[str, Any]
    ) -> Union[ToolInvokeMessage, list[ToolInvokeMessage]]:
        """
        Invoke the ApplyGuardrail tool
        """
        try:
            params = GuardrailParameters(**tool_parameters)
            bedrock_client = boto3.client("bedrock-runtime", region_name=params.aws_region)
            response = bedrock_client.apply_guardrail(
                guardrailIdentifier=params.guardrail_id,
                guardrailVersion=params.guardrail_version,
                source=params.source,
                content=[{"text": {"text": params.text}}],
            )
            logger.info(f"Raw response from AWS: {json.dumps(response, indent=2)}")
            if not response:
                return self.create_text_message(text="Received empty response from AWS Bedrock.")
            action = response.get("action", "No action specified")
            outputs = response.get("outputs", [])
            output = outputs[0].get("text", "No output received") if outputs else "No output received"
            assessments = response.get("assessments", [])
            formatted_assessments = []
            for assessment in assessments:
                for policy_type, policy_data in assessment.items():
                    if isinstance(policy_data, dict) and "topics" in policy_data:
                        for topic in policy_data["topics"]:
                            formatted_assessments.append(
                                f"Policy: {policy_type}, Topic: {topic['name']}, Type: {topic['type']}, Action: {topic['action']}"
                            )
                    else:
                        formatted_assessments.append(f"Policy: {policy_type}, Data: {policy_data}")
            result = f"Action: {action}\n "
            result += f"Output: {output}\n "
            if formatted_assessments:
                result += "Assessments:\n " + "\n ".join(formatted_assessments) + "\n "
            return self.create_text_message(text=result)
        except BotoCoreError as e:
            error_message = f"AWS service error: {str(e)}"
            logger.error(error_message, exc_info=True)
            return self.create_text_message(text=error_message)
        except json.JSONDecodeError as e:
            error_message = f"JSON parsing error: {str(e)}"
            logger.error(error_message, exc_info=True)
            return self.create_text_message(text=error_message)
        except Exception as e:
            error_message = f"An unexpected error occurred: {str(e)}"
            logger.error(error_message, exc_info=True)
            return self.create_text_message(text=error_message)
