import json
import logging

from confluent_kafka import Producer, Consumer, KafkaError, KafkaException
# from kafka import KafkaConsumer, KafkaProducer
# from kafka.admin import KafkaAdminClient, NewTopic


"""
class MessageProducer1:
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
"""


class MessageConsumer:
    def __init__(self, cfg):
        conf = {
            "bootstrap.servers": ",".join(cfg["broker"]),
            "group.id": cfg["group_id"],
            "auto.offset.reset": "smallest"
        }
        consumer = Consumer(conf)
        self.consumer = consumer
        self.consumer.subscribe(cfg["topic"])
    
    def listening(self, msg_process):
        logging.info("[Consumer] start listening...")
        try:
            while True:
                message = self.consumer.poll(timeout=1.0)
                if message is None: 
                    continue 
                if message.error():
                    if message.error().code() == KafkaError._PARTITION_EOF:
                        pass
                    elif message.error():
                        raise KafkaException(message.error())
                else:
                    data = json.loads(message.value().decode('utf-8'))
                    logging.info("[Consumer] Received message: {}".format(data))
                    status = msg_process(data)
                    if status:
                        self.consumer.commit(asynchronous=False)

        except KeyboardInterrupt:
            logging.error("[Consumer] close server")
        finally:
            self.consumer.close()

        

class MessageProducer:
    def __init__(self, cfg):
        conf = {
            "bootstrap.servers": ",".join(cfg["broker"]),
            "acks": cfg["producer_acks"]
        }
        self.producer = Producer(conf)
        self.topic = cfg["topic"] 
        logging.info("[Producer] connect with kafka server...")
    
    def acked(self, err, msg):
        if err is not None:
            logging.info("[Producer][Fail] can't deliver message: %s: %s" % (str(msg), str(err)))
        else:
            logging.info("[Producer] Message produced: %s" % (str(msg)))

    def delivery_message(self, data):
        logging.info("[Producer] Send messages: {}".format(str(data)))
        try:
            self.producer.produce(
                self.topic, value=json.dumps(data).encode('utf-8')    
            )
            self.producer.flush()
            
        except Exception as err:
            logging.error("[Producer][Fail] Send messages: {}: {}".format(str(data), str(err)))
   

"""
def create_topics(topic_names, partition_num, replication_factor):
    admin_client = KafkaAdminClient(bootstrap_servers=brokers)
    topic = NewTopic(name=topic_names, num_partitions=partition_num, replication_factor=replication_factor)
    admin_client.create_topics(new_topics=[topic], validate_only=False)

"""