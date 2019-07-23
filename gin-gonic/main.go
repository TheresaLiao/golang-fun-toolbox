package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("start api")
	router := gin.Default()

	// curl  http://localhost:80/
	router.GET("/", test)
	router.GET("/getJson", test)
	router.Run(":80")
}


func test(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello",})
}