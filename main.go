package main

import (
            // Import for HTTP status codes
	"github.com/gin-gonic/gin"
	"instashop_ecommerce/util"
	"instashop_ecommerce/controllers" // Import the controllers package
	"github.com/joho/godotenv"
	"log"
)



func main() {
	// Initialize the database connection
	util.ConnectDatabase()

	r := gin.Default()

	// Public routes
	r.POST("/login", controllers.LoginHandler)

	r.POST("/register", controllers.RegisterHandler)

	// Product routes
	//Here the admin middleware is use to restrict acces to only admin
    r.POST("/products", util.AdminMiddleware(), controllers.CreateProductHandler)       // Create
    r.GET("/products", util.AdminMiddleware(), controllers.GetAllProductsHandler )       // Read all
    r.PUT("/products/:id", util.AdminMiddleware(), controllers.UpdateProductHandler)    // Update
    r.DELETE("/products/:id", util.AdminMiddleware(), controllers.DeleteProductHandler) // Delete


	// Order routes
	r.POST("/orders", controllers.PlaceOrderHandler)     // Place order
	   // List all orders
	   r.GET("/orders", util.JWTMiddleware(), controllers.ListOrdersHandler) 
	r.PUT("/orders/:id/status", util.AdminMiddleware(), controllers.UpdateOrderStatusHandler) // Update order status
	r.DELETE("/orders/:id", controllers.CancelOrderHandler) // Cancel order


	// Protected routes
	auth := r.Group("/auth")
	auth.Use(util.JWTMiddleware()) // Use the JWT middleware
	{
	 
		auth.GET("/profile", profileHandler)
	}

	// Start the server
	r.Run(":8080")
}



// profileHandler handles user profile
func profileHandler(c *gin.Context) {
	// Example response for profile
	c.JSON(200, gin.H{"message": "Welcome to your profile!"})
}

func init() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file")
    }
}













