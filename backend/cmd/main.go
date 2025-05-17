package main

import (
	"net/http"
	"payd-mini-project/api"
	"payd-mini-project/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := db.NewSQLite("payd.db")
	if err != nil {
		panic(err)
	}

	// Existing Gin setup
	engine := gin.Default()

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, world!",
		})
	})

	shiftAPI := api.NewShiftAPI(db)
	shiftAPI.RegisterRoutes(engine)

	if err := engine.Run(":8080"); err != nil {
		panic(err)
	}
}
