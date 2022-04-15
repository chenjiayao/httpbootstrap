package controllers

import "github.com/gin-gonic/gin"

type IndexController struct {
}

func (c IndexController) Test(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
