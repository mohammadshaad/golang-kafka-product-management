package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"net/http/httptest"
	"os"
	"testing"
	"log"

	"github.com/stretchr/testify/assert"
	"github.com/mohammadshaad/zocket/tests/testutils"
	"github.com/mohammadshaad/zocket/internal/db"
	"github.com/mohammadshaad/zocket/internal/cache"
	"github.com/mohammadshaad/zocket/internal/queue"
	"github.com/mohammadshaad/zocket/config"
	"github.com/joho/godotenv"
)

func setup() {
	loadEnv()
	config.LoadConfig()
	db.InitDatabase()
	initSchema()
	initCache()
	initQueue()
}

func loadEnv() {
	err := godotenv.Load("/Users/mohammadshaad/Desktop/zocket/.env")
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}
	printEnvVariables()
}

func printEnvVariables() {
	log.Printf("KAFKA_BROKERS: %s", os.Getenv("KAFKA_BROKERS"))
	log.Printf("KAFKA_TOPIC: %s", os.Getenv("KAFKA_TOPIC"))
	log.Printf("REDIS_ADDR: %s", os.Getenv("REDIS_ADDR"))
	log.Printf("REDIS_PASSWORD: %s", os.Getenv("REDIS_PASSWORD"))
	log.Printf("REDIS_USERNAME: %s", os.Getenv("REDIS_USERNAME"))
}

func initSchema() {
	// Create the products table
	db.DB.Exec(`CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		user_id INT,
		product_name TEXT,
		product_description TEXT,
		product_images TEXT[],
		compressed_product_images TEXT[],
		product_price FLOAT,
		created_at TIMESTAMP
	)`)
}

func initCache() {
	REDIS_ADDR := os.Getenv("REDIS_ADDR")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	USERNAME := os.Getenv("REDIS_USERNAME")
	cache.InitRedis(REDIS_ADDR, USERNAME, REDIS_PASSWORD)
}

func initQueue() {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		log.Fatal("KAFKA_TOPIC environment variable is not set")
	}
	queue.InitProducerWithTopic([]string{brokers}, topic)
}

func TestAddProduct(t *testing.T) {
	setup()
	router := testutils.SetupTestRouter()
	
	jsonData, _ := json.Marshal(testutils.TestProduct)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	product, ok := response["product"].(map[string]interface{})
	assert.True(t, ok, "Product should be present in the response")
	log.Printf("Product in response: %+v", product)
	for key := range product {
		log.Printf("Key: %s, Value: %v", key, product[key])
	}
	assert.Equal(t, testutils.TestProduct.ProductName, product["ProductName"].(string))
	assert.Equal(t, testutils.TestProduct.ProductPrice, product["ProductPrice"].(float64))
}

func TestGetProductByID(t *testing.T) {
	setup()
	router := testutils.SetupTestRouter()
	
	var createdProduct db.Product
	db.DB.Create(&testutils.TestProduct).Scan(&createdProduct)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products/"+strconv.Itoa(int(createdProduct.ID)), nil)
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response db.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, testutils.TestProduct.ProductName, response.ProductName)
}