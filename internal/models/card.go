package models

import (
    "time"
    "github.com/google/uuid"
)

type Card struct {
    ID              uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID          uuid.UUID `json:"user_id" gorm:"not null;index"`
    CardholderName  string    `json:"cardholder_name" gorm:"not null" validate:"required,min=1,max=100"`
    CardNumber      string    `json:"-" gorm:"not null"`
    ExpiryMonth     int       `json:"expiry_month" validate:"required,min=1,max=12"`
    ExpiryYear      int       `json:"expiry_year" validate:"required,min=2024"`
    CVV             string    `json:"-" gorm:"not null"`
    CardType        string    `json:"card_type" gorm:"not null"`
    IsActive        bool      `json:"is_active" gorm:"default:true"`
    KeyVersion      int       `json:"-" gorm:"not null;default:1"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type CardResponse struct {
    ID             uuid.UUID `json:"id"`
    UserID         uuid.UUID `json:"user_id"`
    CardholderName string    `json:"cardholder_name"`
    MaskedNumber   string    `json:"masked_number"`
    ExpiryMonth    int       `json:"expiry_month"`
    ExpiryYear     int       `json:"expiry_year"`
    CardType       string    `json:"card_type"`
    IsActive       bool      `json:"is_active"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type CardRequest struct {
    CardholderName string `json:"cardholder_name" validate:"required,min=1,max=100"`
    CardNumber     string `json:"card_number" validate:"required,creditcard"`
    ExpiryMonth    int    `json:"expiry_month" validate:"required,min=1,max=12"`
    ExpiryYear     int    `json:"expiry_year" validate:"required,min=2024"`
    CVV            string `json:"cvv" validate:"required,len=3"`
}

type BatchUpdateRequest struct {
    Cards []BatchCardUpdate `json:"cards" validate:"required,dive"`
}

type BatchCardUpdate struct {
    ID             uuid.UUID `json:"id" validate:"required"`
    CardholderName *string   `json:"cardholder_name,omitempty"`
    ExpiryMonth    *int      `json:"expiry_month,omitempty" validate:"omitempty,min=1,max=12"`
    ExpiryYear     *int      `json:"expiry_year,omitempty" validate:"omitempty,min=2024"`
}

type BatchUpdateResponse struct {
    CardID  uuid.UUID `json:"card_id"`
    Status  string    `json:"status"`
    Error   string    `json:"error,omitempty"`
}