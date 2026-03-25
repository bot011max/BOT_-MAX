package audit

import (
    "crypto/sha256"
    "encoding/hex"
    "sync"
    "time"
)

type Block struct {
    Index        int       `json:"index"`
    Timestamp    time.Time `json:"timestamp"`
    EventType    string    `json:"event_type"`
    UserID       string    `json:"user_id"`
    Action       string    `json:"action"`
    Details      string    `json:"details"`
    PreviousHash string    `json:"previous_hash"`
    Hash         string    `json:"hash"`
}

type Blockchain struct {
    chain []Block
    mu    sync.RWMutex
}

func NewBlockchain() *Blockchain {
    bc := &Blockchain{
        chain: make([]Block, 0),
    }
    bc.createGenesisBlock()
    return bc
}

func (bc *Blockchain) createGenesisBlock() {
    genesis := Block{
        Index:        0,
        Timestamp:    time.Now(),
        EventType:    "GENESIS",
        UserID:       "system",
        Action:       "init",
        Details:      "Blockchain audit system initialized",
        PreviousHash: "0",
    }
    genesis.Hash = bc.calculateHash(genesis)
    bc.chain = append(bc.chain, genesis)
}

func (bc *Blockchain) calculateHash(block Block) string {
    data := string(rune(block.Index)) +
        block.Timestamp.String() +
        block.EventType +
        block.UserID +
        block.Action +
        block.Details +
        block.PreviousHash
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) AddEvent(eventType, userID, action, details string) {
    bc.mu.Lock()
    defer bc.mu.Unlock()

    lastBlock := bc.chain[len(bc.chain)-1]
    newBlock := Block{
        Index:        len(bc.chain),
        Timestamp:    time.Now(),
        EventType:    eventType,
        UserID:       userID,
        Action:       action,
        Details:      details,
        PreviousHash: lastBlock.Hash,
    }
    newBlock.Hash = bc.calculateHash(newBlock)
    bc.chain = append(bc.chain, newBlock)
}

func (bc *Blockchain) Verify() bool {
    bc.mu.RLock()
    defer bc.mu.RUnlock()

    for i := 1; i < len(bc.chain); i++ {
        current := bc.chain[i]
        previous := bc.chain[i-1]

        if current.PreviousHash != previous.Hash {
            return false
        }

        if current.Hash != bc.calculateHash(current) {
            return false
        }
    }
    return true
}

func (bc *Blockchain) GetEvents() []Block {
    bc.mu.RLock()
    defer bc.mu.RUnlock()
    return bc.chain
}
