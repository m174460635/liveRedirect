package main

import (
	"fmt"
	"github.com/asmcos/requests"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"liveRedirect/service"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	enableCacheSet = hashset.New()
	// 设置超时时间和清理时间
	ccache = cache.New(20*time.Second, 21*time.Second)
)

func setupEnableCache() {
	enableCacheSet.Add("huya")
}
func main() {
	//设置允许缓存的直播平台
	setupEnableCache()

	//初始化服务列表
	serviceMap := service.GetServiceMap()

	//启动web服务
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	handleService(r, serviceMap)
	r.Run(":5000")
}

func handleService(r *gin.Engine, serviceMap map[string]service.LiveService) gin.IRoutes {
	return r.GET("/:key/:id", func(c *gin.Context) {
		key := c.Param("key")
		roomId := c.Param("id")

		_, ok := serviceMap[key]
		if !ok {
			key = "huya"
		}

		uurl := ""
		if enableCacheSet.Contains(key) {
			kkey := key + "__" + roomId
			if x, found := ccache.Get(kkey); found {
				uurl = x.(string)
			} else {
				url, err := serviceMap[key].GetPlayUrl(roomId)
				if err != nil {
					fmt.Println(err.Error())
					c.String(200, err.Error())
					return
				}
				uurl = url
				ccache.Set(kkey, url, time.Second*20)
			}
		} else {
			url, err := serviceMap[key].GetPlayUrl(roomId)
			if err != nil {
				fmt.Println(err.Error())
				c.String(200, err.Error())
				return
			}
			uurl = url
		}

		if key == "huya" {
			processHuya(c, uurl)
			return
		}

		c.Redirect(http.StatusFound, uurl)
	})
}

func processHuya(c *gin.Context, url string) bool {
	i := strings.LastIndex(url, "/")
	if i > -1 {
		urlPrefix := url[0 : i+1]
		resp, err := requests.Get(strings.TrimSpace(url), requests.Header{"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent": "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36 ",
		})

		if err != nil {
			fmt.Println(err)
			c.Redirect(http.StatusFound, url)
			return true
		}
		ss := strings.Split(resp.Text(), "\n")
		s := ""
		for _, v := range ss {
			if strings.HasPrefix(v, "#") {
				s = s + v + "\r\n"
			} else {
				s = s + urlPrefix + v + "\r\n"
			}
		}
		c.Header("Content-type", "application/vnd.apple.mpegurl")
		c.Header("Content-Length", strconv.Itoa(len([]rune(s))))
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")
		c.Writer.WriteString(s)

		return true
	}
	return false
}
