package controllers

import (
	"net/http"               // For HTTP status codes
	"github.com/gin-gonic/gin" // Gin framework
	"instashop_ecommerce/util" // Database utility
)


type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}


func GetAllProductsHandler(c *gin.Context) {
    rows, err := util.DB.Query("SELECT id, name, description, price, quantity FROM products")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
        return
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var product Product
        if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse product data"})
            return
        }
        products = append(products, product)
    }

    c.JSON(http.StatusOK, products)
}


//Update Products
func UpdateProductHandler(c *gin.Context) {
    var product Product

    // Bind JSON input
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
        return
    }

    // Update database
    _, err := util.DB.Exec("UPDATE products SET name = ?, description = ?, price = ?, quantity = ? WHERE id = ?",
        product.Name, product.Description, product.Price, product.Quantity, product.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}



//Delete Products
func DeleteProductHandler(c *gin.Context) {
    id := c.Param("id")

    // Delete from database
    _, err := util.DB.Exec("DELETE FROM products WHERE id = ?", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}


func CreateProductHandler(c *gin.Context) {
    var product Product

    // Bind JSON input
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
        return
    }

    // Insert into database
    result, err := util.DB.Exec("INSERT INTO products (name, description, price, quantity) VALUES (?, ?, ?, ?)",
        product.Name, product.Description, product.Price, product.Quantity)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
        return
    }

    id, _ := result.LastInsertId()
    product.ID = int(id)
    c.JSON(http.StatusCreated, product)
}