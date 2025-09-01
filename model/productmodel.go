package model

import (
	"time"
)

type Product struct {
	Product_name string
	Category     string `gorm:"type:enum('daging', 'sayur','buah');not null"`
	StoreID      int
	Stock        int
	Price        int
	Description  string
	Created_at   time.Time
	Updated_at   time.Time
	Image_path   string
}
