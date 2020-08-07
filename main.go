package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"liveRedirect/service"
	"net/http"
)

func main() {
	//初始化服务列表
	serviceMap := service.GetServiceMap()

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
