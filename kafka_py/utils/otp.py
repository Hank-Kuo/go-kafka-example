import re
import base64
from datetime import datetime

def get_token(input_string):
    s = input_string.lower() + encode_datetime()
    pattern = r'[^a-z]+'
    cleaned_string = re.sub(pattern, '', s)
    return cleaned_string

def encode_datetime():
    dct = {
        0: "a",
        1: "b",
        2: "c",
        3: "d",
        4: "e",
        5: "f",
        6: "g",
        7: "h",
        8: "i",
        9: "j",
    }
    
    date = datetime.utcnow().strftime("%Y-%m-%d")
    s = ""
    for _x in date:
        try:
            x = int(_x)
            if x in dct:
                s += dct[x]
        except:
            continue
    return s
