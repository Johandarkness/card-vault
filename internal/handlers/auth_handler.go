package handlers

import (
    "net/http"
    "card-vault/internal/middleware"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

func GenerateTestToken(c *gin.Context) {
    // Solo para desarrollo - en producción usar un sistema de auth real
    userID := uuid.New()
    
    token, err := middleware.GenerateToken(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token":   token,
        "user_id": userID,
        "message": "Test token generated successfully",
    })
}