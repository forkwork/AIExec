from typing import Any, Generator

from aiexec_plugin import Tool
from aiexec_plugin.entities.tool import ToolInvokeMessage
from lark_api_utils import LarkRequest


class AddEventAttendeesTool(Tool):
    def _invoke(
        self, tool_parameters: dict[str, Any]
    ) -> Generator[ToolInvokeMessage, None, None]:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = LarkRequest(app_id, app_secret)
        event_id = tool_parameters.get("event_id")
        attendee_phone_or_email = tool_parameters.get("attendee_phone_or_email")
        need_notification = tool_parameters.get("need_notification", True)
        res = client.add_event_attendees(
            event_id, attendee_phone_or_email, need_notification
        )
        yield self.create_json_message(res)
