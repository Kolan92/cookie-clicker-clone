package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kolan92/hunter/game"
)

type GameController struct {
	resourcesStock *game.ResourcesStock
}

type UpgradeRequest struct {
	Resource string `json:"resource" binding:"required"`
}

func NewGameController(resourcesStock *game.ResourcesStock) *GameController {
	return &GameController{
		resourcesStock,
	}
}

func (gc *GameController) RegisterRouter(routerGroup *gin.RouterGroup) {
	game := routerGroup.Group("/game")
	{
		game.GET("/dashboard", func(c *gin.Context) {
			gc.GetDashboard(c)
		})

		game.POST("/upgrade", func(c *gin.Context) {
			gc.UpgradeFactory(c)
		})
	}
}

func (gc *GameController) GetDashboard(g *gin.Context) {

	g.JSON(http.StatusOK, gc.resourcesStock)
}

func (gc *GameController) UpgradeFactory(g *gin.Context) {

	input := UpgradeRequest{}

	err := g.ShouldBindJSON(&input)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"message": "enum is not valid"})
		return
	}

	resourceName := strings.ToLower(input.Resource)
	resource, isFound := game.ResourceFromString[resourceName]
	if !isFound {
		g.JSON(http.StatusBadRequest, gin.H{"message": "enum is not valid"})
		return
	}

	err = gc.resourcesStock.Upgrade(resource)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully upgraded %s", resourceName)})
}
