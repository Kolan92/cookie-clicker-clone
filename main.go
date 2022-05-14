package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kolan92/hunter/controllers"
	"github.com/kolan92/hunter/game"
)

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	game := game.NewGame()
	controller := controllers.NewGameController(game)

	controller.RegisterRouter(v1)

	router.Run(":8082")
}
