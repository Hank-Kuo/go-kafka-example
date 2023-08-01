import re

def remove_non_alphabetic(input_string):
    pattern = r'[^a-z]+'
    cleaned_string = re.sub(pattern, '', input_string.lower())
    return cleaned_string