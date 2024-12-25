package controllers

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"instashop_ecommerce/util"
)

type Order struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	ProductID  int    `json:"product_id"`
	Quantity   int    `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Status     string `json:"status"` // Pending, Shipped, Delivered, Canceled
}

func PlaceOrderHandler(c *gin.Context) {
	var order Order
	// Bind JSON input
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	// Fetch product price
	var price float64
	err := util.DB.QueryRow("SELECT price FROM products WHERE id = ?", order.ProductID).Scan(&price)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product price"})
		return
	}

	// Calculate total price
	order.TotalPrice = price * float64(order.Quantity)

	// Insert order into database
	result, err := util.DB.Exec("INSERT INTO orders (user_id, product_id, quantity, total_price, status) VALUES (?, ?, ?, ?, ?)",
		order.UserID, order.ProductID, order.Quantity, order.TotalPrice, "Pending")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
		return
	}

	// Get the order ID
	id, _ := result.LastInsertId()
	order.ID = int(id)

	// Respond with the order details
	c.JSON(http.StatusCreated, order)
}

func ListOrdersHandler(c *gin.Context) {
	rows, err := util.DB.Query("SELECT id, user_id, product_id, quantity, total_price, status FROM orders")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.TotalPrice, &order.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse order data"})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

func CancelOrderHandler(c *gin.Context) {
	id := c.Param("id")

	// Check if order exists and fetch current status
	var status string
	err := util.DB.QueryRow("SELECT status FROM orders WHERE id = ?", id).Scan(&status)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order status"})
		return
	}

	// Ensure the order is not already shipped or delivered
	if status == "Shipped" || status == "Delivered" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be canceled after shipping or delivery"})
		return
	}

	// Update the order status to 'Canceled'
	_, err = util.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", "Canceled", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order canceled successfully"})
}

func UpdateOrderStatusHandler(c *gin.Context) {
	var request struct {
		Status string `json:"status" binding:"required"`
	}

	// Bind the request body to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	id := c.Param("id")

	// Update the order status
	_, err := util.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", request.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
