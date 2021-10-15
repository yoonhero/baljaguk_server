package blockchain

import (
	"strings"
	"time"

	"github.com/yoonhero/baljaguk_server/utils"
)

type UserBlock struct {
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

type UserData struct {
	Address     string `json:"address"`
	PrivateKey  string `json:"privateKey"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

// mine the block
func (b *UserBlock) mine() {
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
func (b *UserBlock) restore(data []byte) {
	utils.FromBytes(b, data)
}

// persist data
func persistUserBlock(b *UserBlock) {
	// db.SaveBlock(b.Hash, utils.ToBytes(b))
	dbStorage.SaveUserBlock(b.Hash, utils.ToBytes(b))
}

// find block by hash
func FindUserBlock(hash string) (*UserBlock, error) {
	// blockBytes := db.Block(hash)
	blockBytes := dbStorage.FindUserBlock(hash)

	// if that block don't exist
	if blockBytes == nil {
		// return nil with error
		return nil, ErrNotFound
	}

	block := &UserBlock{}
	// restore the block data
	block.restore(blockBytes)

	return block, nil
}

// create block
func createUserBlock(prevHash string, height int, diff int, data UserData) *UserBlock {
	block := &UserBlock{
		Hash:        "",
		PrevHash:    prevHash,
		Height:      height,
		Difficulty:  diff,
		Nonce:       0,
		Address:     data.Address,
		PrivateKey:  data.PrivateKey,
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
	}

	// block.Transactions = Mempool().TxToConfirm(from)

	// mining the block
	block.mine()

	// persist the block
	persistUserBlock(block)

	return block
}
