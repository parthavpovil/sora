package models

import "time"
 
type User struct{
    Id int 	`gorm:"primaryKey;autoIncrement" json:"id"`
    Username string `gorm:"unique;not null" json:"user_name"`
    Password_hash string `gorm:"not null" json:"password_hash"`
}

type Message struct{
    Id int	`gorm:"primaryKey;autoIncrement" json:"id"`
    User_id int `gorm:"not null" json:"user_id"`
    User User `gorm:"foreignKey:User_id" json:"user,omitempty"`  
    Username string `gorm:"not null" json:"user_name"`
    Room_id int `gorm:"not null" json:"room_id"`
    Room Room `gorm:"foreignKey:Room_id" json:"room,omitempty"`  
    Content string `gorm:"not null" json:"content"`
    Created_at time.Time `gorm:"not null" json:"created_at"`
}

type Room struct{
    Id int `gorm:"primaryKey;autoIncrement" json:"id"`
    Name string `gorm:"not null;unique" json:"name"`
    Is_private bool `gorm:"not null" json:"is_private"`
    CreatedBy *int `gorm:"not null" json:"created_by"`
    Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

type RoomMembers struct{
    Id int `gorm:"primaryKey;autoIncrement" json:"id"`
    Room_id *int `gorm:"not null" json:"room_id"`
    Room Room `gorm:"foreignKey:Room_id" json:"room,omitempty"`
    User_id *int `gorm:"not null" json:"user_id"`
    User User `gorm:"foreignKey:User_id" json:"user,omitempty"`
    Joined_at time.Time `gorm:"not null" json:"joined_at"`
}