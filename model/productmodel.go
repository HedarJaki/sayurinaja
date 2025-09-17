package model

import (
	"time"
)

type Product struct {
	ProductID           int       `gorm:"column:productID;primaryKey;autoIncrement"`
	Product_name        string    `gorm:"not null"`
	Category            string    `gorm:"type:enum('daging', 'sayur','buah');not null"`
	StoreID             int       `gorm:"column:storeID;not null"`
	Stock               int       `gorm:"column:stock;not null"`
	Price_Each          float64   `gorm:"column:price_each;not null"`
	Product_description string    `gorm:"column:product_description"`
	Created_at          time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at          time.Time `gorm:"column:updated_at" json:"updated_at"`
	Image_path          string
	Reviews             []Productreview `gorm:"foreignKey:ProductID;references:ProductID"`
	//Store               Store           `gorm:"foreignKey:StoreID;references:StoreID"`
}

type Store struct {
	StoreID         int     `gorm:"column:storeID;primaryKey;autoIncrement"`
	UserID          int     `gorm:"column:userID;not null"`
	StoreName       string  `gorm:"not null"`
	Rating          float64 `gorm:"column:rating"`
	StoreDesription string
	Created_at      time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at      time.Time `gorm:"column:updated_at" json:"updated_at"`
	StoreAddress    string
	LogoPath        string
	Products        []Product     `gorm:"foreignKey:StoreID;references:StoreID"`
	Reviews         []StoreReview `gorm:"foreignKey:StoreID;references:StoreID"`
	Orders          []Order       `gorm:"foreignKey:StoreID;references:StoreID"`
	//User            User          `gorm:"foreignKey:UserID;references:UserID"`
}

type Productreview struct {
	ReviewID    int    `gorm:"column:reviewID;primaryKey;autoIncrement"`
	ProductID   int    `gorm:"column:productID;not null"`
	UserID      int    `gorm:"column:userID;not null"`
	Star        int    `gorm:"column:star;not null"`
	Description string `gorm:"column:description;not null"`
}
