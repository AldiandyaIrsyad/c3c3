package main

import (
	"github.com/gin-gonic/gin"

	"github.com/aldiandyaIrsyad/c3c2/handlers"
)

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		authorized := api.Group("/", handlers.JWTAuthMiddleware())
		{
			authorized.POST("/products", handlers.CreateProduct)
			authorized.GET("/products", handlers.GetProducts)
			authorized.GET("/products/:id", handlers.GetProductByID)
			authorized.PUT("/products/:id", handlers.UpdateProduct)
			authorized.DELETE("/products/:id", handlers.DeleteProduct)
		}
	}

	router.Run(":8080")
}
