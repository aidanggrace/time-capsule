package main

import (
	"time-capsule/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgresql://postgres.ylpdfqxoawztrtqksnag:BzfmdfyLU6cUP2jv@aws-1-us-east-2.pooler.supabase.com:5432/postgres"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// create table (migration)
	db.AutoMigrate(&models.Capsule{})

	r := gin.Default()

	// allow frontend to talk to backend
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Next()
	})

	// route
	r.GET("/capsules", func(c *gin.Context) {
		var capsules []models.Capsule
		db.Find(&capsules) // get all records
		c.JSON(200, capsules)
	})

	r.Run(":8080")
}
