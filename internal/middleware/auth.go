package middleware

import (
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/repository"
    "github.com/bot011max/medical-bot/internal/security"
)

func AuthMiddleware(userRepo *repository.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }
        
        token := parts[1]
        claims, err := security.ValidateJWT(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }
        
        userID, ok := claims["user_id"].(string)
        if !ok {
            c.JSON(401, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }
        
        user, err := userRepo.FindByID(userID)
        if err != nil {
            c.JSON(401, gin.H{"error": "User not found"})
            c.Abort()
            return
        }
        
        c.Set("user_id", user.ID)
        c.Set("user", user)
        c.Next()
    }
}
