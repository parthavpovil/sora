package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parthavpovil/sora/models"
	"github.com/parthavpovil/sora/database"
	"github.com/gorilla/websocket"

)

type Client struct{


}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {return true},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast =make(chan []byte)

func HandleConnections(c *gin.Context){
	ws,err :=upgrader.Upgrade(c.Writer,c.Request,nil)

	if err !=nil{
		return
	}
	defer ws.Close()
	clients[ws]=true

	for{
		_,msg,err :=ws.ReadMessage()
		if err !=nil{
			delete(clients,ws)
			break
		}
		broadcast <-msg
	}
}

func HandleMessage(){
	for{
		msg := <-broadcast
		for client :=range clients{
			err :=client.WriteMessage(websocket.TextMessage,msg)
			if err !=nil{
				client.Close()
				delete(clients,client)
			}
		}
	}
}

func SentMsg() gin.HandlerFunc{
	return  func(c *gin.Context){

		var msg models.Message
		err :=c.BindJSON(&msg)
		if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}
		userIdInterface,exit :=c.Get("userId")
		if !exit{
			   c.JSON(http.StatusUnauthorized, gin.H{
                "error": "user not authenticated",
            })
            return
        }
		
		userId, _ := userIdInterface.(float64)
		userIdInt := int(userId) 
		msg.User_id =&userIdInt
		
		now := time.Now()
		msg.Created_at = &now

		err =database.DB.Create(&msg).Error
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error": "failed to save message",
			})
			return 
		}
		c.JSON(http.StatusOK, gin.H{
            "message": "message sent successfully",
            "data": msg,
        })



	}
}