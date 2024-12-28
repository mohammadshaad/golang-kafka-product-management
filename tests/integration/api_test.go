package integration

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "strconv"

    "encoding/json"
    "github.com/stretchr/testify/assert"
    "github.com/mohammadshaad/zocket/internal/db"
    "github.com/mohammadshaad/zocket/tests/testutils"
    "github.com/mohammadshaad/zocket/config"
    "github.com/mohammadshaad/zocket/internal/cache"
    "github.com/joho/godotenv"
    "log"
    "os"
)

func setup() {
    loadEnv()
    config.LoadConfig()
    db.InitDatabase()
    initSchema()
    initCache()
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
    log.Printf("DATABASE_DSN: %s", os.Getenv("DATABASE_DSN"))
    log.Printf("KAFKA_BROKERS: %s", os.Getenv("KAFKA_BROKERS"))
    log.Printf("KAFKA_TOPIC: %s", os.Getenv("KAFKA_TOPIC"))
    log.Printf("REDIS_ADDR: %s", os.Getenv("REDIS_ADDR"))
    log.Printf("REDIS_PASSWORD: %s", os.Getenv("REDIS_PASSWORD"))
    log.Printf("REDIS_USERNAME: %s", os.Getenv("REDIS_USERNAME"))
    log.Printf("AWS_REGION: %s", os.Getenv("AWS_REGION"))
    log.Printf("S3_BUCKET: %s", os.Getenv("S3_BUCKET"))
    log.Printf("AWS_ACCESS_KEY_ID: %s", os.Getenv("AWS_ACCESS_KEY_ID"))
    log.Printf("AWS_SECRET_ACCESS_KEY: %s", os.Getenv("AWS_SECRET_ACCESS_KEY"))
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

func TestGetProductByIDWithCache(t *testing.T) {
    setup()
    router := testutils.SetupTestRouter()
    
    var product db.Product
    db.DB.Create(&testutils.TestProduct).Scan(&product)
    
    // First request (cache miss)
    start := time.Now()
    w1 := httptest.NewRecorder()
    req1, _ := http.NewRequest("GET", "/api/v1/products/"+strconv.FormatUint(uint64(product.ID), 10), nil)
    router.ServeHTTP(w1, req1)
    cacheMissTime := time.Since(start)
    
    // Second request (cache hit)
    start = time.Now()
    w2 := httptest.NewRecorder()
    req2, _ := http.NewRequest("GET", "/api/v1/products/"+strconv.FormatUint(uint64(product.ID), 10), nil)
    router.ServeHTTP(w2, req2)
    cacheHitTime := time.Since(start)
    
    assert.Less(t, cacheHitTime, cacheMissTime)
}

func TestGetAllProducts(t *testing.T) {
    setup()
    router := testutils.SetupTestRouter()
    
    products := []db.Product{
        testutils.TestProduct,
        {
            UserID:      1,
            ProductName: "Test Product 2",
            ProductPrice: 199.99,
        },
    }
    
    for _, p := range products {
        db.DB.Create(&p)
    }
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/products", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response []db.Product
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.GreaterOrEqual(t, len(response), 2)
}