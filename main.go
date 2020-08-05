package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"liveRedirect/service"
	"net/http"
)

func initServiceMap() map[string]service.LiveService {
	//服务列表
	serviceMap := make(map[string]service.LiveService)
	serviceMap["huya"] = new(service.HuyaLiveService)
	serviceMap["yy"] = new(service.YYLiveService)
	serviceMap["huajiao"] = new(service.HuajiaoLiveService)
	serviceMap["2cp"] = new(service.SpunSugarLiveService)
	return serviceMap
}
func main() {
	//初始化服务列表
	serviceMap := initServiceMap()

	//启动web服务
	r := gin.Default()
	r.GET("/:key/:id", func(c *gin.Context) {
		key := c.Param("key")
		roomId := c.Param("id")

		_, ok := serviceMap[key]
		if !ok {
			key = "huya"
		}
		url, err := serviceMap[key].GetPlayUrl(roomId)
		if err != nil {
			fmt.Println(err.Error())
			c.String(200, err.Error())
			return
		}
		fmt.Println(url)
		c.Redirect(http.StatusMovedPermanently, url)
	})
	r.Run(":5000")
}