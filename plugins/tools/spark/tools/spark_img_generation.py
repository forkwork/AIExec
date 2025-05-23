import hmac
import json
from hashlib import sha256
from base64 import b64decode, b64encode
from datetime import datetime
from time import mktime
from typing import Any,Generator
from urllib.parse import urlencode
from wsgiref.handlers import format_date_time
import requests
from aiexec_plugin.entities.tool import ToolInvokeMessage
from aiexec_plugin import Tool


class AssembleHeaderError(Exception):
    def __init__(self, msg):
        self.message = msg


class Url:
    def __init__(self, host, path, schema):
        self.host = host
        self.path = path
        self.schema = schema


def sha256base64(data):
    sha = sha256()
    sha.update(data)
    digest = b64encode(sha.digest()).decode(encoding="utf-8")
    return digest


def parse_url(request_url):
    stidx = request_url.index("://")
    host = request_url[stidx + 3 :]
    schema = request_url[: stidx + 3]
    edidx = host.index("/")
    if edidx <= 0:
        raise AssembleHeaderError("invalid request url:" + request_url)
    path = host[edidx:]
    host = host[:edidx]
    u = Url(host, path, schema)
    return u


def assemble_ws_auth_url(request_url, method="GET", api_key="", api_secret=""):
    u = parse_url(request_url)
    host = u.host
    path = u.path
    now = datetime.now()
    date = format_date_time(mktime(now.timetuple()))
    signature_origin = "host: {}\ndate: {}\n{} {} HTTP/1.1".format(host, date, method, path)
    signature_sha = hmac.new(
        api_secret.encode("utf-8"), signature_origin.encode("utf-8"), digestmod=sha256
    ).digest()
    signature_sha = b64encode(signature_sha).decode(encoding="utf-8")
    authorization_origin = (
        f'api_key="{api_key}", algorithm="hmac-sha256", headers="host date request-line", signature="{signature_sha}"'
    )
    authorization = b64encode(authorization_origin.encode("utf-8")).decode(encoding="utf-8")
    values = {"host": host, "date": date, "authorization": authorization}
    return request_url + "?" + urlencode(values)


def get_body(appid, text):
    body = {
        "header": {"app_id": appid, "uid": "123456789"},
        "parameter": {"chat": {"domain": "general", "temperature": 0.5, "max_tokens": 4096}},
        "payload": {"message": {"text": [{"role": "user", "content": text}]}},
    }
    return body


def spark_response(text, appid, apikey, apisecret):
    host = "http://spark-api.cn-huabei-1.xf-yun.com/v2.1/tti"
    url = assemble_ws_auth_url(host, method="POST", api_key=apikey, api_secret=apisecret)
    content = get_body(appid, text)
    response = requests.post(url, json=content, headers={"content-type": "application/json"}).text
    return response


class SparkImgGeneratorTool(Tool):
    def _invoke(
        self,  tool_parameters: dict[str, Any]
    ) -> Generator[ToolInvokeMessage, None, None]:
        """
        invoke tools
        """
        if "APPID" not in self.runtime.credentials or not self.runtime.credentials.get("APPID"):
            yield self.create_text_message("APPID  is required.")
        if "APISecret" not in self.runtime.credentials or not self.runtime.credentials.get("APISecret"):
            yield self.create_text_message("APISecret  is required.")
        if "APIKey" not in self.runtime.credentials or not self.runtime.credentials.get("APIKey"):
            yield self.create_text_message("APIKey  is required.")
        prompt = tool_parameters.get("prompt", "")
        if not prompt:
            yield self.create_text_message("Please input prompt")
        res = self.img_generation(prompt)
        for image in res:
            if isinstance(image, dict) and "base64_image" in image:
                yield self.create_blob_message(
                    blob=b64decode(image["base64_image"]),
                    meta={"mime_type": "image/png"},
                )
            else:
                yield self.create_text_message("Invalid image data received.")

    def img_generation(self, prompt):
        response = spark_response(
            text=prompt,
            appid=self.runtime.credentials.get("APPID"),
            apikey=self.runtime.credentials.get("APIKey"),
            apisecret=self.runtime.credentials.get("APISecret"),
        )
        data = json.loads(response)
        code = data["header"]["code"]
        if code != 0:
            return self.create_text_message(f"error: {code}, {data}")
        else:
            text = data["payload"]["choices"]["text"]
            image_content = text[0]
            image_base = image_content["content"]
            json_data = {"base64_image": image_base}
        return [json_data]
