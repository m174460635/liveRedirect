package service

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type HuyaLiveService struct {
}

func live(e string) string {
	i := strings.Split(e, "?")[0]
	b := strings.Split(e, "?")[1]
	r := strings.Split(i, "/")
	re := regexp.MustCompile(".(flv|m3u8)")
	s := re.ReplaceAllString(r[len(r)-1], "")
	c := strings.SplitN(b, "&", 4)
	cc := c[:0]
	n := make(map[string]string)
	for _, x := range c {
		if len(x) > 0 {
			cc = append(cc, x)
			ss := strings.Split(x, "=")
			n[ss[0]] = ss[1]
		}
	}
	c = cc
	fm, _ := url.QueryUnescape(n["fm"])
	uu, _ := base64.StdEncoding.DecodeString(fm)
	u := string(uu)
	p := strings.Split(u, "_")[0]
	f := strconv.FormatInt(time.Now().UnixNano()/100, 10)
	l := n["wsTime"]
	t := "0"
	h := p + "_" + t + "_" + s + "_" + f + "_" + l
	m := GetMD5Hash(h)
	y := c[len(c)-1]
	url := fmt.Sprintf("%s?wsSecret=%s&wsTime=%s&u=%s&seqid=%s&%s", i, m, l, t, f, y)

	return url
}
func (s *HuyaLiveService) GetPlayUrl(key string) (string, error) {
	//roomUrl := "https://m.huya.com/" + key
	//resp, err := requests.Get(roomUrl, requests.Header{"Content-Type": "application/x-www-form-urlencoded",
	//	"User-Agent": "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36 ",
	//})
	//if err != nil {
	//	fmt.Print(err.Error())
	//	return "", err
	//}
	//pageResult := resp.Text()
	//re := regexp.MustCompile(`"liveLineUrl":"([\s\S]*?)",`)
	//res := re.FindStringSubmatch(pageResult)
	//if len(res) > 0 { //有直播链接
	//	u := res[1]
	//	if len(u) > 0 {
	//		decodedRet, _ := base64.StdEncoding.DecodeString(u)
	//		decodedUrl := string(decodedRet)
	//		if strings.Contains(decodedUrl, "replay") { //重播
	//			return "https:" + u, nil
	//		} else {
	//			liveLineUrl := live(decodedUrl)
	//			liveLineUrl = strings.Replace(liveLineUrl, "hls", "flv", -1)
	//            		liveLineUrl = strings.Replace(liveLineUrl, "m3u8", "flv", -1)
	//			return "https:" + liveLineUrl, nil
	//		}
	//	}
	//}
	urls, err := GetHuyaStreamUrls("https://www.huya.com/" + key)
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	for _, v := range urls {
		if strings.Contains(v, "http://al") {
			return v, nil
		}
	}
	for _, v := range urls {
		if strings.HasPrefix(v, "http") {
			return v, nil
		}
	}
	return "", nil
}
func init() {
	RegisterService("huya", new(HuyaLiveService))
}
