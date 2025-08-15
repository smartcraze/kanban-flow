package main

import "github.com/gin-gonic/gin"

func main() {
	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	app.Run(":8000")

}
