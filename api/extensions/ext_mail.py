import logging
from typing import Optional

from flask import Flask

from configs import aiexec_config
from aiexec_app import AiexecApp


class Mail:
    def __init__(self):
        self._client = None
        self._default_send_from = None

    def is_inited(self) -> bool:
        return self._client is not None

    def init_app(self, app: Flask):
        mail_type = aiexec_config.MAIL_TYPE
        if not mail_type:
            logging.warning("MAIL_TYPE is not set")
            return

        if aiexec_config.MAIL_DEFAULT_SEND_FROM:
            self._default_send_from = aiexec_config.MAIL_DEFAULT_SEND_FROM

        match mail_type:
            case "resend":
                import resend

                api_key = aiexec_config.RESEND_API_KEY
                if not api_key:
                    raise ValueError("RESEND_API_KEY is not set")

                api_url = aiexec_config.RESEND_API_URL
                if api_url:
                    resend.api_url = api_url

                resend.api_key = api_key
                self._client = resend.Emails
            case "smtp":
                from libs.smtp import SMTPClient

                if not aiexec_config.SMTP_SERVER or not aiexec_config.SMTP_PORT:
                    raise ValueError("SMTP_SERVER and SMTP_PORT are required for smtp mail type")
                if not aiexec_config.SMTP_USE_TLS and aiexec_config.SMTP_OPPORTUNISTIC_TLS:
                    raise ValueError("SMTP_OPPORTUNISTIC_TLS is not supported without enabling SMTP_USE_TLS")
                self._client = SMTPClient(
                    server=aiexec_config.SMTP_SERVER,
                    port=aiexec_config.SMTP_PORT,
                    username=aiexec_config.SMTP_USERNAME or "",
                    password=aiexec_config.SMTP_PASSWORD or "",
                    _from=aiexec_config.MAIL_DEFAULT_SEND_FROM or "",
                    use_tls=aiexec_config.SMTP_USE_TLS,
                    opportunistic_tls=aiexec_config.SMTP_OPPORTUNISTIC_TLS,
                )
            case _:
                raise ValueError("Unsupported mail type {}".format(mail_type))

    def send(self, to: str, subject: str, html: str, from_: Optional[str] = None):
        if not self._client:
            raise ValueError("Mail client is not initialized")

        if not from_ and self._default_send_from:
            from_ = self._default_send_from

        if not from_:
            raise ValueError("mail from is not set")

        if not to:
            raise ValueError("mail to is not set")

        if not subject:
            raise ValueError("mail subject is not set")

        if not html:
            raise ValueError("mail html is not set")

        self._client.send(
            {
                "from": from_,
                "to": to,
                "subject": subject,
                "html": html,
            }
        )


def is_enabled() -> bool:
    return aiexec_config.MAIL_TYPE is not None and aiexec_config.MAIL_TYPE != ""


def init_app(app: AiexecApp):
    mail.init_app(app)


mail = Mail()
