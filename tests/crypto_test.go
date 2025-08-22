package tests

import (
    "card-vault/internal/crypto"
    "crypto/rand"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestEncryptionService(t *testing.T) {
    // Generate test key
    key := make([]byte, 32)
    rand.Read(key)

    encSvc, err := crypto.NewEncryptionService(key)
    assert.NoError(t, err)

    testData := "4111111111111111"

    // Test encryption
    encrypted, err := encSvc.Encrypt(testData)
    assert.NoError(t, err)
    assert.NotEmpty(t, encrypted)
    assert.NotEqual(t, testData, encrypted)

    // Test decryption
    decrypted, err := encSvc.Decrypt(encrypted)
    assert.NoError(t, err)
    assert.Equal(t, testData, decrypted)
}

func TestKeyManager(t *testing.T) {
    km := crypto.NewKeyManager()

    currentKey, version := km.GetCurrentKey()
    assert.NotNil(t, currentKey)
    assert.Equal(t, 1, version)

    // Test key rotation
    err := km.RotateKey()
    assert.NoError(t, err)

    newKey, newVersion := km.GetCurrentKey()
    assert.NotNil(t, newKey)
    assert.Equal(t, 2, newVersion)
    assert.NotEqual(t, currentKey, newKey)

    previousKey := km.GetPreviousKey()
    assert.Equal(t, currentKey, previousKey)
}