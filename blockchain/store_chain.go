package blockchain

import "github.com/yoonhero/baljaguk_server/utils"

// variable blockchain pointers
var store_b *blockchain

// add block to blockchain
func (b *blockchain) AddStoreBlock(Address string, PrivateKey string, PhoneNumber string) *StoreBlock {
	data := StoreData{
		Address:     Address,
		PrivateKey:  PrivateKey,
		PhoneNumber: PhoneNumber,
	}
	// createBlock
	block := createStoreBlock(b.NewestHash, b.Height+1, getStoreDifficulty(b), data)

	// set newesthash new block's hash
	b.NewestHash = block.Hash
	// set height new block's height
	b.Height = block.Height

	b.CurrentDifficulty = block.Difficulty

	// persist the blockchain
	persistStoreBlockchain(b)
	return block
}

// all blocks
func StoreBlocks(b *blockchain) []*StoreBlock {
	b.m.Lock()
	defer b.m.Unlock()
	var blocks []*StoreBlock

	// start newesthash and its prevhash and find block
	// if prevhash dont exist = genesis block break
	hashCursor := b.NewestHash

	for {
		block, _ := FindStoreBlock(hashCursor)
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
func persistStoreBlockchain(b *blockchain) {
	// db.SaveCheckpoint(utils.ToBytes(b))
	dbStorage.SaveStoreChain((utils.ToBytes(b)))
}

func StoreBlockchain() *blockchain {
	// run only one time
	once.Do(func() {
		// initial blockchain struct
		store_b = &blockchain{Height: 0}

		// search for checkpoint on the db
		// checkpoint := db.Checkpoint()
		checkpoint := dbStorage.LoadStoreChain()

		if checkpoint == nil {
			// if blockchain don't exist create block
			store_b.AddStoreBlock("", "", "")
		} else {
			// restore data from db
			store_b.restore(checkpoint)
		}
	})
	// return type blockchain struct
	return store_b
}

// recalculate difficulty of block by timestamp
func recalculateStoreDifficulty(b *blockchain) int {
	// get all blocks
	allBlocks := StoreBlocks(b)
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

func getStoreDifficulty(b *blockchain) int {
	// if genesis block or not
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return recalculateUserDifficulty(b)
	} else {
		if b.CurrentDifficulty <= 5 {
			return b.CurrentDifficulty
		} else {
			return 5
		}
	}
}
