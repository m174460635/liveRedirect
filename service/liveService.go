package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/markbates/pkger"
	"io/ioutil"
)

var serviceMap = make(map[string]LiveService)

func initServiceMap() map[string]LiveService {
	//服务列表
	serviceMap["huya"] = new(HuyaLiveService)
	serviceMap["yy"] = new(YYLiveService)
	serviceMap["huajiao"] = new(HuajiaoLiveService)
	serviceMap["2cp"] = new(SpunSugarLiveService)
	serviceMap["zhanqi"] = new(ZhanqiLiveService)
	serviceMap["kugou"] = new(KugouLiveService)
	serviceMap["douyu"] = new(DouyuLiveService)
	serviceMap["51lm"] = new(LMLiveService)
	//serviceMap["iqiyi"] = new(IqiyiLiveService)
	return serviceMap
}
func GetServiceMap() map[string]LiveService {
	if len(serviceMap) == 0 {
		initServiceMap()
	}
	return serviceMap

}

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
