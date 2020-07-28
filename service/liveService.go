package service

import (
	"crypto/md5"
	"encoding/hex"
)

type LiveService interface {
	GetPlayUrl(key string) (string, error)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
