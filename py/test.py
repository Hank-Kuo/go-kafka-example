import os 
import smtplib
from dotenv import load_dotenv

load_dotenv()

MAIL_SMTP = os.getenv("MAIL_SMTP")
MAIL_SMTP_PORT = os.getenv("MAIL_SMTP_PORT")
MAIL_FROM = os.getenv("MAIL_FROM")
MAIL_PASSWORD = os.getenv("MAIL_PASSWORD")
MAIL_TO = os.getenv("MAIL_TO").split(",")

def connect(smtp, port, mail_from, mail_password):

    smtp = smtplib.SMTP(smtp, port)
    smtp.ehlo()
    smtp.starttls()
    smtp.login(mail_from, mail_password)
    return smtp

smtp = connect(MAIL_SMTP, MAIL_SMTP_PORT, MAIL_FROM, MAIL_PASSWORD)

SUBJECT = "Get mails by Gmail"
BODY = "Subscribe!!!"

MESSAGE = """\
From: %s
To: %s
Subject: %s

%s
""" % (MAIL_FROM, ", ".join(MAIL_TO), SUBJECT, BODY)


status = smtp.sendmail(MAIL_FROM, MAIL_TO, MESSAGE)

if status=={}:
    print("郵件傳送成功!")
else:
    print("郵件傳送失敗!")

smtp.quit()
