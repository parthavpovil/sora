package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthavpovil/sora/database"
	"github.com/parthavpovil/sora/routes"
)

func main() {

	database.InitDb()
	log.Println("✅ Migration complete, DB ready")
	r := gin.Default()
	routes.UserRoutes(r)
	
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
