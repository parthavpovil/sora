package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthavpovil/sora/controllers"
	"github.com/parthavpovil/sora/database"
	"github.com/parthavpovil/sora/routes"
)

func main() {

	database.InitDb()
	log.Println("✅ Migration complete, DB ready")
	r := gin.Default()
	routes.UserRoutes(r)
	routes.RoomRoutes(r)
	r.GET("/ws", controllers.HandleConnections)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	go controllers.HandleMessage()
	r.Run()
}
