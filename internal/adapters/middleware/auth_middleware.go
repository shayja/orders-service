package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "msg": "Authorization token required"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure the token's signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			fmt.Println("JWT", err, token.Valid)
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "msg": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract the user ID from the token (assuming it's in the 'sub' field)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["sub"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "msg": "Invalid token claims"})
			c.Abort()
			return
		}

		// Set the user ID in the context to use it later in the handler
		userId := claims["sub"].(string)
		c.Set("userId", userId)

		// Continue with the next middleware/handler
		c.Next()
	}
}
