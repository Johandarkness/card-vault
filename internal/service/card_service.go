package service

import (
    "errors"
    "fmt"
    "regexp"
    "strings"
    "sync"
    "card-vault/internal/crypto"
    "card-vault/internal/models"
    "card-vault/internal/repository"
    
    "github.com/google/uuid"
)

type CardService interface {
    CreateCard(userID uuid.UUID, req *models.CardRequest) (*models.CardResponse, error)
    GetCard(cardID, userID uuid.UUID) (*models.CardResponse, error)
    GetUserCards(userID uuid.UUID) ([]models.CardResponse, error)
    UpdateCard(cardID, userID uuid.UUID, req *models.CardRequest) (*models.CardResponse, error)
    DeleteCard(cardID, userID uuid.UUID) error
    BatchUpdateCards(userID uuid.UUID, req *models.BatchUpdateRequest) ([]models.BatchUpdateResponse, error)
    RotateKeys() ([]models.BatchUpdateResponse, error)
}

type cardService struct {
    repo      repository.CardRepository
    encSvc    *crypto.EncryptionService
    keyMgr    *crypto.KeyManager
    mu        sync.RWMutex
}

func NewCardService(repo repository.CardRepository, encSvc *crypto.EncryptionService, keyMgr *crypto.KeyManager) CardService {
    return &cardService{
        repo:   repo,
        encSvc: encSvc,
        keyMgr: keyMgr,
    }
}

func (s *cardService) CreateCard(userID uuid.UUID, req *models.CardRequest) (*models.CardResponse, error) {
    cardNumber := strings.ReplaceAll(req.CardNumber, " ", "")
    if !s.isValidCardNumber(cardNumber) {
        return nil, errors.New("invalid card number")
    }

    encryptedNumber, err := s.encSvc.Encrypt(cardNumber)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt card number: %w", err)
    }

    encryptedCVV, err := s.encSvc.Encrypt(req.CVV)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt CVV: %w", err)
    }

    _, keyVersion := s.keyMgr.GetCurrentKey()

    card := &models.Card{
        UserID:         userID,
        CardholderName: req.CardholderName,
        CardNumber:     encryptedNumber,
        ExpiryMonth:    req.ExpiryMonth,
        ExpiryYear:     req.ExpiryYear,
        CVV:            encryptedCVV,
        CardType:       s.detectCardType(cardNumber),
        KeyVersion:     keyVersion,
    }

    if err := s.repo.Create(card); err != nil {
        return nil, fmt.Errorf("failed to create card: %w", err)
    }

    return s.toCardResponse(card, cardNumber), nil
}

func (s *cardService) GetCard(cardID, userID uuid.UUID) (*models.CardResponse, error) {
    card, err := s.repo.GetByID(cardID, userID)
    if err != nil {
        return nil, fmt.Errorf("card not found: %w", err)
    }

    decryptedNumber, err := s.decryptCardNumber(card)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt card data: %w", err)
    }

    return s.toCardResponse(card, decryptedNumber), nil
}

func (s *cardService) GetUserCards(userID uuid.UUID) ([]models.CardResponse, error) {
    cards, err := s.repo.GetAllByUserID(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user cards: %w", err)
    }

    responses := make([]models.CardResponse, len(cards))
    for i, card := range cards {
        decryptedNumber, err := s.decryptCardNumber(&card)
        if err != nil {
            return nil, fmt.Errorf("failed to decrypt card data: %w", err)
        }
        responses[i] = *s.toCardResponse(&card, decryptedNumber)
    }

    return responses, nil
}

