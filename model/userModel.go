package model

import (
	"time"
)

type User struct {
	UserID       int             `gorm:"column:userID;primaryKey;autoIncrement"`
	Username     string          `gorm:"unique;not null" json:"username"`
	Email        string          `gorm:"unique;not null" json:"email"`
	Password     string          `gorm:"not null" json:"password"`
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at" json:"updated_at"`
	Is_seller    bool            `gorm:"default:false" json:"is_seller"`
	ProdReviews  []Productreview `gorm:"foreignKey:UserID;references:UserID"`
	StoreReviews []StoreReview   `gorm:"foreignKey:StoreID;references:UserID"`
}

type Order struct {
	OrderID      int           `gorm:"column:orderID;primaryKey;autoIncrement"`
	UserID       int           `gorm:"column:userID;not null"`
	StoreID      int           `gorm:"column:storeID;not null"`
	AddressID    int           `gorm:"column:addressID;not null"`
	Total_price  float64       `gorm:"column:total_price;not null"`
	Order_date   time.Time     `gorm:"column:order_at" json:"created_at"`
	Updated_at   time.Time     `gorm:"column:updated_at" json:"updated_at"`
	Status       string        `gorm:"type:enum('dalam proses','sedang dikirim','telah sampai');default:'dalam proses';not null"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID;references:OrderID"`
	//User         User          `gorm:"foreignKey:UserID;references:UserID"`
	//Store        Store         `gorm:"foreignKey:StoreID;references:StoreID"`
}

type OrderDetail struct {
	OrderDetailID int       `gorm:"column:orderdetailID;primaryKey;autoIncrement"`
	ProductID     int       `gorm:"column:productID;not null"`
	OrderID       int       `gorm:"column:orderID;not null"`
	Quantity      int       `gorm:"column:quantity;not null"`
	Price         float64   `gorm:"column:price;not null"`
	Note          string    `gorm:"column:note"`
	Created_at    time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Cart struct {
	CartID     int        `gorm:"column:cartID;primaryKey;autoIncrement"`
	UserID     int        `gorm:"column:userID;not null"`
	Is_active  bool       `gorm:"column:is_active;Default:'true'"`
	Created_at time.Time  `gorm:"column:created_at" json:"created_at"`
	Updated_at time.Time  `gorm:"column:updated_at" json:"updated_at"`
	CartItems  []CartItem `gorm:"foreignKey:CartID;references:CartID"`
}

type CartItem struct {
	CartItemID int       `gorm:"column:cartItemID;primaryKey;autoIncrement"`
	ProductID  int       `gorm:"column:productID;not null"`
	CartID     int       `gorm:"column:cartID"`
	Quantity   int       `gorm:"column:quantity;not null"`
	Price      float64   `gorm:"column:price;not null"`
	Note       string    `gorm:"column:note"`
	Created_at time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type StoreReview struct {
	StorereviewID     int    `gorm:"column:storereviewID;primaryKey;autoIncrement"`
	StoreID           int    `gorm:"column:storeID;not null"`
	UserID            int    `gorm:"column:userID;not null"`
	Star              int    `gorm:"column:star;not null"`
	Store_review_desc string `gorm:"column:store_review_desc;not null"`
}

type Address struct {
	AddressID   int     `gorm:"column:addressID;primaryKey;autoIncrement"`
	UserID      int     `gorm:"column:userID;not null"`
	Fulladdress string  `gorm:"column:fulladdress;not null"`
	Isdefault   bool    `gorm:"column:isDefault;not null"`
	Latitude    float64 `gorm:"column:latitude;not null"`
	Longitude   float64 `gorm:"column:longitude;not null"`
}
