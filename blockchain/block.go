package blockchain

import (
	"github.com/yoonhero/baljaguk_server/utils"
)

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
}

// restore data
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}
