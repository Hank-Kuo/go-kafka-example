import os 
import smtplib
from dotenv import load_dotenv
from kafka import KafkaProducer
import threading
import json

class MessageProducer:
    broker = ""
    topic = ""
    producer = None

    def __init__(self, broker, topic):
        self.broker = broker
        self.topic = topic
        self.producer = KafkaProducer(
            bootstrap_servers = self.broker,
        
            value_serializer = lambda v: json.dumps(v).encode('utf-8'),
            
            acks='all', retries = 3,
        )


    def send_msg(self, data):
        print("sending message...")
        try:
            future = self.producer.send(self.topic, data)
            self.producer.flush()
            future.get()
            print("message sent successfully...")
            return {'status_code': 200, 'error': None }
        except Exception as ex:
            return ex

load_dotenv()
broker = 'localhost:9092'
topic = 'topic1'
message_producer = MessageProducer(broker, topic)
data = {'name':'abc', 'email':'abc@example.com'}
resp = message_producer.send_msg(data)
print(resp)
