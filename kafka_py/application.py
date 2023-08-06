import asyncio
from fastapi import FastAPI

from kafka_py.extensions.handler_extension import register_handler
from kafka_py.jobs.mail_job import mail_job

app = FastAPI()

register_handler(app)
