import asyncio
import threading
import pyotp
import logging

from kafka_py.utils.email import Email
from kafka_py.utils.otp import get_token
from kafka_py.config.config import CONFIG


def mail_job(message_consumer):
    def msg_process(data):
        try:
            email_server = Email(CONFIG["mail_smtp"], CONFIG["mail_smtp_port"], CONFIG["mail_from"], CONFIG["mail_password"])
            email_server.connect()

            name = data["name"]
            to_email = data["email"]
            token = get_token(to_email)
            
            otp_code = pyotp.TOTP(token, interval=60*15)
            email_server.send_mail(to_email, name, otp_code.now())
            email_server.quit()
            return True
        except Exception as err:
            logging.info(err)
            return False

    message_consumer.listening(msg_process)
    

