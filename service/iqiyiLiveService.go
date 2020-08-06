package service

import (
	"errors"
	"fmt"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"github.com/markbates/pkger"
	"io/ioutil"
	"net/url"
	"regexp"
	"rogchap.com/v8go"
	"strings"
	"time"
)

type IqiyiLiveService struct {
}

func (s *IqiyiLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://m-gamelive.iqiyi.com/w/" + key
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	responseText := res.Text()
	reg := regexp.MustCompile(`"qipuId":(\d*?),"roomId`)
	matchedStrs := reg.FindStringSubmatch(responseText)
	if matchedStrs == nil {
		return "", errors.New("qipuId find error")
	}
	qipuId := matchedStrs[1]
	callback := "jsonp_" + fmt.Sprintf("%d", time.Now().UnixNano()/1000000) + "_0000"
	params := url.Values{}

	params.Add("lp", qipuId)
	params.Add("src", "01010031010000000000")
	params.Add("rateVers", "H5_QIYI")
	params.Add("qd_v", "1")
	params.Add("callback", callback)
	ba := "/jp/live?" + params.Encode()
	f, err := pkger.Open("/resources/js/iqiyi.js")
	if err != nil {
		return "", err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	iqiyijsText := string(content)
	ctx, _ := v8go.NewContext(nil)
	ctx.RunScript(iqiyijsText, "test.js")
	vf, err := ctx.RunScript("cmd5x('"+ba+"')", "test.js")
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	p := requests.Params{
		"vf": vf.String(),
	}
	res, err = requests.Get("https://live.video.iqiyi.com"+ba, p)
	if err != nil {
		return "", err
	}
	responseText = res.Text()
	reg = regexp.MustCompile(`try{.*?\((.*)\);}catch\(e\){};`)
	urlJsons := reg.FindStringSubmatch(responseText)
	if urlJsons == nil {
		return "", nil
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	urlJsonData := make(map[string]interface{})
	if err = json.UnmarshalFromString(urlJsons[1], &urlJsonData); err != nil {
		return "", err
	}
	realUrl := json.Get([]byte(urlJsons[1]), "data", "streams", 0, "url").ToString()
	realUrl = strings.Replace(realUrl, "hlslive.video.iqiyi.com", "m3u8live.video.iqiyi.com", 1)
	return realUrl, nil
}
