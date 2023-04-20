package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/aldiandyaIrsyad/c3c2/database"
	"github.com/aldiandyaIrsyad/c3c2/jwt"
	"github.com/aldiandyaIrsyad/c3c2/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/register", Register)
		api.POST("/login", Login)

		authorized := api.Group("/", JWTAuthMiddleware())
		{
			authorized.GET("/products/:id", GetProductByID)
			authorized.GET("/products", GetProducts)
		}
	}

	return router
}

func createTestUser() models.User {
	db, _ := database.Connect()
	defer db.Close()

	user := models.User{
		Username: "testuser",
		Password: "testpassword",
		Role:     models.Creator,
	}

	db.Create(&user)

	return user
}

func createTestProduct(user models.User) models.Product {
	db, _ := database.Connect()
	defer db.Close()

	product := models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
		UserID:      user.ID,
	}

	db.Create(&product)

	return product
}

func generateTestToken(user models.User) string {
	token, _ := jwt.GenerateToken(user.Username, user.Role)
	return token
}

func TestGetProductByID(t *testing.T) {
	router := setupTestRouter()

	// Prepare test data
	user := createTestUser()
	product := createTestProduct(user)
	token := generateTestToken(user)

	// Test get product by ID (found)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/products/"+strconv.Itoa(int(product.ID)), nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), product.Name)

	// Test get product by ID (not found)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/999999", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllProducts(t *testing.T) {
	router := setupTestRouter()

	// Prepare test data
	user := createTestUser()
	createTestProduct(user)
	token := generateTestToken(user)

	// Test get all products
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/products", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Product")
}
