
CONFIG = {
    "broker": ["localhost:9092"],
    "topic": "user_email", 
    "producer_acks": "all",
    "producer_retries": 3,

    "group_id": "consumer_python", 
    "consumer_timeout": 60000,
    "enable_auto_commit": False,
    "auto_offset_reset": "earliest",

    "logger_path": "./assets/run.log"
}

