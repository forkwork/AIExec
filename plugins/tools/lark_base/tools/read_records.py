from typing import Any, Generator

from aiexec_plugin import Tool
from aiexec_plugin.entities.tool import ToolInvokeMessage
from lark_api_utils import LarkRequest


class ReadRecordsTool(Tool):
    def _invoke(
        self, tool_parameters: dict[str, Any]
    ) -> Generator[ToolInvokeMessage, None, None]:
        app_id = self.runtime.credentials.get("app_id")
        app_secret = self.runtime.credentials.get("app_secret")
        client = LarkRequest(app_id, app_secret)
        app_token = tool_parameters.get("app_token")
        table_id = tool_parameters.get("table_id")
        table_name = tool_parameters.get("table_name")
        record_ids = tool_parameters.get("record_ids")
        user_id_type = tool_parameters.get("user_id_type", "open_id")
        res = client.read_records(
            app_token, table_id, table_name, record_ids, user_id_type
        )
        yield self.create_json_message(res)
