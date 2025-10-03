package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthavpovil/sora/database"
	"github.com/parthavpovil/sora/models"
)

func CreateRoom() gin.HandlerFunc{

	
	return  func(c *gin.Context) {
		var room models.Room
		err :=c.BindJSON(&room)
		if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}
		userIdInterface, exists := c.Get("userId")
		 if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "user not authenticated",
            })
            return
        }

		userId :=int(userIdInterface.(float64))
		room.CreatedBy=&userId

		err =database.DB.Create(&room).Error
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error creating room",
				
			})
			return 
		}
		c.JSON(http.StatusOK,gin.H{
			"message":"room created",
			"details":room,
		})

	}
}