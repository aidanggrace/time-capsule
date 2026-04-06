package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"time-capsule/db"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	InitJWKS()

	r := gin.Default()

	// allow frontend to talk to backend
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/capsules", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				return
			}
			ownerID := fmt.Sprintf("%v", userID)

			capsules, err := queries.ListCapsules(c, ownerID)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, capsules)
		})

		protected.GET("/received-capsules", func(c *gin.Context) {
			email, _ := c.Get("userEmail")
			emailStr := fmt.Sprintf("%v", email)

			capsules, err := queries.ListReceivedCapsules(c, emailStr)
			if err != nil {
				fmt.Println("received-capsules error:", err)
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, capsules)
		})

		protected.POST("/capsules", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				return
			}

			var input struct {
				RecipientEmail string `json:"recipient_email"`
				Message        string `json:"message"`
				UnlockAt       string `json:"unlock_at"`
			}
			if err := c.BindJSON(&input); err != nil {
				c.JSON(400, gin.H{"error": "Bad request"})
				return
			}

			unlockTime, err := time.Parse(time.RFC3339, input.UnlockAt)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid date format"})
				return
			}

			senderEmail, _ := c.Get("userEmail")
			capsule, err := queries.CreateCapsule(c, db.CreateCapsuleParams{
				OwnerID:        fmt.Sprintf("%v", userID),
				SenderEmail:    fmt.Sprintf("%v", senderEmail),
				RecipientEmail: input.RecipientEmail,
				Message:        input.Message,
				MediaUrl:       "",
				UnlockAt:       unlockTime,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(201, capsule)
		})
	}

	r.Run(":8080")
}
