package performance

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/mohammadshaad/zocket/internal/api"
    "github.com/mohammadshaad/zocket/internal/cache"
    "github.com/mohammadshaad/zocket/internal/db"
    "github.com/mohammadshaad/zocket/tests/testutils"
    "github.com/mohammadshaad/zocket/config"
    "github.com/joho/godotenv"
    "log"
    "os"
)

const (
    maxResponseTimeWithoutCache = 200 * time.Millisecond
    maxResponseTimeWithCache    = 50 * time.Millisecond
    minThroughput               = 1000 // requests per second
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

func BenchmarkGetProductByID(b *testing.B) {
    setup()
    router := gin.New()
    api.SetupRoutes(router)
    
    var product db.Product
    db.DB.Create(&testutils.TestProduct).Scan(&product)
    
    b.Run("Without Cache", func(b *testing.B) {
        cache.InvalidateProductCache(strconv.FormatUint(uint64(product.ID), 10))
        for i := 0; i < b.N; i++ {
            start := time.Now()
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("GET", "/api/v1/products/"+strconv.FormatUint(uint64(product.ID), 10), nil)
            router.ServeHTTP(w, req)
            duration := time.Since(start)
            if duration > maxResponseTimeWithoutCache {
                b.Errorf("Response time exceeded threshold: %v", duration)
            }
        }
    })
    
    b.Run("With Cache", func(b *testing.B) {
        // First request to populate cache
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/api/v1/products/"+strconv.FormatUint(uint64(product.ID), 10), nil)
        router.ServeHTTP(w, req)
        
        // Benchmark cached requests
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            start := time.Now()
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("GET", "/api/v1/products/"+strconv.FormatUint(uint64(product.ID), 10), nil)
            router.ServeHTTP(w, req)
            duration := time.Since(start)
            if duration > maxResponseTimeWithCache {
                b.Errorf("Response time exceeded threshold: %v", duration)
            }
        }
    })
}