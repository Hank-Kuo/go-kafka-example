import re
import base64
from datetime import datetime

def remove_non_alphabetic(input_string):

    s = input_string.lower() + encode_datetime()
    pattern = r'[^a-z]+'
    cleaned_string = re.sub(pattern, '', s)
    return cleaned_string

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

def encode_datetime():
    date = datetime.utcnow().strftime("%Y%m%d")
    s = ""
    for x in date :
        if x in dct:
            s += dct[x]
    return s
