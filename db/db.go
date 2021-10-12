package db

type DB struct {
}

func (DB) FindUserBlock(hash string) []byte {
	return findUserBlockInSQL(hash)
}

func (DB) FindStoreBlock(hash string) []byte {
	return findStoreBlockInSQL(hash)
}

func (DB) FindBaljaguk(hash string) []byte {
	return findBaljagukBlockInSQL(hash)
}

func (DB) LoadUserChain() []byte {
	return loadUserChainInSQL()
}

func (DB) LoadStoreChain() []byte {
	return loadStoreChainInSQL()
}

func (DB) LoadBaljagukChain() []byte {
	return loadBaljagukChainInSQL()
}

func (DB) SaveUserChain(data []byte) {
	saveUserChainInSQL(data)
}

func (DB) SaveStoreChain(data []byte) {
	saveStoreChainInSQL(data)
}

func (DB) SaveBaljagukChain(data []byte) {
	saveBaljagukChainInSQL(data)
}

func (DB) SaveUserBlock(hash string, data []byte) {
	saveUserBlockInSQL(hash, data)
}

func (DB) SaveStoreBlock(hash string, data []byte) {
	saveStoreBlockInSQL(hash, data)
}

func (DB) SaveBaljagukBlock(hash string, data []byte) {
	saveBaljagukBlockInSQL(hash, data)
}

func (DB) DeleteAllBlocks() {
	emptyUserBlocksInSQL()
	emptyStoreBlocksInSQL()
	emptyBaljagukBlocksInSQL()
}
