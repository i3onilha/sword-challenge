package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// contextKey type prevents collisions
type contextKey string

const UserIDKey contextKey = "userID"

// Custom error types for better error handling
type AuthError struct {
	Message string
	Err     error
}

func (e *AuthError) Error() string {
	return e.Message
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

// ContainsRole checks if the user has at least one required role
func ContainsRole(userRoles []string, requiredRoles []string) bool {
	roleSet := make(map[string]struct{})
	for _, r := range userRoles {
		roleSet[r] = struct{}{}
	}
	for _, req := range requiredRoles {
		if _, ok := roleSet[req]; ok {
			return true
		}
	}
	return false
}

// getJWTSecret returns the JWT secret key from environment variables
func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, &AuthError{Message: "authentication not configured"}
	}
	if len(secret) < 32 {
		return nil, &AuthError{Message: "authentication configuration error"}
	}
	return []byte(secret), nil
}

// GinAuthMiddleware returns a Gin middleware that validates JWT tokens
func GinAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret, err := getJWTSecret()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication error"})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("userID", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

// RequireRole returns a Gin middleware that checks if the user has at least one of the required roles
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication error"})
			c.Abort()
			return
		}

		if !ContainsRole(userRoles, requiredRoles) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
