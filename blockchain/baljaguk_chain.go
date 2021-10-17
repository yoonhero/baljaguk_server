package blockchain

import (
	"sync"

	"github.com/yoonhero/baljaguk_server/db"
	"github.com/yoonhero/baljaguk_server/utils"
)

var baljagukDBStorage storage = db.DB{}

// variable struct that play func only one time
var baljagukOnce sync.Once

// type blockchain
// blocks is slice of []Block
type baljagukBlockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

func (b *baljagukBlockchain) LockBlockchain() {
	b.m.Lock()
	defer b.m.Unlock()
}

func (b *baljagukBlockchain) restore(data []byte) {
	// decoder := gob.NewDecoder(bytes.NewReader(data))
	// decoder.Decode(b)
	utils.FromBytes(b, data)
}

// variable blockchain pointers
var baljaguk_b *baljagukBlockchain

// add block to blockchain
func (b *baljagukBlockchain) AddBaljagukBlock(StoreHash string, UserHash string, Latitude string, Longitude string) *BaljagukBlock {
	b.LockBlockchain()

	data := BaljagukData{
		StoreHash: StoreHash,
		UserHash:  UserHash,
		Latitude:  Latitude,
		Longitude: Longitude,
	}
	// createBlock
	block := createBaljagukBlock(b.NewestHash, b.Height+1, getBaljagukDifficulty(b), data)

	// set newesthash new block's hash
	b.NewestHash = block.Hash
	// set height new block's height
	b.Height = block.Height

	b.CurrentDifficulty = block.Difficulty

	// persist the blockchain
	persistBaljagukBlockchain(b)
	return block
}

// all blocks
func BaljagukBlocks(b *baljagukBlockchain) []*BaljagukBlock {
	b.LockBlockchain()

	var blocks []*BaljagukBlock

	// start newesthash and its prevhash and find block
	// if prevhash dont exist = genesis block break
	hashCursor := b.NewestHash

	for {
		block, _ := FindBaljagukBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

// persist the blockchain data
func persistBaljagukBlockchain(b *baljagukBlockchain) {
	// db.SaveCheckpoint(utils.ToBytes(b))
	baljagukDBStorage.SaveBaljagukChain((utils.ToBytes(b)))
}

func BaljagukBlockchain() *baljagukBlockchain {
	// run only one time
	baljagukOnce.Do(func() {
		// initial blockchain struct
		baljaguk_b = &baljagukBlockchain{Height: 0}

		// search for checkpoint on the db
		// checkpoint := db.Checkpoint()
		checkpoint := baljagukDBStorage.LoadBaljagukChain()

		if checkpoint == nil {
			// if blockchain don't exist create block
			baljaguk_b.AddBaljagukBlock("", "", "", "")
		} else {
			// reBaljaguk data from db
			baljaguk_b.restore(checkpoint)
		}
	})
	// return type blockchain struct
	return baljaguk_b
}

// recalculate difficulty of block by timestamp
func recalculateBaljagukDifficulty(b *baljagukBlockchain) int {
	// get all blocks
	allBlocks := BaljagukBlocks(b)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if b.CurrentDifficulty > 5 {
		return b.CurrentDifficulty - 1
	}
	if actualTime <= (expectedTime - allowedRange) {
		// if acuaultime < 8 difficulty + 1
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		// if actualtime >= 12 difficulty - 1
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func getBaljagukDifficulty(b *baljagukBlockchain) int {
	// if genesis block or not
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return recalculateBaljagukDifficulty(b)
	} else {
		if b.CurrentDifficulty <= 5 {
			return b.CurrentDifficulty
		} else {
			return 5
		}
	}
}
