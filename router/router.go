package router

import (
	"net/http"
	"ws101/views"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) *gin.Engine {
	r.Use(cors.Default())
	r.GET("/hi", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok": 200,
		})
	})

	t := views.ToDoView{}
	toDoRouters := r.Group("todo")
	toDoRouters.GET("", t.Load)
	toDoRouters.GET("/ws", t.Ws)
	toDoRouters.POST("", t.Add)
	toDoRouters.PUT(":id", t.Update)
	toDoRouters.DELETE(":id", t.Delete)

	return r
}
