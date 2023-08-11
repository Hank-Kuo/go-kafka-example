import threading
import uvicorn

from kafka_py.adapter.kafka_adapter import MessageConsumer
from kafka_py.config.config import CONFIG
from kafka_py.jobs.mail_job import mail_job
from kafka_py.utils import logger 


if __name__ == "__main__":
    logger.set_logger(CONFIG["logger_path"])
    # background
    message_consumer = MessageConsumer(CONFIG)
    t1 = threading.Thread(target=mail_job, daemon=True, args=(message_consumer,))
    t1.start()

    uvicorn.run("kafka_py.application:app", host=CONFIG["host"], port=CONFIG["port"], reload=CONFIG["reload"], workers=CONFIG["worker"])
    
    message_consumer.consumer.close()