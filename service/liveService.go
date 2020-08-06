package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/markbates/pkger"
	"io/ioutil"
)

type LiveService interface {
	GetPlayUrl(key string) (string, error)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func GetFromResource(path string) (string, error) {
	f, err := pkger.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	text := string(content)
	return text, nil
}
