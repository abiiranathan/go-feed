package main

import (
	"net/http"

	"github.com/abiiranathan/gin_api/db"
	"github.com/abiiranathan/gin_api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	db.Connect()
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.43.223"})

	db := db.DBConn

	r.LoadHTMLGlob("./public/templates/*")
	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	routes.SetupRoutes(r, db)

	r.Run(":8080")
}
