import os
from dotenv import load_dotenv

load_dotenv()

CONFIG = {
    "host": os.getenv("CONSUMER_HOST"),
    "port": int(os.getenv("CONSUMER_PORT")),
    "reload": (os.getenv('CONSUMER_DEBUG', 'false') == 'true'),
    "worker": int(os.getenv("CONSUMER_WORKER")),

    "broker": os.getenv("KAFKA_BROKER").split(","),
    "topic": ["user_email"], 
    "producer_acks": "all",
    "producer_retries": 3,

    "group_id": "consumer_python", 
    "consumer_timeout": 60000,
    "enable_auto_commit": False,
    "auto_offset_reset": "earliest",

    "logger_path": "./assets/run.log",

    "mail_smtp": os.getenv("MAIL_SMTP"),
    "mail_smtp_port": os.getenv("MAIL_SMTP_PORT"),
    "mail_from": os.getenv("MAIL_FROM"),
    "mail_password": os.getenv("MAIL_PASSWORD"),
    "mail_to": os.getenv("MAIL_TO"),
}


