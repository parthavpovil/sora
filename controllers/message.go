package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/parthavpovil/sora/database"
	"github.com/parthavpovil/sora/middleware"
	"github.com/parthavpovil/sora/models"
)

type Client struct {
	Conn     *websocket.Conn
	UserID   int
	Username string
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

		message := &models.Message{
			User_id:    client.UserID,
			Username:   client.Username,
			Content:    string(msg),
			Created_at: time.Now(),
		}
		err = database.DB.Create(&message).Error
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("failed to add to db"))
		}

		broadcast <- *message
	}
}

func HandleMessage() {
	for {
		msg := <-broadcast
		for client := range clients {
			jsonMsg, _ := json.Marshal(msg)
			err := client.WriteMessage(websocket.TextMessage, jsonMsg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
