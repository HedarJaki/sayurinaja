package model

import (
	"time"

	"gorm.io/gorm"
)

type User_app struct {
	gorm.Model
	Username   string    `gorm:"primaryKey;unique" json:"username"`
	Email      string    `gorm:"unique;not null" json:"email"`
	Password   string    `gorm:"not null" json:"Password"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Is_seller  bool      `gorm:"default:false" json:"is_seller"`
}

func (User_app) Tablename() string {
	return "user_app"
}
