import os 
import smtplib
from dotenv import load_dotenv
from kafka import KafkaConsumer
import threading
import json

load_dotenv()

BOOTSTRAP_SERVERS = ['localhost:9092']

def register_kafka_listener(topic, listener):
    def poll():
        # Initialize consumer Instance
        consumer = KafkaConsumer(topic, bootstrap_servers=BOOTSTRAP_SERVERS)

        print("About to start polling for topic:", topic)
        consumer.poll(timeout_ms=6000)
        print("Started Polling for topic:", topic)
        for msg in consumer:
            print("Entered the loop\nKey: ",msg.key," Value:", msg.value)
            kafka_listener(msg)
    print("About to register listener to topic:", topic)
    t1 = threading.Thread(target=poll)
    t1.start()
    print("started a background thread")

def kafka_listener(data):
    print("Image Ratings:\n", data.value.decode("utf-8"))


class MessageConsumer:
    broker = ""
    topic = ""
    group_id = ""
    logger = None

    def __init__(self, broker, topic, group_id):
        self.broker = broker
        self.topic = topic
        self.group_id = group_id

    def activate_listener(self):
        consumer = KafkaConsumer(bootstrap_servers=self.broker,
                                 group_id=group_id,
                                 consumer_timeout_ms=60000,
                                 auto_offset_reset='earliest',
                                 enable_auto_commit=False,
                                 value_deserializer=lambda m: json.loads(m.decode('utf-8')))

        consumer.subscribe(self.topic)
        print("consumer is listening....")
        try:
            for message in consumer:
                print("received message = ", message)
                consumer.commit()
        except KeyboardInterrupt:
            print("Aborted by user...")
        finally:
            consumer.close()

# register_kafka_listener('topic1', kafka_listener)

broker = 'localhost:9092'
topic = 'topic1'
group_id = 'consumer-1'

consumer1 = MessageConsumer(broker,topic,group_id)
consumer1.activate_listener()
