package tasks

import "github.com/gin-gonic/gin"

func Routes(route *gin.Engine) {
	r := route.Group("/tasks")
	{
		r.GET("/result/:UUID", GetResult)
		r.GET("/log/:UUID", GetLog)
		r.GET("/status/:UUID", GetStatus)
	}
}
