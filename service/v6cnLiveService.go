package service

import (
	"fmt"
	"github.com/asmcos/requests"
	"regexp"
)

type V6CNLiveService struct{}

func (V6CNLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://v.6.cn/" + key
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	uidMatchRes := regexp.MustCompile(`"flvtitle":"v(\d*?)-(\d*?)"`).FindStringSubmatch(res.Text())
	if uidMatchRes == nil {
		return "", nil
	}
	println(uidMatchRes[1])
	println(uidMatchRes[2])
	uid := uidMatchRes[1]
	flvTitle := fmt.Sprintf("v%s-%s", uidMatchRes[1], uidMatchRes[2])
	realUrl := "https://rio.6rooms.com/live/?s=" + uid
	res, err = requests.Get(realUrl)
	if err != nil {
		return "", err
	}
	hipMatchRes := regexp.MustCompile(`<watchip>(.*\.xiu123\.cn).*</watchip>`).FindStringSubmatch(res.Text())
	if hipMatchRes == nil {
		return "", nil
	}
	hip := "https://" + hipMatchRes[1]
	return hip + "/" + flvTitle + "/playlist.m3u8", nil
}

func init() {
	RegisterService("v6cn", new(V6CNLiveService))
}
