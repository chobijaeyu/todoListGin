package views

import (
	"log"
	"net/http"
	"ws101/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type ToDoView struct{}

func (t ToDoView) Ws(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	// ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	// goods := models.Goods{}
	// goods, err = goods.GetAllGoods(ctx, client)
	if err != nil {

		return
	}
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func (t ToDoView) Load(c *gin.Context) {
	todo := models.ToDo{}

	db, ctx := models.GetClient()
	defer db.Disconnect(ctx)

	r, err := todo.LoadRecord(db)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	if len(r) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "no content",
		})
		return
	}

	c.JSON(http.StatusOK, r)
}

func (t ToDoView) Add(c *gin.Context) {
	todo := models.ToDo{}
	if err := c.Bind(&todo); err != nil {
		log.Printf("Failed bindjson with Goods: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Error Get Body",
		})
		return
	}

	if todo.Desc == "" || todo.Done != false {
		log.Printf("new todo invalid data")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "new todo invalid data",
		})
		return
	}

	db, ctx := models.GetClient()
	defer db.Disconnect(ctx)

	r, err := todo.AddRecord(db, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"res": r,
	})

}

func (t ToDoView) Update(c *gin.Context) {
	todo := models.ToDo{}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "update must have id",
		})
		return
	}

	db, ctx := models.GetClient()
	defer db.Disconnect(ctx)

	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Error Get Body",
		})
		return
	}

	r, err := todo.UpateRecord(db, todo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"res": r,
	})

}

func (t ToDoView) Delete(c *gin.Context) {
	todo := models.ToDo{}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "update must have id",
		})
		return
	}

	db, ctx := models.GetClient()
	defer db.Disconnect(ctx)

	r, err := todo.DeleteRecord(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"res": r,
	})
}
