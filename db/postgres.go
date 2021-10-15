package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/yoonhero/baljaguk_server/utils"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "randompassword"
	dbname   = "baljaguk"
)

var sqlDB *sql.DB

//////////////////////////////// Basic Functions //////////////////////////////////
// database url to connect
func dsn() string {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		return databaseURL
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func createBlocksTable() {
	stmt, err := sqlDB.Prepare("CREATE TABLE IF NOT EXISTS UserBlocks (Hash varchar(111) NOT NULL, Data bytea NOT NULL)")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)

	stmt, err = sqlDB.Prepare("CREATE TABLE IF NOT EXISTS StoreBlocks (Hash varchar(111) NOT NULL, Data bytea NOT NULL)")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)

	stmt, err = sqlDB.Prepare("CREATE TABLE IF NOT EXISTS BaljagukBlocks (Hash varchar(111) NOT NULL, Data bytea NOT NULL)")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
}

func createCheckpointTable() {
	stmt, err := sqlDB.Prepare("CREATE TABLE IF NOT EXISTS UserCheckpoint (Data bytea NOT NULL)")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)

	stmt, err = sqlDB.Prepare("CREATE TABLE IF NOT EXISTS StoreCheckpoint (Data bytea NOT NULL)")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)

	stmt, err = sqlDB.Prepare("CREATE TABLE IF NOT EXISTS BaljagukCheckpoint (Data bytea NOT NULL)")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
}

////////////////////////////    Basic Setting of Database     //////////////////////////////////

func CloseSqlDB() {
	sqlDB.Close()
}

func InitPostgresDB() {
	if sqlDB == nil {
		db, err := sql.Open("postgres", dsn())
		utils.HandleErr(err)

		err = db.Ping()
		utils.HandleErr(err)

		log.Printf("Connected to DB %s successfully\n", dbname)

		sqlDB = db

		createBlocksTable()
		createCheckpointTable()
	}
}

//////////////////////////////////    Main Functions     //////////////////////////////////////////////

//////////////////////////////// User ////////////////////////////////
// empty chain table
func emptyUserChainTable() {
	stmt, err := sqlDB.Prepare("DROP TABLE UserCheckpoint")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
	createCheckpointTable()
}

// save chain
func saveUserChainInSQL(data []byte) {
	emptyUserChainTable()

	_, err := sqlDB.Exec("INSERT INTO UserCheckpoint(Data) values($1)", data)
	utils.HandleErr(err)
}

// load chain
func loadUserChainInSQL() []byte {
	var data []byte

	rows, err := sqlDB.Query("SELECT Data FROM UserCheckpoint")
	utils.HandleErr(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&data)
		utils.HandleErr(err)
	}

	return data
}

// save block data
func saveUserBlockInSQL(hash string, data []byte) {
	// update database
	_, err := sqlDB.Exec("INSERT INTO UserBlocks(Hash, Data) values($1, $2)", hash, data)
	utils.HandleErr(err)
}

func findUserBlockInSQL(hash string) []byte {
	var data []byte

	err := sqlDB.QueryRow("SELECT Data FROM UserBlocks WHERE Hash = $1", hash).Scan(&data)
	utils.HandleErr(err)

	return data
}

//////////////////////////////// Store ////////////////////////////////
// empty chain table
func emptyStoreChainTable() {
	stmt, err := sqlDB.Prepare("DROP TABLE StoreCheckpoint")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
	createCheckpointTable()
}

// save chain
func saveStoreChainInSQL(data []byte) {
	emptyStoreChainTable()

	_, err := sqlDB.Exec("INSERT INTO StoreCheckpoint(Data) values($1)", data)
	utils.HandleErr(err)
}

// load chain
func loadStoreChainInSQL() []byte {
	var data []byte

	rows, err := sqlDB.Query("SELECT Data FROM StoreCheckpoint")
	utils.HandleErr(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&data)
		utils.HandleErr(err)
	}

	return data
}

func saveStoreBlockInSQL(hash string, data []byte) {
	// update database
	_, err := sqlDB.Exec("INSERT INTO StoreBlocks(Hash, Data) values($1, $2)", hash, data)
	utils.HandleErr(err)
}

func findStoreBlockInSQL(hash string) []byte {
	var data []byte

	err := sqlDB.QueryRow("SELECT Data FROM StoreBlocks WHERE Hash = $1", hash).Scan(&data)
	utils.HandleErr(err)

	return data
}

//////////////////////////////// Baljaguk ////////////////////////////////
// empty chain table
func emptyBaljagukChainTable() {
	stmt, err := sqlDB.Prepare("DROP TABLE BaljagukCheckpoint")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
	createCheckpointTable()
}

// save chain
func saveBaljagukChainInSQL(data []byte) {
	emptyBaljagukChainTable()

	_, err := sqlDB.Exec("INSERT INTO BaljagukCheckpoint(Data) values($1)", data)
	utils.HandleErr(err)
}

// load chain
func loadBaljagukChainInSQL() []byte {
	var data []byte

	rows, err := sqlDB.Query("SELECT Data FROM BaljagukCheckpoint")
	utils.HandleErr(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&data)
		utils.HandleErr(err)
	}

	return data
}

func saveBaljagukBlockInSQL(hash string, data []byte) {
	// update database
	_, err := sqlDB.Exec("INSERT INTO BaljagukBlocks(Hash, Data) values($1, $2)", hash, data)
	utils.HandleErr(err)
}

func findBaljagukBlockInSQL(hash string) []byte {
	var data []byte

	err := sqlDB.QueryRow("SELECT Data FROM BaljagukBlocks WHERE Hash = $1", hash).Scan(&data)
	utils.HandleErr(err)

	return data
}

//////////////////////////////// Else //////////////////////////////
func emptyUserBlocksInSQL() {
	stmt, err := sqlDB.Prepare("DROP TABLE UserBlocks")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
	createBlocksTable()
}

func emptyStoreBlocksInSQL() {
	stmt, err := sqlDB.Prepare("DROP TABLE StoreBlocks")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
	createBlocksTable()
}

func emptyBaljagukBlocksInSQL() {
	stmt, err := sqlDB.Prepare("DROP TABLE BaljagukBlocks")
	utils.HandleErr(err)

	_, err = stmt.Exec()
	utils.HandleErr(err)
	createBlocksTable()
}
