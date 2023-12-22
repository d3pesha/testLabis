package main

import (
	"github.com/d3pesha/testLabis/database"
	_ "github.com/d3pesha/testLabis/docs"
	"github.com/d3pesha/testLabis/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @version 1.0
// @description API for managing clients, objects, and contracts.

// @host localhost:8080
// @BasePath /api/v1

// @title Medium API
func main() {
	router := gin.Default()

	config := database.NewPostgresConfig()
	postgres := database.Connect(config)

	clientService := handlers.NewClientService(postgres.GetDB())
	objectService := handlers.NewObjectService(postgres.GetDB())
	contractService := handlers.NewContractService(postgres.GetDB())

	api := router.Group("/api/v1")
	{
		api.GET("/clients", clientService.GetClients)
		api.GET("/clients/:id", clientService.GetClientByID)
		api.POST("/clients", clientService.CreateClient)
		api.DELETE("/clients/:id", clientService.DeleteClient)
		api.PUT("/clients/:id", clientService.UpdateClient)

		api.GET("/objects/:id", objectService.GetObjects)
		api.POST("/objects", objectService.CreateObjects)
		api.DELETE("/objects/:id", objectService.DeleteObjects)
		api.PUT("/objects/:id", objectService.UpdateObject)

		api.GET("/contracts/:id", contractService.GetContract)
		api.POST("/contracts", contractService.CreateContract)
		api.DELETE("/contracts/:id", contractService.DeleteContract)
		api.PUT("/contracts/:id", contractService.UpdateContract)

	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
