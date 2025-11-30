package routes

import (
	"fmt"
	"go-rest-api/config"
	"go-rest-api/controllers"
	"go-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if config.AppConfig.ENVMode == "production" {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("Running in production mode")
		fmt.Printf("Listening and serving HTTP on Port :%s\n", config.AppConfig.APPPort)
	}

	r := gin.Default()

	// Register the IPLogger middleware
	r.Use(middlewares.IPLogger())

	//Products Routes
	r.GET("/products", controllers.FindProducts)
	r.POST("/products", controllers.CreateProduct)
	r.GET("/products/:id", controllers.FindProduct)
	r.PUT("/products/:id", controllers.UpdateProduct)
	r.DELETE("/products/:id", controllers.DeleteProduct)

	// Menu Category Routes
	r.GET("/menu_categories", controllers.FindMenuCategories)
	r.POST("/menu_category", controllers.CreateMenuCategory)
	r.GET("/menu_category/:id", controllers.FindMenuCategory)
	r.PUT("/menu_category/:id", controllers.UpdateMenuCategory)
	r.DELETE("/menu_category/:id", controllers.DeleteMenuCategory)
	r.GET("/get_menus", controllers.FindMenuCategoriesWithMenuDetails)

	// Menu Detail Routes
	r.GET("/menu_details", controllers.FindMenuDetails)

	// Seed database route
	r.GET("/seed_dummy", controllers.SeedDummyData)

	// Transaction Routes
	r.POST("/sync-transaction", controllers.SyncTransactions)
	r.GET("/transactions", controllers.ListTransactions)
	return r
}
