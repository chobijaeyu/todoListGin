package views

import (
	"log"
	"net/http"
	"ws101/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type ToDoView struct{}

func (t ToDoView) Ws(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	// ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrader error:$s\n", err.Error())
		return
	}
	defer log.Println("ws closed")
	defer ws.Close()
	var (
		mt      = 1
		todo    = models.ToDo{}
		db, ctx = models.GetClient()
		docChan = make(chan []byte, 10)
	)

	defer db.Disconnect(ctx)

	go todo.WsRecord(db, docChan)
	log.Println("ws db")

	go func() {
		for {
			mt, _, err = ws.ReadMessage()
			if err != nil {
				log.Println("ws read message err", err.Error())
				break
			}
		}
	}()

	for {
		d, ok := <-docChan
		if !ok {
			break
		}
		err = ws.WriteMessage(mt, d)
		if err != nil {
			log.Println(err.Error())
			break
		}
	}
}

func (t ToDoView) Load(c *gin.Context) {
	todo := models.ToDo{}

	db, ctx := models.GetClient()
	defer db.Disconnect(ctx)

	q := c.Query("query")
	p := c.Query("param")

	if q == "" && p == "" {
		log.Printf("not query get request")
		r, err := todo.LoadRecord(db)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"err": err.Error(),
			})
			return
		}
		if len(r) == 0 {
			c.JSON(http.StatusNoContent, gin.H{
				"err": "no content",
			})
			return
		}

		c.JSON(http.StatusOK, r)
		return
	}
	r, err := todo.QueryRecord(db, q, p)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	if len(r) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"err": "no content",
		})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (t ToDoView) Add(c *gin.Context) {
	todo := models.ToDo{}
	if err := c.Bind(&todo); err != nil {
		log.Printf("Failed bindjson with todo: %v\n", err)
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
	todo.ID = primitive.NewObjectID()
	_, err := todo.AddRecord(db, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, todo)

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
