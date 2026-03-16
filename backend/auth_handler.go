package main

import (
	"context"
	"net/http"
	"os"
	"time"
	"time-capsule/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Simple struct for both login and register
type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func Register(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AuthInput
		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Bad request"})
			return
		}

		// hash the password
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

		_, err := queries.CreateUser(context.Background(), db.CreateUserParams{
			Email:        input.Email,
			PasswordHash: string(hashed),
		})

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not create user"})
			return
		}

		c.JSON(200, gin.H{"message": "User created!"})
	}
}

func Login(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AuthInput
		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Bad request"})
			return
		}

		user, err := queries.GetUserByEmail(context.Background(), input.Email)
		if err != nil {
			c.JSON(401, gin.H{"error": "User not found"})
			return
		}

		// check password
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
		if err != nil {
			c.JSON(401, gin.H{"error": "Wrong password"})
			return
		}

		// create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, _ := token.SignedString(secretKey)

		// send cookie
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 0, "/", "", false, true)

		c.JSON(200, gin.H{"message": "Logged in!"})
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if token == nil || !token.Valid {
			c.AbortWithStatus(401)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["sub"])

		c.Next()
	}
}
