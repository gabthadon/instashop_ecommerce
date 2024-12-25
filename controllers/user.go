package controllers


import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"instashop_ecommerce/util"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) { 
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind the JSON request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Fetch the user and role from the database
	var hashedPassword, role string
	err := util.DB.QueryRow("SELECT password, role FROM users WHERE username = ?", req.Username).Scan(&hashedPassword, &role)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := util.GenerateJWT(req.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token
	c.JSON(http.StatusOK, gin.H{"token": token})
}


// registerHandler handles user registration
func RegisterHandler(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
        Role     string `json:"role"` // Optional role field, defaults to "customer"
    }

    // Bind JSON input
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // Validate role
    if req.Role != "" && req.Role != "admin" && req.Role != "customer" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
        return
    }

    if req.Role == "" {
        req.Role = "customer" // Default to customer role
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    // Insert the user into the database
    _, err = util.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", req.Username, string(hashedPassword), req.Role)
    if err != nil {
        if sql.ErrNoRows == err {
            c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
        return
    }

    // Success response
    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
