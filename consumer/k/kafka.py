from kafka import KafkaProducer
import json


class MessageProducer:
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


class MessageConsumer:
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
                                 value_deserializer=lambda m: m.decode('utf-8'))#json.loads(

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
