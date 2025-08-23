package main

import (
    "log"
    "os"
    "card-vault/internal/config"
    "card-vault/internal/crypto"
    "card-vault/internal/handlers"
    "card-vault/internal/middleware"
    "card-vault/internal/repository"
    "card-vault/internal/service"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "golang.org/x/time/rate"
)

func main() {
    // Cargar variables de entorno
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found: %v", err)
    }

    // Inicializar base de datos
    db := config.InitDatabase()

    // Inicializar gestión de claves y cifrado
    keyManager := crypto.NewKeyManager()
    currentKey, _ := keyManager.GetCurrentKey()
    encryptionService, err := crypto.NewEncryptionService(currentKey)
    if err != nil {
        log.Fatal("Failed to initialize encryption service:", err)
    }

    // Inicializar capas
    cardRepo := repository.NewCardRepository(db)
    cardService := service.NewCardService(cardRepo, encryptionService, keyManager)
    cardHandler := handlers.NewCardHandler(cardService)

    // Configurar rate limiter
    rateLimiter := middleware.NewIPRateLimiter(rate.Limit(100), 20) // 100 requests per second, burst of 20

    // Configurar Gin
    if os.Getenv("GIN_MODE") == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.Default()

    // Middleware global
    r.Use(gin.Recovery())
    r.Use(middleware.RateLimitMiddleware(rateLimiter))
    r.Use(middleware.SecurityHeaders())
    r.Use(middleware.CORS())

    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })

    // API routes con autenticación
    api := r.Group("/api/v1")
    api.Use(middleware.AuthMiddleware())
    {
        cards := api.Group("/cards")
        {
            cards.POST("", cardHandler.CreateCard)
            cards.GET("", cardHandler.GetUserCards)
            cards.GET("/:id", cardHandler.GetCard)
            cards.PUT("/:id", cardHandler.UpdateCard)
            cards.DELETE("/:id", cardHandler.DeleteCard)
            cards.PATCH("/batch-update", cardHandler.BatchUpdateCards)
        }

        // Endpoint administrativo para rotación de claves
        admin := api.Group("/admin")
        {
            admin.POST("/cards/rotate-keys", cardHandler.RotateKeys)
        }
    }

    // Endpoint público para generar token (solo para testing)
    if os.Getenv("ENABLE_TEST_AUTH") == "true" {
        r.POST("/auth/test-token", handlers.GenerateTestToken)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}