package blockchain

import (
	"errors"
)

//////////////////////////////// Basic Block Structure //////////////////////////////////////
// type block
// data is data for block
// hash is sha256.Sum256([]byte(Data+PrevHash))
// prevHash is previous block's hash
// Height is id of block
type Block struct {
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
	Timestamp  int    `json:"timestamp"`

	Address     string `json:"address"`
	PrivateKey  string `json:"privateKey"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

var ErrNotFound = errors.New("Block not Found")
