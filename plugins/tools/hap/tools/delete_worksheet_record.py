from typing import Any, Generator
import httpx
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool


class DeleteWorksheetRecordTool(Tool):
    def _invoke(
        self, tool_parameters: dict[str, Any]
    ) -> Generator[ToolInvokeMessage, None, None]:
        appkey = tool_parameters.get("appkey", "")
        if not appkey:
            yield self.create_text_message("Invalid parameter App Key")
        sign = tool_parameters.get("sign", "")
        if not sign:
            yield self.create_text_message("Invalid parameter Sign")
        worksheet_id = tool_parameters.get("worksheet_id", "")
        if not worksheet_id:
            yield self.create_text_message("Invalid parameter Worksheet ID")
        row_id = tool_parameters.get("row_id", "")
        if not row_id:
            yield self.create_text_message("Invalid parameter Record Row ID")
        host = tool_parameters.get("host", "")
        if not host:
            host = "https://api.mingdao.com"
        elif not host.startswith(("http://", "https://")):
            yield self.create_text_message("Invalid parameter Host Address")
        else:
            host = f"{host.removesuffix('/')}/api"
        url = f"{host}/v2/open/worksheet/deleteRow"
        headers = {"Content-Type": "application/json"}
        payload = {"appKey": appkey, "sign": sign, "worksheetId": worksheet_id, "rowId": row_id}
        try:
            res = httpx.post(url, headers=headers, json=payload, timeout=30)
            res.raise_for_status()
            res_json = res.json()
            if res_json.get("error_code") != 1:
                yield self.create_text_message(f"Failed to delete the record. {res_json['error_msg']}")
            yield self.create_text_message("Successfully deleted the record.")
        except httpx.RequestError as e:
            yield self.create_text_message(f"Failed to delete the record, request error: {e}")
        except Exception as e:
            yield self.create_text_message(f"Failed to delete the record, unexpected error: {e}")
