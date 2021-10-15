package main

import (
	"os"
	"strconv"

	"github.com/yoonhero/baljaguk_server/db"
	"github.com/yoonhero/baljaguk_server/rest"
	"github.com/yoonhero/baljaguk_server/utils"
)

func main() {

	// close db to protect db file data
	defer db.CloseSqlDB()
	db.InitPostgresDB()
	port := os.Getenv("PORT")

	if port == "" {
		rest.Start(4000)
	} else {
		sv, err := strconv.Atoi(port)
		utils.HandleErr(err)
		rest.Start(sv)
	}
}
