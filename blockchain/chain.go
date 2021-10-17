// persistence of block
// connected to DB to save data
// using bolt DB (bitcoin levelDB)
package blockchain

import (
	"encoding/json"
	"net/http"

	"github.com/yoonhero/baljaguk_server/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 1
	allowedRange       int = 2
)

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

func Status(b *baljagukBlockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()

	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}
