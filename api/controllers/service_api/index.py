from flask_restful import Resource

from configs import aiexec_config
from controllers.service_api import api


class IndexApi(Resource):
    def get(self):
        return {
            "welcome": "Aiexec OpenAPI",
            "api_version": "v1",
            "server_version": aiexec_config.CURRENT_VERSION,
        }


api.add_resource(IndexApi, "/")