func (s *cardService) UpdateCard(cardID, userID uuid.UUID, req *models.CardRequest) (*models.CardResponse, error) {
    card, err := s.repo.GetByID(cardID, userID)
    if err != nil {
        return nil, fmt.Errorf("card not found: %w", err)
    }

    cardNumber := strings.ReplaceAll(req.CardNumber, " ", "")
    if !s.isValidCardNumber(cardNumber) {
        return nil, errors.New("invalid card number")
    }

    encryptedNumber, err := s.encSvc.Encrypt(cardNumber)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt card number: %w", err)
    }

    encryptedCVV, err := s.encSvc.Encrypt(req.CVV)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt CVV: %w", err)
    }

    card.CardholderName = req.CardholderName
    card.CardNumber = encryptedNumber
    card.ExpiryMonth = req.ExpiryMonth
    card.ExpiryYear = req.ExpiryYear
    card.CVV = encryptedCVV
    card.CardType = s.detectCardType(cardNumber)

    if err := s.repo.Update(card); err != nil {
        return nil, fmt.Errorf("failed to update card: %w", err)
    }

    return s.toCardResponse(card, cardNumber), nil
}

func (s *cardService) DeleteCard(cardID, userID uuid.UUID) error {
    return s.repo.Delete(cardID, userID)
}

func (s *cardService) BatchUpdateCards(userID uuid.UUID, req *models.BatchUpdateRequest) ([]models.BatchUpdateResponse, error) {
    responses := make([]models.BatchUpdateResponse, len(req.Cards))
    var wg sync.WaitGroup
    
    for i, updateReq := range req.Cards {
        wg.Add(1)
        go func(index int, cardUpdate models.BatchCardUpdate) {
            defer wg.Done()
            
            card, err := s.repo.GetByID(cardUpdate.ID, userID)
            if err != nil {
                responses[index] = models.BatchUpdateResponse{
                    CardID: cardUpdate.ID,
                    Status: "failed",
                    Error:  "card not found",
                }
                return
            }

            if cardUpdate.CardholderName != nil {
                card.CardholderName = *cardUpdate.CardholderName
            }
            if cardUpdate.ExpiryMonth != nil {
                card.ExpiryMonth = *cardUpdate.ExpiryMonth
            }
            if cardUpdate.ExpiryYear != nil {
                card.ExpiryYear = *cardUpdate.ExpiryYear
            }

            if err := s.repo.Update(card); err != nil {
                responses[index] = models.BatchUpdateResponse{
                    CardID: cardUpdate.ID,
                    Status: "failed",
                    Error:  err.Error(),
                }
                return
            }

            responses[index] = models.BatchUpdateResponse{
                CardID: cardUpdate.ID,
                Status: "success",
            }
        }(i, updateReq)
    }
    
    wg.Wait()
    return responses, nil
}

func (s *cardService) RotateKeys() ([]models.BatchUpdateResponse, error) {
    cards, err := s.repo.GetAllCards()
    if err != nil {
        return nil, fmt.Errorf("failed to get cards: %w", err)
    }

    if err := s.keyMgr.RotateKey(); err != nil {
        return nil, fmt.Errorf("failed to rotate key: %w", err)
    }

    newKey, newVersion := s.keyMgr.GetCurrentKey()
    newEncSvc, err := crypto.NewEncryptionService(newKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create new encryption service: %w", err)
    }

    responses := make([]models.BatchUpdateResponse, len(cards))
    oldKey := s.keyMgr.GetPreviousKey()
    oldEncSvc, _ := crypto.NewEncryptionService(oldKey)

    for i, card := range cards {
        cardNumber, err := oldEncSvc.Decrypt(card.CardNumber)
        if err != nil {
            responses[i] = models.BatchUpdateResponse{
                CardID: card.ID,
                Status: "failed",
                Error:  "failed to decrypt with old key",
            }
            continue
        }

        cvv, err := oldEncSvc.Decrypt(card.CVV)
        if err != nil {
            responses[i] = models.BatchUpdateResponse{
                CardID: card.ID,
                Status: "failed",
                Error:  "failed to decrypt CVV with old key",
            }
            continue
        }

        newEncryptedNumber, err := newEncSvc.Encrypt(cardNumber)
        if err != nil {
            responses[i] = models.BatchUpdateResponse{
                CardID: card.ID,
                Status: "failed",
                Error:  "failed to encrypt with new key",
            }
            continue
        }

        newEncryptedCVV, err := newEncSvc.Encrypt(cvv)
        if err != nil {
            responses[i] = models.BatchUpdateResponse{
                CardID: card.ID,
                Status: "failed",
                Error:  "failed to encrypt CVV with new key",
            }
            continue
        }

        card.CardNumber = newEncryptedNumber
        card.CVV = newEncryptedCVV
        card.KeyVersion = newVersion

        if err := s.repo.Update(&card); err != nil {
            responses[i] = models.BatchUpdateResponse{
                CardID: card.ID,
                Status: "failed",
                Error:  "failed to update card",
            }
            continue
        }

        responses[i] = models.BatchUpdateResponse{
            CardID: card.ID,
            Status: "success",
        }
    }

    s.mu.Lock()
    s.encSvc = newEncSvc
    s.mu.Unlock()

    return responses, nil
}

