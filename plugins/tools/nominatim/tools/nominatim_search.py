import json
from typing import Any, Generator

import requests
from aiexec_plugin import Tool
from aiexec_plugin.entities.tool import ToolInvokeMessage


class NominatimSearchTool(Tool):
    def _invoke(
            self, tool_parameters: dict[str, Any]
    ) -> Generator[ToolInvokeMessage, None, None]:
        query = tool_parameters.get("query", "")
        limit = tool_parameters.get("limit", 10)
        if not query:
            yield self.create_text_message("Please input a search query")
        params = {"q": query, "format": "json", "limit": limit, "addressdetails": 1}
        yield self._make_request("search", params)

    def _make_request(self, endpoint: str, params: dict) -> ToolInvokeMessage:
        base_url = self.runtime.credentials.get("base_url", "https://nominatim.openstreetmap.org")
        try:
            headers = {"User-Agent": "AiexecNominatimTool/1.0"}
            s = requests.session()
            response = s.request(method="GET", headers=headers, url=f"{base_url}/{endpoint}", params=params)
            response_data = response.json()
            if response.status_code == 200:
                s.close()
                return self.create_text_message(
                    self.session.model.summary.invoke(text=json.dumps(response_data, ensure_ascii=False),
                                                      instruction="")
                )
            else:
                return self.create_text_message(f"Error: {response.status_code} - {response.text}")
        except Exception as e:
            return self.create_text_message(f"An error occurred: {str(e)}")
