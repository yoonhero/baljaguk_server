// persistence of block
// connected to DB to save data
// using bolt DB (bitcoin levelDB)
package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/yoonhero/baljaguk_server/db"
	"github.com/yoonhero/baljaguk_server/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 1
	allowedRange       int = 2
)

// type blockchain
// blocks is slice of []Block
type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

type storage interface {
	FindBlock(hash string) []byte
	LoadChain() []byte
	SaveBlock(hash string, data []byte)
	SaveChain(data []byte)
	DeleteAllBlocks()
}

// variable blockchain pointers
var b *blockchain

// variable struct that play func only one time
var once sync.Once

var dbStorage storage = db.DB{}

func (b *blockchain) restore(data []byte) {
	// decoder := gob.NewDecoder(bytes.NewReader(data))
	// decoder.Decode(b)
	utils.FromBytes(b, data)
}

// add block to blockchain
func (b *blockchain) AddBlock(from string) *Block {
	// createBlock
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b), from)

	// set newesthash new block's hash
	b.NewestHash = block.Hash
	// set height new block's height
	b.Height = block.Height

	b.CurrentDifficulty = block.Difficulty

	// persist the blockchain
	persistBlockchain(b)
	return block
}

// all blocks
func Blocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()
	var blocks []*Block

	// start newesthash and its prevhash and find block
	// if prevhash dont exist = genesis block break
	hashCursor := b.NewestHash

	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

// get the latest 6 blocks
func LatestBlock(b *blockchain, rw http.ResponseWriter) {
	var blocks []*Block
	// for _, v := range Blocks(b) {
	// 	h := fmt.Sprintf("%s", v.Hash[0:7]) + "..."
	// 	v.Hash = h
	// 	if len(v.PrevHash) > 7 {
	// 		ph := fmt.Sprintf("%s", v.PrevHash[0:7]) + "..."
	// 		v.PrevHash = ph
	// 	}
	// 	blocks = append(blocks, v)
	// }
	blocks = Blocks(b)
	if len(blocks) > 6 {
		blocks = blocks[0:6]
	}
	utils.HandleErr(json.NewEncoder(rw).Encode(blocks))
}

// persist the blockchain data
func persistBlockchain(b *blockchain) {
	// db.SaveCheckpoint(utils.ToBytes(b))
	dbStorage.SaveChain((utils.ToBytes(b)))
}

func Blockchain() *blockchain {
	// run only one time
	once.Do(func() {
		// initial blockchain struct
		b = &blockchain{Height: 0}

		// search for checkpoint on the db
		// checkpoint := db.Checkpoint()
		checkpoint := dbStorage.LoadChain()

		if checkpoint == nil {
			// if blockchain don't exist create block
			b.AddBlock("")
		} else {
			// restore data from db
			b.restore(checkpoint)
		}
	})
	// return type blockchain struct
	return b
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()

	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}

// recalculate difficulty of block by timestamp
func recalculateDifficulty(b *blockchain) int {
	// get all blocks
	allBlocks := Blocks(b)
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

func getDifficulty(b *blockchain) int {
	// if genesis block or not
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return recalculateDifficulty(b)
	} else {
		if b.CurrentDifficulty <= 5 {
			return b.CurrentDifficulty
		} else {
			return 5
		}
	}
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	dbStorage.DeleteAllBlocks()
	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) LockBlockchain() {
	b.m.Lock()
	defer b.m.Unlock()
}
