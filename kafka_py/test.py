import logging
import threading
import pyotp

from kafka_py.adapter.kafka_adapter import MessageProducer, MessageConsumer
from kafka_py.config.config import  CONFIG
from kafka_py.utils import logger 
from kafka_py.utils.email import Email
from kafka_py.utils.otp import get_token

token = get_token("hank_kuo@trendmicro.com")
print(token)
otp_code = pyotp.TOTP(token, interval=60*15)
print(otp_code.now())