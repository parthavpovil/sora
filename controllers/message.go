package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/parthavpovil/sora/database"
	"github.com/parthavpovil/sora/middleware"
	"github.com/parthavpovil/sora/models"
	"gorm.io/gorm"
)

type Client struct {
	Conn     *websocket.Conn
	UserID   int
	Username string
}
type IncomingMessage struct {
    RoomID  int    `json:"room_id"`
    Content string `json:"content"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]*Client)
var broadcast = make(chan models.Message)

func HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		return
	}
	defer ws.Close()

	tokenString := c.Query("token")
	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		ws.WriteMessage(websocket.TextMessage, []byte("missing token"))
		ws.Close()
		return
	}
	claims, err := middleware.ParseToken(tokenString)

	if err != nil {
		fmt.Printf("Token parse error: %v\n", err)
		fmt.Printf("Token string: %s\n", tokenString)
		ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		ws.Close()
		return
	}

	client := &Client{
		Conn:     ws,
		UserID:   int(claims["userId"].(float64)),
		Username: claims["username"].(string),
	}
	clients[ws] = client

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}

		var incoming IncomingMessage
		err=json.Unmarshal(msg,&incoming)
		if err!=nil{
			ws.WriteMessage(websocket.TextMessage, []byte("invalid message format cant sent"))
        continue
		}
		err =database.DB.Where("user_id=? AND room_id=?",client.UserID,incoming.RoomID).First(&models.RoomMembers{}).Error

		if errors.Is(err,gorm.ErrRecordNotFound){
			ws.WriteMessage(websocket.TextMessage,[]byte("user not in room"))
			continue
		}

		message := &models.Message{
			User_id:    client.UserID,
			Username:   client.Username,
			Room_id: incoming.RoomID,
			Content:    incoming.Content,
			Created_at: time.Now(),
		}
		
		err = database.DB.Create(&message).Error
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("failed to add to db"))
			continue
		}

		broadcast <- *message
	}
}

func HandleMessage() {
	for {
		msg := <-broadcast
		var members []int
		database.DB.Model(&models.RoomMembers{}).Where("room_id=?",msg.Room_id).Pluck("user_id",&members)

		memberMap := make(map[int]bool)
		for _,userID:= range members{
			memberMap[userID]=true
		} 
	
		for conn,client := range clients {
			if memberMap[client.UserID]{
				jsonMsg, _ := json.Marshal(msg)
				err := conn.WriteMessage(websocket.TextMessage, jsonMsg)
				if err != nil {
					conn.Close()
					delete(clients, conn)
				}
		}
	}
	}
}
