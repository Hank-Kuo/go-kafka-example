import logging
import smtplib
from email.mime.text import MIMEText

class Email:
    def __init__(self, smtp_host, smtp_port, mail_from, mail_password):
        self.smtp_host = smtp_host
        self.smtp_port = smtp_port
        self.mail_from = mail_from
        self.mail_password = mail_password
        self.smtp = None
    
    def connect(self):
        smtp = smtplib.SMTP(self.smtp_host, self.smtp_port)
        smtp.ehlo()
        smtp.starttls()
        smtp.login(self.mail_from, self.mail_password)
        self.smtp = smtp

    def send_mail(self, mail_to, name, opt_code):
        content = """
        <html>
            <body>
        <br/>
        Hi <b>%s</b>,<br/> <br/>
        You're trying to register this server, please use this OTP code (one-time password) to active your account.<br/><br/><br/>
        The OTP code is: <font size="20px"><b>%s</b></font><br/>
        <br/>
        Please enter the code in <b>15 mins</b>.
        </body>
        </html>
        """ % (name, opt_code)
        msg = MIMEText(content, 'html')

        msg['Subject'] = "Activate your account"
        msg['From'] = self.mail_from
        msg['To'] = mail_to
        
        status = self.smtp.sendmail(self.mail_from, mail_to, msg.as_string())
        
        if status!={}:
            raise Exception("[SMTP][Fail] can't send email: {}".format(str(status)))
    
    def quit(self):
        self.smtp.quit()


if __name__ == "__main__":
    email = Email(_, _, _, _)
    email.connect()
    try:
        email.sendmail(_)
    except Exception as err:
        logging.error(err)

    email.quit()


