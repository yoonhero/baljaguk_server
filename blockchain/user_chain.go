package blockchain

import (
	"sync"

	"github.com/yoonhero/baljaguk_server/db"
	"github.com/yoonhero/baljaguk_server/utils"
)

var userDBStorage storage = db.DB{}

// variable struct that play func only one time
var userOnce sync.Once

// type blockchain
// blocks is slice of []Block
type userBlockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

func (b *userBlockchain) LockBlockchain() {
	b.m.Lock()
	defer b.m.Unlock()
}

func (b *userBlockchain) restore(data []byte) {
	// decoder := gob.NewDecoder(bytes.NewReader(data))
	// decoder.Decode(b)
	utils.FromBytes(b, data)
}

// variable blockchain pointers
var user_b *userBlockchain

// add block to blockchain
func (b *userBlockchain) AddUserBlock(Address string, PrivateKey string, PhoneNumber string, Email string) *UserBlock {
	b.LockBlockchain()

	data := UserData{
		Address:     Address,
		PrivateKey:  PrivateKey,
		PhoneNumber: PhoneNumber,
		Email:       Email,
	}
	// createBlock
	block := createUserBlock(b.NewestHash, b.Height+1, getUserDifficulty(b), data)

	// set newesthash new block's hash
	b.NewestHash = block.Hash
	// set height new block's height
	b.Height = block.Height

	b.CurrentDifficulty = block.Difficulty

	// persist the blockchain
	persistUserBlockchain(b)
	return block
}

// all blocks
func UserBlocks(b *userBlockchain) []*UserBlock {
	b.LockBlockchain()
	var blocks []*UserBlock

	// start newesthash and its prevhash and find block
	// if prevhash dont exist = genesis block break
	hashCursor := b.NewestHash

	for {
		block, _ := FindUserBlock(hashCursor)
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
func persistUserBlockchain(b *userBlockchain) {
	// db.SaveCheckpoint(utils.ToBytes(b))
	userDBStorage.SaveUserChain((utils.ToBytes(b)))
}

func UserBlockchain() *userBlockchain {
	// run only one time
	userOnce.Do(func() {
		// initial blockchain struct
		user_b = &userBlockchain{Height: 0}

		// search for checkpoint on the db
		// checkpoint := db.Checkpoint()
		checkpoint := userDBStorage.LoadUserChain()

		if checkpoint == nil {
			// if blockchain don't exist create block
			user_b.AddUserBlock("", "", "", "")
		} else {
			// restore data from db
			user_b.restore(checkpoint)
		}
	})
	// return type blockchain struct
	return user_b
}

// recalculate difficulty of block by timestamp
func recalculateUserDifficulty(b *userBlockchain) int {
	// get all blocks
	allBlocks := UserBlocks(b)
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

func getUserDifficulty(b *userBlockchain) int {
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
