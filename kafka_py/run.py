import logging
from kafka_py.adapter.kafka_adapter import MessageProducer, MessageConsumer
from kafka_py.config.config import  CONFIG

from kafka_py.utils import logger 



if __name__ == "__main__":
    logger.set_logger(CONFIG["logger_path"])
    message_producer = MessageProducer(CONFIG)
    