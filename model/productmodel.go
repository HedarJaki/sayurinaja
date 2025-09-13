package model

import (
	"time"
)

type Product struct {
	ProductID    int     `gorm:"column:productID;primaryKey;autoIncrement"`
	Product_name string  `gorm:"not null"`
	Category     string  `gorm:"type:enum('daging', 'sayur','buah');not null"`
	StoreID      int     `gorm:"column:storeID;not null"`
	Stock        int     `gorm:"column:stock;not null"`
	Price_Each   float64 `gorm:"column:price_each;not null"`
	Description  string
	Created_at   time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at   time.Time `gorm:"column:updated_at" json:"updated_at"`
	Image_path   string
}

type Store struct {
	StoreID         int    `gorm:"column:storeID;primaryKey;autoIncrement"`
	UserID          int    `gorm:"column:userID;not null"`
	StoreName       string `gorm:"not null"`
	StoreDesription string
	Created_at      time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at      time.Time `gorm:"column:updated_at" json:"updated_at"`
	StoreAddress    string
	LogoPath        string
}
