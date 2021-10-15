package blockchain

import (
	"strings"
	"time"

	"github.com/yoonhero/baljaguk_server/utils"
)

type BaljagukBlock struct {
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
	Timestamp  int    `json:"timestamp"`
	StoreHash  string `json:"storehash"`
	UserHash   string `json:"userhash"`
}

type BaljagukData struct {
	StoreHash string `json:"storehash"`
	UserHash  string `json:"userhash"`
}

// persist data
func persistBaljagukBlock(b *BaljagukBlock) {
	// db.SaveBlock(b.Hash, utils.ToBytes(b))
	dbStorage.SaveBaljagukBlock(b.Hash, utils.ToBytes(b))
}

// mine the block
func (b *BaljagukBlock) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			// fmt.Printf("Target:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
			break
		} else {
			b.Nonce++
		}
	}
}

// restore data
func (b *BaljagukBlock) restore(data []byte) {
	utils.FromBytes(b, data)
}

// find block by hash
func FindBaljagukBlock(hash string) (*BaljagukBlock, error) {
	// blockBytes := db.Block(hash)
	blockBytes := dbStorage.FindBaljaguk(hash)

	// if that block don't exist
	if blockBytes == nil {
		// return nil with error
		return nil, ErrNotFound
	}

	block := &BaljagukBlock{}
	// restore the block data
	block.restore(blockBytes)

	return block, nil
}

// create block
func createBaljagukBlock(prevHash string, height int, diff int, data BaljagukData) *BaljagukBlock {
	block := &BaljagukBlock{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
		StoreHash:  data.StoreHash,
		UserHash:   data.UserHash,
	}

	// block.Transactions = Mempool().TxToConfirm(from)

	// mining the block
	block.mine()

	// persist the block
	persistBaljagukBlock(block)

	return block
}
