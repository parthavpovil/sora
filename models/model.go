package models

import "time"
 
type User struct{
	Id *int 	`gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"unique,not null" json:"user_name"`
	Password_hash string `gorm:"not null" json:"password_hash"`
}

type Message struct{
	Id int	`gorm:"primaryKey;auto increment" json:"id"`
	User_id string `gorm:"not null" json:"user_id"`
	Content string `gorm:"not null" json:"content"`
	Created_at time.Time `gorm:"not null" json:"created_at"`
}