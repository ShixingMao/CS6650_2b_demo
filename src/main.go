package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Product structure — matches spec requirements
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`     // searchable
	Category    string `json:"category"` // searchable
	Description string `json:"description"`
	Brand       string `json:"brand"`
}

// Thread-safe in-memory store using sync.Map
var (
	productStore  sync.Map
	totalProducts = 100000
	productIDs    []int // fixed list of all IDs to allow deterministic iteration
)

// GenerateProducts creates 100,000 products with rotating categories and brands
func generateProducts() {
	brands := []string{"Alpha", "Beta", "Gamma", "Delta", "Omega", "Nova", "Apex", "Orion", "Epsilon", "Zeta"}
	categories := []string{"Electronics", "Books", "Home", "Toys", "Sports", "Clothing", "Beauty", "Garden", "Office", "Grocery"}

	for i := 1; i <= totalProducts; i++ {
		brand := brands[i%len(brands)]
		category := categories[i%len(categories)]
		name := fmt.Sprintf("Product %s %d", brand, i)

		p := Product{
			ID:          i,
			Name:        name,
			Category:    category,
			Description: "A reliable product for everyday use.",
			Brand:       brand,
		}

		productStore.Store(i, p)
		productIDs = append(productIDs, i)
	}

	fmt.Printf("✅ Generated %d products successfully.\n", totalProducts)
}

func main() {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})
	// Generate 100k products at startup
	generateProducts()

	// Routes
	router.GET("/products/:id", getProduct)
	router.GET("/products", getSampleProducts) // quick test
	router.GET("/products/search", searchProducts)

	// Run the server
	router.Run("0.0.0.0:8080")
}

// GET /products/:id — retrieve product by ID
func getProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_INPUT",
			"message": "Product ID must be a positive integer",
		})
		return
	}

	v, ok := productStore.Load(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "NOT_FOUND",
			"message": "Product not found",
		})
		return
	}

	product := v.(Product)
	c.JSON(http.StatusOK, product)
}

// GET /products — quick check endpoint to view first few products
func getSampleProducts(c *gin.Context) {
	products := []Product{}
	count := 0

	productStore.Range(func(key, value interface{}) bool {
		if count >= 5 {
			return false
		}
		p := value.(Product)
		products = append(products, p)
		count++
		return true
	})

	c.JSON(http.StatusOK, gin.H{
		"total_generated": totalProducts,
		"sample_products": products,
	})
}

// GET /products/search?q={query}
// Each search checks exactly 100 products, counts every check, returns up to 20 matches
func searchProducts(c *gin.Context) {
	query := strings.TrimSpace(strings.ToLower(c.Query("q")))
	const checkLimit = 100
	const maxResults = 20

	// Validate that query param exists (you can relax this if "" is valid)
	if query == "" {
		// Allow empty searches if you want — comment/uncomment as needed:
		// c.JSON(http.StatusBadRequest, gin.H{"error": "MISSING_QUERY", "message": "Missing ?q parameter"})
		// return
	}

	start := time.Now()
	results := []Product{}
	totalFound := 0
	checked := 0

	// pick a stable start index (avoid out-of-range)
	startIndex := int(time.Now().UnixNano() % int64(totalProducts))

	for i := 0; i < checkLimit; i++ {
		index := (startIndex + i) % totalProducts
		id := productIDs[index]

		v, ok := productStore.Load(id)
		if !ok {
			continue
		}
		checked++

		p := v.(Product)
		if query == "" ||
			strings.Contains(strings.ToLower(p.Name), query) ||
			strings.Contains(strings.ToLower(p.Category), query) {
			totalFound++
			if len(results) < maxResults {
				results = append(results, p)
			}
		}
	}

	elapsed := time.Since(start)

	c.JSON(http.StatusOK, gin.H{
		"products":    results,
		"total_found": totalFound,
		"search_time": elapsed.String(),
	})
}
