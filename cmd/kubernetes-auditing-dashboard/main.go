package main

import (
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	apiGroup := app.Group("/api")
	apiGroup.POST("/audit-webhook", func(c *gin.Context) {
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestBody))
		c.Status(200)
	})
	app.Run("0.0.0.0:23333")
}
