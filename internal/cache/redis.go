package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/mohammadshaad/zocket/internal/db"
)

var rdb *redis.Client
var ctx = context.Background()

const (
    defaultTTL = 1 * time.Hour
    productKeyPrefix = "product:"
)

func InitRedis(addr, username, password string) {
    rdb = redis.NewClient(&redis.Options{
        Addr:     addr,
        Username: username,
        Password: password,
        DB:       0,
    })

    // Test the connection
    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Println("Successfully connected to Redis")
}

// GetProductFromCache retrieves a product from cache
func GetProductFromCache(productID string) (*db.Product, error) {
    key := fmt.Sprintf("%s%s", productKeyPrefix, productID)
    data, err := rdb.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return nil, nil // Cache miss
        }
        return nil, err
    }

    var product db.Product
    if err := json.Unmarshal([]byte(data), &product); err != nil {
        return nil, err
    }

    return &product, nil
}

// SetProductInCache stores a product in cache
func SetProductInCache(product *db.Product) error {
    key := fmt.Sprintf("%s%d", productKeyPrefix, product.ID)
    data, err := json.Marshal(product)
    if err != nil {
        return err
    }

    return rdb.Set(ctx, key, data, defaultTTL).Err()
}

// InvalidateProductCache removes a product from cache
func InvalidateProductCache(productID string) error {
    key := fmt.Sprintf("%s%s", productKeyPrefix, productID)
    return rdb.Del(ctx, key).Err()
}
