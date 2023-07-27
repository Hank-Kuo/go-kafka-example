import json
import logging

from kafka import KafkaConsumer, KafkaProducer
from kafka.admin import KafkaAdminClient, NewTopic

class MessageConsumer:
    def __init__(self, broker, topic, group_id, cfg):
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


class MessageProducer:
    def __init__(self, cfg):
        self.broker = cfg["broker"]
        self.topic = cfg["topic"] 
        self.producer = KafkaProducer(
            bootstrap_servers = self.broker,
            value_serializer = lambda v: json.dumps(v).encode('utf-8'),
            acks=cfg["producer_acks"], retries = cfg["producer_retries"],
        )
        logging.info("Connect with kafka server...")


    def delivery_message(self, data):
        logging.info("[Producer] Send messages: {}".format(str(data)))
        try:
            future = self.producer.send(self.topic, data)
            self.producer.flush()
            future.get()
        except Exception as err:
            logging.error("[Producer][Fail] Send messages: {}".format(str(err)))
   


def create_topics(topic_names):
    
    admin_client = KafkaAdminClient(bootstrap_servers=brokers)
    topic_list = []
    topic_list.append(NewTopic(name="example_topic", num_partitions=1, replication_factor=1))

    admin_client.create_topics(new_topics=topic_list, validate_only=False)

    