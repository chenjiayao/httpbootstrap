package routes

import (
	"httptemplate/controllers"
	"httptemplate/globals"
	"httptemplate/middlewares"
)

var (
	regionController = &controllers.IndexController{}
)

func LoadRoute() {

	g1 := globals.Engine.Use(middlewares.TestMiddleware)
	{
		g1.GET("/ping", regionController.Test)
	}
}
