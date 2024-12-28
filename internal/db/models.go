package db

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
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
	CreatedAt              time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
}

type GormStringList []string

// Scan implements the Scanner interface for GormStringList
func (list *GormStringList) Scan(value interface{}) error {
	if value == nil {
		*list = []string{}
		return nil
	}

	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("cannot scan type %T into GormStringList", value)
	}

	// Parse PostgreSQL array format: {"value1","value2"}
	str = strings.Trim(str, "{}") // Remove the curly braces
	if str == "" {
		*list = []string{}
	} else {
		*list = strings.Split(str, ",")
		for i, s := range *list {
			(*list)[i] = strings.Trim(s, `"`) // Remove quotes around each string
		}
	}

	return nil
}

// Value implements the Valuer interface for GormStringList
func (list GormStringList) Value() (driver.Value, error) {
	if list == nil {
		return "{}", nil
	}
	// Convert Go slice to PostgreSQL array format
	return "{" + strings.Join(list, ",") + "}", nil
}


func Migrate() {
	DB.AutoMigrate(&User{}, &Product{})
}
