package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Product_name string
	StoreID      int
	Stock        int
	Price        int
	Description  string
	Created_at   time.Time
	Updated_at   time.Time
	Image_path   string
}
