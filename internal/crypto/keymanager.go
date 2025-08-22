package crypto

import (
    "crypto/rand"
    "sync"
    "time"
)

type KeyManager struct {
    currentKey    []byte
    previousKey   []byte
    keyVersion    int
    rotationTime  time.Time
    mu           sync.RWMutex
}

func NewKeyManager() *KeyManager {
    key := make([]byte, 32)
    rand.Read(key)
    
    return &KeyManager{
        currentKey:   key,
        keyVersion:   1,
        rotationTime: time.Now(),
    }
}

func (km *KeyManager) GetCurrentKey() ([]byte, int) {
    km.mu.RLock()
    defer km.mu.RUnlock()
    return km.currentKey, km.keyVersion
}

func (km *KeyManager) RotateKey() error {
    km.mu.Lock()
    defer km.mu.Unlock()
    
    newKey := make([]byte, 32)
    if _, err := rand.Read(newKey); err != nil {
        return err
    }
    
    km.previousKey = km.currentKey
    km.currentKey = newKey
    km.keyVersion++
    km.rotationTime = time.Now()
    
    return nil
}

func (km *KeyManager) GetPreviousKey() []byte {
    km.mu.RLock()
    defer km.mu.RUnlock()
    return km.previousKey
}