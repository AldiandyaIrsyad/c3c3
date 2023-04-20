package handlers

import (
	"net/http"

	"github.com/aldiandyaIrsyad/c3c2/database"
	"github.com/aldiandyaIrsyad/c3c2/jwt"
	"github.com/aldiandyaIrsyad/c3c2/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Register(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	newUser.Role = models.Creator

	db, _ := database.Connect()
	defer db.Close()

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var credentials models.User
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	db, _ := database.Connect()
	defer db.Close()

	var user models.User
	if err := db.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user"})
		}
		return
	}

	if user.Password != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	token, err := jwt.GenerateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Logged in successfully",
	})
}

func CreateProduct(c *gin.Context) {
	user, _ := c.Get("user")

	var newProduct models.Product
	if err := c.BindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	newProduct.UserID = user.(*models.User).ID

	db, _ := database.Connect()
	defer db.Close()

	if err := db.Create(&newProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func GetProducts(c *gin.Context) {
	user, _ := c.Get("user")

	db, _ := database.Connect()
	defer db.Close()

	var products []models.Product
	if user.(*models.User).Role == models.Admin {
		db.Find(&products)
	} else {
		db.Where("user_id = ?", user.(*models.User).ID).Find(&products)
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func UpdateProduct(c *gin.Context) {
	user, _ := c.Get("user")
	productID := c.Param("id")

	db, _ := database.Connect()
	defer db.Close()

	var product models.Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query product"})
		}
		return
	}

	if user.(*models.User).Role != models.Admin && product.UserID != user.(*models.User).ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this product"})
		return
	}

	var updatedProduct models.Product
	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	db.Model(&product).Updates(updatedProduct)
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func DeleteProduct(c *gin.Context) {
	user, _ := c.Get("user")
	productID := c.Param("id")

	db, _ := database.Connect()
	defer db.Close()

	var product models.Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query product"})
		}
		return
	}

	if user.(*models.User).Role != models.Admin && product.UserID != user.(*models.User).ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this product"})
		return
	}

	db.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		user, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func GetProductByID(c *gin.Context) {
	user, _ := c.Get("user")
	productID := c.Param("id")

	db, _ := database.Connect()
	defer db.Close()

	var product models.Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query product"})
		}
		return
	}

	if user.(*models.User).Role != models.Admin && product.UserID != user.(*models.User).ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view this product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}
