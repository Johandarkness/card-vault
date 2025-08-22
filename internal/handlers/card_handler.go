package handlers

import (
    "net/http"
    "card-vault/internal/models"
    "card-vault/internal/service"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/google/uuid"
)

type CardHandler struct {
    cardService service.CardService
    validator   *validator.Validate
}

func NewCardHandler(cardService service.CardService) *CardHandler {
    return &CardHandler{
        cardService: cardService,
        validator:   validator.New(),
    }
}

// CreateCard - crea una tarjeta nueva
func (h *CardHandler) CreateCard(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var req models.CardRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    if err := h.validator.Struct(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    card, err := h.cardService.CreateCard(userID.(uuid.UUID), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, card)
}

// GetCard - obtiene una tarjeta por ID
func (h *CardHandler) GetCard(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    cardIDStr := c.Param("id")
    cardID, err := uuid.Parse(cardIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
        return
    }

    card, err := h.cardService.GetCard(userID.(uuid.UUID), cardID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, card)
}

// GetCards - obtiene todas las tarjetas de un usuario
func (h *CardHandler) GetCards(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    cards, err := h.cardService.GetCards(userID.(uuid.UUID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, cards)
}

// UpdateCard - actualiza una tarjeta existente
func (h *CardHandler) UpdateCard(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    cardIDStr := c.Param("id")
    cardID, err := uuid.Parse(cardIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
        return
    }

    var req models.CardRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    if err := h.validator.Struct(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updatedCard, err := h.cardService.UpdateCard(userID.(uuid.UUID), cardID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updatedCard)
}

// DeleteCard - elimina una tarjeta por ID
func (h *CardHandler) DeleteCard(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    cardIDStr := c.Param("id")
    cardID, err := uuid.Parse(cardIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
        return
    }

    if err := h.cardService.DeleteCard(userID.(uuid.UUID), cardID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}
