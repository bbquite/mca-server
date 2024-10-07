package utils

import (
	"crypto/hmac"
	"crypto/sha256"
)

func MakeHMACSign(key string, data []byte) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	sign := h.Sum(nil)
	return sign
}

func CheckHMACEqual(key string, sign []byte, data []byte) bool {
	return hmac.Equal(MakeHMACSign(key, data), sign)
}
