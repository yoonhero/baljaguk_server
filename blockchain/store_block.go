package blockchain

import (
	"strings"
	"time"

	"github.com/yoonhero/baljaguk_server/utils"
)

type StoreBlock struct {
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
	Timestamp  int    `json:"timestamp"`

	Address     string `json:"address"`
	PrivateKey  string `json:"privateKey"`
	PhoneNumber string `json:"phoneNumber"`
}

type StoreData struct {
	Address     string `json:"address"`
	PrivateKey  string `json:"privateKey"`
	PhoneNumber string `json:"phoneNumber"`
}

// persist data
func persistStoreBlock(b *StoreBlock) {
	// db.SaveBlock(b.Hash, utils.ToBytes(b))
	dbStorage.SaveStoreBlock(b.Hash, utils.ToBytes(b))
}

// mine the block
func (b *StoreBlock) mine() {
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
func (b *StoreBlock) restore(data []byte) {
	utils.FromBytes(b, data)
}

// find block by hash
func FindStoreBlock(hash string) (*StoreBlock, error) {
	// blockBytes := db.Block(hash)
	blockBytes := dbStorage.FindStoreBlock(hash)

	// if that block don't exist
	if blockBytes == nil {
		// return nil with error
		return nil, ErrNotFound
	}

	block := &StoreBlock{}
	// restore the block data
	block.restore(blockBytes)

	return block, nil
}

// create block
func createStoreBlock(prevHash string, height int, diff int, data StoreData) *StoreBlock {
	block := &StoreBlock{
		Hash:        "",
		PrevHash:    prevHash,
		Height:      height,
		Difficulty:  diff,
		Nonce:       0,
		Address:     data.Address,
		PhoneNumber: data.PhoneNumber,
		PrivateKey:  data.PrivateKey,
	}

	// block.Transactions = Mempool().TxToConfirm(from)

	// mining the block
	block.mine()

	// persist the block
	persistStoreBlock(block)

	return block
}
