"""
This file is used to hold the integration config for plugin testing.
"""

import shutil
import subprocess

from packaging.version import Version
from pydantic import Field, field_validator
from pydantic_settings import BaseSettings, SettingsConfigDict

_PLUGIN_NAMES = [
    "aiexec",
    "aiexec.exe",
    "aiexec-plugin",
    "aiexec-plugin.exe",
    "aiexec-plugin-darwin-amd64",
    "aiexec-plugin-darwin-arm64",
    "aiexec-plugin-linux-amd64",
    "aiexec-plugin-linux-arm64",
    "aiexec-plugin-windows-amd64.exe",
    "aiexec-plugin-windows-arm64.exe",
]


class IntegrationConfig(BaseSettings):
    aiexec_cli_path: str = Field(default="", description="The path to the aiexec cli")

    @field_validator("aiexec_cli_path")
    @classmethod
    def validate_aiexec_cli_path(cls, v):
        # find the aiexec cli path
        if not v:
            for plugin_name in _PLUGIN_NAMES:
                v = shutil.which(plugin_name)
                if v:
                    break

            if not v:
                raise ValueError("aiexec cli not found")

        # check aiexec version
        version = subprocess.check_output([v, "version"]).decode("utf-8")  # noqa: S603

        try:
            version = Version(version)
        except Exception as e:
            raise ValueError("aiexec cli version is not valid") from e

        if version < Version("0.1.0"):
            raise ValueError("aiexec cli version must be greater than 0.1.0 to support plugin run")

        return v

    model_config = SettingsConfigDict(env_file=".env")
