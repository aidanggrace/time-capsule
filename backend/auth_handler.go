package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwks *keyfunc.JWKS

func InitJWKS() {
	supabaseURL := os.Getenv("SUPABASE_URL")
	jwksURL := supabaseURL + "/auth/v1/.well-known/jwks.json"

	var err error
	jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		log.Fatalf("Failed to get JWKS: %v", err)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)

		if err != nil || !token.Valid {
			fmt.Println("JWT error:", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Set("userID", claims["sub"])
		c.Set("userEmail", claims["email"])
		c.Next()
	}
}
