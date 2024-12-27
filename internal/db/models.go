package db

import (
	"encoding/json"
	"database/sql/driver"
)

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Email string `gorm:"size:100;unique"`
}

type Product struct {
	ID                    uint           `gorm:"primaryKey"`
	UserID                uint           `gorm:"index"`
	ProductName           string         `gorm:"size:255"`
	ProductDescription    string         `gorm:"type:text"`
	ProductImages         GormStringList `gorm:"type:text[]"`
	CompressedProductImages GormStringList `gorm:"type:text[]"`
	ProductPrice          float64        `gorm:"type:decimal(10,2)"`
	CreatedAt             string
}

type GormStringList []string

func (list *GormStringList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &list)
}

func (list GormStringList) Value() (driver.Value, error) {
	return json.Marshal(list)
}

// AutoMigrate function
func Migrate() {
	DB.AutoMigrate(&User{}, &Product{})
}
