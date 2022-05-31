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

		url, err := serviceMap[key].GetPlayUrl(roomId)
		if err != nil {
			fmt.Println(err.Error())
			c.String(200, err.Error())
			return
		}

		c.Redirect(http.StatusFound, url)
	})
}
