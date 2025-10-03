package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parthavpovil/sora/controllers"
	"github.com/parthavpovil/sora/middleware"
)

func RoomRoutes(incomingRoutes *gin.Engine) {
	protected := incomingRoutes.Group("/")
	protected.Use(middleware.JWTverify())
	{
		protected.POST("rooms/create", controllers.CreateRoom())
		protected.POST("/rooms/:roomid/join",controllers.JoinRoom())
		protected.POST("/rooms/:roomid/leave",controllers.LeaveRoom())
	}
}
