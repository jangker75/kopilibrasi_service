package routes

import (
	"go-rest-api/controllers"
	"go-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
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

	return r
}
