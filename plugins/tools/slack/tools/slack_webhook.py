from typing import Any, Generator
import httpx
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool


class SlackWebhookTool(Tool):
    def _invoke(
        self, tool_parameters: dict[str, Any]
    ) -> Generator[ToolInvokeMessage, None, None]:
        """
        Incoming Webhooks
        API Document: https://api.slack.com/messaging/webhooks
        """
        content = tool_parameters.get("content", "")
        if not content:
            yield self.create_text_message("Invalid parameter content")
        webhook_url = tool_parameters.get("webhook_url", "")
        if not webhook_url.startswith("https://hooks.slack.com/"):
            yield self.create_text_message(
                f"Invalid parameter webhook_url ${webhook_url}, not a valid Slack webhook URL"
            )
        headers = {"Content-Type": "application/json"}
        params = {}
        payload = {"text": content}
        try:
            res = httpx.post(webhook_url, headers=headers, params=params, json=payload)
            if res.is_success:
                yield self.create_text_message("Text message was sent successfully")
            else:
                yield self.create_text_message(
                    f"Failed to send the text message, status code: {res.status_code}, response: {res.text}"
                )
        except Exception as e:
            yield self.create_text_message("Failed to send message through webhook. {}".format(e))
