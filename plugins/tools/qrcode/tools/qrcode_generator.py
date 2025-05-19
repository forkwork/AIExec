import io
import logging
from typing import Any, Generator

from aiexec_plugin import Tool
from aiexec_plugin.entities.tool import ToolInvokeMessage
from qrcode.constants import ERROR_CORRECT_H, ERROR_CORRECT_L, ERROR_CORRECT_M, ERROR_CORRECT_Q
from qrcode.image.base import BaseImage
from qrcode.image.pure import PyPNGImage
from qrcode.main import QRCode


class QRCodeGeneratorTool(Tool):
    error_correction_levels: dict[str, int] = {
        "L": ERROR_CORRECT_L,
        "M": ERROR_CORRECT_M,
        "Q": ERROR_CORRECT_Q,
        "H": ERROR_CORRECT_H,
    }

    def _invoke(
        self, tool_parameters: dict[str, Any]
    ) -> Generator[ToolInvokeMessage, None, None]:
        """
        invoke tools
        """
        content = tool_parameters.get("content", "")
        if not content:
            yield self.create_text_message("Invalid parameter content")
        border = tool_parameters.get("border", 0)
        if border < 0 or border > 100:
            yield self.create_text_message("Invalid parameter border")
        error_correction = tool_parameters.get("error_correction", "")
        if error_correction not in self.error_correction_levels:
            yield self.create_text_message("Invalid parameter error_correction")
        try:
            image = self._generate_qrcode(content, border, error_correction)
            image_bytes = self._image_to_byte_array(image)
            yield self.create_blob_message(
                blob=image_bytes, meta={"mime_type": "image/png"}
            )
        except Exception:
            logging.exception(f"Failed to generate QR code for content: {content}")
            yield self.create_text_message("Failed to generate QR code")

    def _generate_qrcode(self, content: str, border: int, error_correction: str) -> BaseImage:
        qr = QRCode(
            image_factory=PyPNGImage, error_correction=self.error_correction_levels.get(error_correction), border=border
        )
        qr.add_data(data=content)
        qr.make(fit=True)
        img = qr.make_image()
        return img

    @staticmethod
    def _image_to_byte_array(image: BaseImage) -> bytes:
        byte_stream = io.BytesIO()
        image.save(byte_stream)
        return byte_stream.getvalue()
