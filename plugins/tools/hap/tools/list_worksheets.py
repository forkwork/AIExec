import json
from typing import Any, Generator
import httpx
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool


class ListWorksheetsTool(Tool):
    def _invoke(
        self, tool_parameters: dict[str, Any]
    ) -> Generator[ToolInvokeMessage, None, None]:
        appkey = tool_parameters.get("appkey", "")
        if not appkey:
            yield self.create_text_message("Invalid parameter App Key")
        sign = tool_parameters.get("sign", "")
        if not sign:
            yield self.create_text_message("Invalid parameter Sign")
        host = tool_parameters.get("host", "")
        if not host:
            host = "https://api.mingdao.com"
        elif not (host.startswith("http://") or host.startswith("https://")):
            yield self.create_text_message("Invalid parameter Host Address")
        else:
            host = f"{host.removesuffix('/')}/api"
        url = f"{host}/v1/open/app/get"
        result_type = tool_parameters.get("result_type", "")
        if not result_type:
            result_type = "table"
        headers = {"Content-Type": "application/json"}
        params = {"appKey": appkey, "sign": sign}
        try:
            res = httpx.get(url, headers=headers, params=params, timeout=30)
            res_json = res.json()
            if res.is_success:
                if res_json["error_code"] != 1:
                    yield self.create_text_message(
                        "Failed to access the application. {}".format(res_json["error_msg"])
                    )
                elif result_type == "json":
                    worksheets = []
                    for section in res_json["data"]["sections"]:
                        worksheets.extend(self._extract_worksheets(section, result_type))
                    yield self.create_text_message(text=json.dumps(worksheets, ensure_ascii=False))
                else:
                    worksheets = "|worksheetId|worksheetName|description|\n|---|---|---|"
                    for section in res_json["data"]["sections"]:
                        worksheets += self._extract_worksheets(section, result_type)
                    yield self.create_text_message(worksheets)
            else:
                yield self.create_text_message(
                    f"Failed to list worksheets, status code: {res.status_code}, response: {res.text}"
                )
        except Exception as e:
            yield self.create_text_message("Failed to list worksheets, something went wrong: {}".format(e))

    def _extract_worksheets(self, section, type):
        items = []
        tables = ""
        for item in section.get("items", []):
            if item.get("type") == 0 and ("notes" not in item or item.get("notes") != "NO"):
                if type == "json":
                    filtered_item = {"id": item["id"], "name": item["name"], "notes": item.get("notes", "")}
                    items.append(filtered_item)
                else:
                    tables += f"\n|{item['id']}|{item['name']}|{item.get('notes', '')}|"
        for child_section in section.get("childSections", []):
            if type == "json":
                items.extend(self._extract_worksheets(child_section, "json"))
            else:
                tables += self._extract_worksheets(child_section, "table")
        return items if type == "json" else tables