func (s *cardService) isValidCardNumber(cardNumber string) bool {
    var sum int
    alternate := false
    
    for i := len(cardNumber) - 1; i >= 0; i-- {
        digit := int(cardNumber[i] - '0')
        if digit < 0 || digit > 9 {
            return false
        }
        
        if alternate {
            digit *= 2
            if digit > 9 {
                digit = digit%10 + digit/10
            }
        }
        
        sum += digit
        alternate = !alternate
    }
    
    return sum%10 == 0
}

func (s *cardService) detectCardType(cardNumber string) string {
    patterns := map[string]*regexp.Regexp{
        "Visa":       regexp.MustCompile(`^4[0-9]{12}(?:[0-9]{3})?$`),
        "Mastercard": regexp.MustCompile(`^5[1-5][0-9]{14}$|^2(?:2(?:2[1-9]|[3-9][0-9])|[3-6][0-9][0-9]|7(?:[01][0-9]|20))[0-9]{12}$`),
        "Amex":       regexp.MustCompile(`^3[47][0-9]{13}$`),
        "Discover":   regexp.MustCompile(`^6(?:011|5[0-9]{2})[0-9]{12}$`),
    }

    for cardType, pattern := range patterns {
        if pattern.MatchString(cardNumber) {
            return cardType
        }
    }
    
    return "Unknown"
}

func (s *cardService) maskCardNumber(cardNumber string) string {
    if len(cardNumber) < 4 {
        return cardNumber
    }
    masked := strings.Repeat("*", len(cardNumber)-4)
    return masked + cardNumber[len(cardNumber)-4:]
}

func (s *cardService) toCardResponse(card *models.Card, decryptedNumber string) *models.CardResponse {
    return &models.CardResponse{
        ID:             card.ID,
        UserID:         card.UserID,
        CardholderName: card.CardholderName,
        MaskedNumber:   s.maskCardNumber(decryptedNumber),
        ExpiryMonth:    card.ExpiryMonth,
        ExpiryYear:     card.ExpiryYear,
        CardType:       card.CardType,
        IsActive:       card.IsActive,
        CreatedAt:      card.CreatedAt,
        UpdatedAt:      card.UpdatedAt,
    }
}

func (s *cardService) decryptCardNumber(card *models.Card) (string, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    _, currentVersion := s.keyMgr.GetCurrentKey()
    
    if card.KeyVersion == currentVersion {
        return s.encSvc.Decrypt(card.CardNumber)
    }
    
    previousKey := s.keyMgr.GetPreviousKey()
    if previousKey != nil {
        oldEncSvc, err := crypto.NewEncryptionService(previousKey)
        if err != nil {
            return "", err
        }
        return oldEncSvc.Decrypt(card.CardNumber)
    }
    
    return "", errors.New("unable to decrypt card with available keys")
}