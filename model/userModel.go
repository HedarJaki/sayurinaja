package model

import (
	"time"
)

type User struct {
	UserID    uint      `gorm:"column:userID;primaryKey;autoIncrement"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Is_seller bool      `gorm:"default:false" json:"is_seller"`
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

type Order struct {
	OrderID    int       `gorm:"column:orderID;primaryKey;autoIncrement"`
	UserId     int       `gorm:"column:userID;not null"`
	Order_date time.Time `gorm:"column:order_at" json:"created_at"`
	Updated_at time.Time `gorm:"column:updated_at" json:"updated_at"`
	Status     string    `gorm:"type:enum('dalam proses','sedang dikirim','telah sampai');defaul:'pending';not null"`
}

type OrderDetail struct {
	OrderDetailID int `gorm:"column:orderdetailID;primaryKey;autoIncrement"`
	ProductID     int `gorm:"column:productID;not null"`
	OrderID       int `gorm:"column:orderIDID;not null"`
	Quantity      int `gorm:"column:quantity;not null"`
	Price         int `gorm:"column:price;not null"`
}

type Cart struct {
	CartID     int       `gorm:"column:cartID;primaryKey;autoIncrement"`
	UserID     int       `gorm:"column:userID;not null"`
	Is_active  bool      `gorm:"column:is_activer;Default:'true'"`
	Created_at time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type CartItem struct {
	CartItemID int       `gorm:"column:cartItemID;primaryKey;autoIncrement"`
	CartID     int       `gorm:"column:cartID"`
	ProductID  int       `gorm:"column:productID;not null"`
	Quantity   int       `gorm:"column:quantity;not null"`
	Price      int       `gorm:"column:price;not null"`
	Created_at time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Address struct {
	AddressID int `gorm:"column:addressID;primaryKey;autoIncrement"`
}
