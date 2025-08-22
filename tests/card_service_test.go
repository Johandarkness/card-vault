package tests

import (
    "card-vault/internal/crypto"
    "card-vault/internal/models"
    "card-vault/internal/service"
    "crypto/rand"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type MockCardRepository struct {
    mock.Mock
}

func (m *MockCardRepository) Create(card *models.Card) error {
    args := m.Called(card)
    return args.Error(0)
}

func (m *MockCardRepository) GetByID(id, userID uuid.UUID) (*models.Card, error) {
    args := m.Called(id, userID)
    return args.Get(0).(*models.Card), args.Error(1)
}

func (m *MockCardRepository) GetAllByUserID(userID uuid.UUID) ([]models.Card, error) {
    args := m.Called(userID)
    return args.Get(0).([]models.Card), args.Error(1)
}

func (m *MockCardRepository) Update(card *models.Card) error {
    args := m.Called(card)
    return args.Error(0)
}

func (m *MockCardRepository) Delete(id, userID uuid.UUID) error {
    args := m.Called(id, userID)
    return args.Error(0)
}

func (m *MockCardRepository) BatchUpdate(cards []models.Card) error {
    args := m.Called(cards)
    return args.Error(0)
}

func (m *MockCardRepository) GetAllCards() ([]models.Card, error) {
    args := m.Called()
    return args.Get(0).([]models.Card), args.Error(1)
}

func (m *MockCardRepository) UpdateKeyVersion(cardID uuid.UUID, version int) error {
    args := m.Called(cardID, version)
    return args.Error(0)
}

func TestCardService_CreateCard(t *testing.T) {
    // Setup
    mockRepo := new(MockCardRepository)
    key := make([]byte, 32)
    rand.Read(key)
    encSvc, _ := crypto.NewEncryptionService(key)
    keyMgr := crypto.NewKeyManager()
    
    cardSvc := service.NewCardService(mockRepo, encSvc, keyMgr)

    userID := uuid.New()
    cardReq := &models.CardRequest{
        CardholderName: "John Doe",
        CardNumber:     "4111111111111111",
        ExpiryMonth:    12,
        ExpiryYear:     2025,
        CVV:            "123",
    }

    // Mock expectations
    mockRepo.On("Create", mock.AnythingOfType("*models.Card")).Return(nil)

    // Execute
    result, err := cardSvc.CreateCard(userID, cardReq)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "John Doe", result.CardholderName)
    assert.Equal(t, "************1111", result.MaskedNumber)
    assert.Equal(t, "Visa", result.CardType)

    mockRepo.AssertExpectations(t)
}