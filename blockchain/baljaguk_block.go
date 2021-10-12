import (
	"errors"
	"strings"
	"time"

	"github.com/yoonhero/baljaguk_server/utils"
)

// persist data
func persistBaljagukBlock(b *Block) {
	// db.SaveBlock(b.Hash, utils.ToBytes(b))
	dbStorage.SaveBaljagukBlock(b.Hash, utils.ToBytes(b))
}

var ErrNotFound = errors.New("Block not Found")

// find block by hash
func FindBaljagukBlock(hash string) (*Block, error) {
	// blockBytes := db.Block(hash)
	blockBytes := dbStorage.FindBaljagukBlock(hash)

	// if that block don't exist
	if blockBytes == nil {
		// return nil with error
		return nil, ErrNotFound
	}

	block := &Block{}
	// restore the block data
	block.restore(blockBytes)

	return block, nil
}

// create block
func createBaljagukBlock(prevHash string, height int, diff int, from string) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
	}

	// block.Transactions = Mempool().TxToConfirm(from)

	// mining the block
	block.mine()

	// persist the block
	persistBaljagukBlock(block)

	return block
}
