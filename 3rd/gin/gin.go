// Copyright © 2022 hops. All rights reserved.

package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func helloGin(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "hello gin",
	})
}

func main() {
	fmt.Println("hello gin")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/hello", helloGin)
	router.Run()
}
