package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateOTP(inputString string) (string, error) {
	otpCode, err := totp.GenerateCodeCustom(inputString, time.Now(), totp.ValidateOpts{
		Period:    60 * 15,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	return otpCode, err
}

func ValidOTP(otpCode string, inputString string) (bool, error) {
	isValid, err := totp.ValidateCustom(otpCode, inputString, time.Now(), totp.ValidateOpts{
		Period:    60 * 15,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return isValid, err
}

func GetToken(inputString string) string {
	s := strings.ToLower(inputString) + encode_datetime()

	pattern := regexp.MustCompile(`[^a-z]+`)
	cleanedString := pattern.ReplaceAllString(s, "")
	return cleanedString
}

func encode_datetime() string {
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
			continue
		}
		alphabet, ok := ALPHABET[k]
		if ok {
			timeStr += alphabet
		}
	}
	return timeStr
}
