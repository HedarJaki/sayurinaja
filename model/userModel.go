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
	StoreID         int `gorm:"column:storeID;primaryKey;autoIncrement"`
	UserID          int `gorm:"column:userID;not null"`
	StoreName       string
	StoreDesription string
	Created_at      time.Time `gorm:"column:created_at" json:"created_at"`
	Updated_at      time.Time `gorm:"column:updated_at" json:"updated_at"`
	StoreAddress    string
	LogoPath        string
}
