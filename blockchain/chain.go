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
	FindUserBlock(hash string) []byte
	FindStoreBlock(hash string) []byte
	FindBaljaguk(hash string) []byte
	LoadUserChain() []byte
	LoadStoreChain() []byte
	LoadBaljagukChain() []byte
	SaveUserBlock(hash string, data []byte)
	SaveStoreBlock(hash string, data []byte)
	SaveBaljagukBlock(hash string, data []byte)
	SaveUserChain(data []byte)
	SaveStoreChain(data []byte)
	SaveBaljagukChain(data []byte)
	DeleteAllBlocks()
}

// variable struct that play func only one time
var once sync.Once

var dbStorage storage = db.DB{}

func (b *blockchain) restore(data []byte) {
	// decoder := gob.NewDecoder(bytes.NewReader(data))
	// decoder.Decode(b)
	utils.FromBytes(b, data)
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()

	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}
