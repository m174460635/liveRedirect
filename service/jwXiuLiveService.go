package service

import (
	"errors"
	"github.com/asmcos/requests"
	"regexp"
)

type JwXiuLiveService struct {

}

func (s *JwXiuLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "http://www.95.cn/" + key + ".html"
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	resText := res.Text()
	statusMatchRes := regexp.MustCompile(`"is_offline":"(\d)"`).FindStringSubmatch(resText)
	if statusMatchRes == nil {
		return "", nil
	}
	status := statusMatchRes[1]
	if status != "0" {
		return "", errors.New("未开播")
	}
	uidMatchRes := regexp.MustCompile(`"uid":(\d+),`).FindStringSubmatch(resText)
	if uidMatchRes == nil {
		return "", nil
	}
	uid := uidMatchRes[1]
	real_url := "http://play.95xiu.com/app/" + uid + ".flv"
	return real_url, nil
}

func init() {
	RegisterService("95xiu", new(JwXiuLiveService))
}
