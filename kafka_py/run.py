import logging
import threading
import pyotp


from kafka_py.adapter.kafka_adapter import MessageProducer, MessageConsumer
from kafka_py.config.config import  CONFIG
from kafka_py.utils import logger 
from kafka_py.utils.email import Email
from kafka_py.utils.normalize import remove_non_alphabetic

def main(consumer, email_server):
    logging.info("[Consumer] start listening...")

    try:
        for message in consumer:
            if message is None:
                continue
            logging.info("[Consumer] Received message: {}".format(message))
            name = message.value["name"]
            to_email = message.value["email"]
            normalize_to_mail = remove_non_alphabetic(to_email)
            otp_code = pyotp.TOTP(normalize_to_mail, interval=60*15)
            email_server.send_mail(to_email, name, otp_code.now())
            consumer.commit()
    except KeyboardInterrupt:
        logging.error("[Consumer] close server")
    finally:
        email_server.quit()
        consumer.close()

if __name__ == "__main__":
    logger.set_logger(CONFIG["logger_path"])

    email_server = Email(CONFIG["mail_smtp"], CONFIG["mail_smtp_port"], CONFIG["mail_from"], CONFIG["mail_password"])
    email_server.connect()
    logging.info("[Consumer] connect SMTP server")

    # message_producer = MessageProducer(CONFIG)
    message_consumer = MessageConsumer(CONFIG)
    consumer = message_consumer.activate_listener()
    
    # run in background
    t1 = threading.Thread(target=main, args=(consumer, email_server))
    t1.start()