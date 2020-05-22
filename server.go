package main

import (
	"ws101/router"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r = router.Router(r)

	r.Run()
}
