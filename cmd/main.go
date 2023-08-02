package main

import (
	"strconv"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GeneratePassCode(s string) (string, error) {
	timeStr, err := encode_datetime()
	if err != nil {
		return "", err
	}

	otpCode, err := totp.GenerateCodeCustom(s+timeStr, time.Now(), totp.ValidateOpts{
		Period:    60 * 15,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	return otpCode, err
}

func encode_datetime() (string, error) {
	ALPHABET := map[int]string{
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

	_t := time.Now().UTC()
	_timeStr := _t.Format("20060102")
	timeStr := ""
	for _, c := range _timeStr {
		k, err := strconv.Atoi(string(c))
		if err != nil {
			return "", err
		}
		alphabet, ok := ALPHABET[k]
		if ok {
			timeStr += alphabet
		}
	}
	return timeStr, nil
}
