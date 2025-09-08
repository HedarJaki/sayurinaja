package model

import (
	"time"
)

type Product struct {
	ProductID    int    `gorm:"column:productID;primaryKey;autoIncrement"`
	Product_name string `gorm:"not null"`
	Category     string `gorm:"type:enum('daging', 'sayur','buah');not null"`
	StoreID      int    `gorm:"column:storeID;not null"`
	Stock        int    `gorm:"column:stock;not null"`
	PriceEach    int    `gorm:"column:price_each;not null"`
	Description  string
	Created_at   time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at   time.Time `gorm:"column:updated_at" json:"updated_at"`
	Image_path   string
}
