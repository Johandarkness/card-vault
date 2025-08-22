package repository

import (
    "card-vault/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type CardRepository interface {
    Create(card *models.Card) error
    GetByID(id, userID uuid.UUID) (*models.Card, error)
    GetAllByUserID(userID uuid.UUID) ([]models.Card, error)
    Update(card *models.Card) error
    Delete(id, userID uuid.UUID) error
    BatchUpdate(cards []models.Card) error
    GetAllCards() ([]models.Card, error)
    UpdateKeyVersion(cardID uuid.UUID, version int) error
}

type cardRepository struct {
    db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
    return &cardRepository{db: db}
}

func (r *cardRepository) Create(card *models.Card) error {
    return r.db.Create(card).Error
}

func (r *cardRepository) GetByID(id, userID uuid.UUID) (*models.Card, error) {
    var card models.Card
    err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&card).Error
    return &card, err
}

func (r *cardRepository) GetAllByUserID(userID uuid.UUID) ([]models.Card, error) {
    var cards []models.Card
    err := r.db.Where("user_id = ?", userID).Find(&cards).Error
    return cards, err
}

func (r *cardRepository) Update(card *models.Card) error {
    return r.db.Save(card).Error
}

func (r *cardRepository) Delete(id, userID uuid.UUID) error {
    return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Card{}).Error
}

func (r *cardRepository) BatchUpdate(cards []models.Card) error {
    tx := r.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    for _, card := range cards {
        if err := tx.Save(&card).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    return tx.Commit().Error
}

func (r *cardRepository) GetAllCards() ([]models.Card, error) {
    var cards []models.Card
    err := r.db.Find(&cards).Error
    return cards, err
}

func (r *cardRepository) UpdateKeyVersion(cardID uuid.UUID, version int) error {
    return r.db.Model(&models.Card{}).Where("id = ?", cardID).Update("key_version", version).Error
}