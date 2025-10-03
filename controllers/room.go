package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parthavpovil/sora/database"
	"github.com/parthavpovil/sora/models"
	"gorm.io/gorm"
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

func JoinRoom() gin.HandlerFunc{
	return func(c *gin.Context) {
		var room models.Room
		roomId,_:=strconv.Atoi(c.Param("roomid"))
		

		userinter,_ :=c.Get("userId")
		userID:=int(userinter.(float64))

		result :=database.DB.Where("id=?",roomId).First(&room)
		if result.Error !=nil{
			if errors.Is(result.Error,gorm.ErrRecordNotFound){
				c.JSON(http.StatusNotFound,gin.H{
					"error":"room not found",
				})
				return 
			}
			  c.JSON(http.StatusInternalServerError,gin.H{
        		"error":"database error",
				})
				return
		}
		if !room.Is_private{
			 var existingMember models.RoomMembers
				checkResult := database.DB.Where("room_id = ? AND user_id = ?", roomId, userID).First(&existingMember)
				if checkResult.Error == nil {
					c.JSON(http.StatusConflict, gin.H{
						"error": "user already in room",
					})
					return
				}
			var roomMember models.RoomMembers
			roomMember.Room_id=&roomId
			roomMember.User_id=&userID
			roomMember.Joined_at=time.Now()


			err :=database.DB.Create(&roomMember).Error
			if err!=nil{
				c.JSON(http.StatusInternalServerError,gin.H{
					"message":"failed added to room ",
				})
				return 
			}
			c.JSON(http.StatusOK,gin.H{
				"message":"user added to room",
				"details":roomMember,
			})
			
		}else{
			c.JSON(http.StatusForbidden,gin.H{
				"error":"cannot join room its private",
			
			})
		}
	}
}

func LeaveRoom() gin.HandlerFunc{
	return func(c *gin.Context) {
		var room models.Room
		roomId:=c.Param("roomid")
		userIdInterface,_:=c.Get("userId")
		userId :=int(userIdInterface.(float64))

		result :=database.DB.Where("id=?",roomId).First(&room)
		if result.Error !=nil{
			if errors.Is(result.Error,gorm.ErrRecordNotFound){
				c.JSON(http.StatusNotFound,gin.H{
					"error":"room not found",
				})
				return 
			}
			  c.JSON(http.StatusInternalServerError,gin.H{
        		"error":"database error",
				})
				return
		}
		var existingMember models.RoomMembers
				checkResult := database.DB.Where("room_id = ? AND user_id = ?", roomId, userId).First(&existingMember)
				if checkResult.Error !=nil{
			if errors.Is(result.Error,gorm.ErrRecordNotFound){
				c.JSON(http.StatusNotFound,gin.H{
					"error":"user not in room not found",
				})
				return 
			}
			  c.JSON(http.StatusInternalServerError,gin.H{
        		"error":"database error",
				})
				return
		}

		err:=database.DB.Where("room_id=? AND user_id=?",roomId,userId ).Delete(&models.RoomMembers{}).Error
		if err != nil {
            c.JSON(http.StatusInternalServerError,gin.H{
                "error":"failed to leave room",
            })
            return
        }

        c.JSON(http.StatusOK,gin.H{
            "message":"successfully left room",
        })
	}
}