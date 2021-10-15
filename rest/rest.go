package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yoonhero/baljaguk_server/blockchain"
	"github.com/yoonhero/baljaguk_server/utils"
	"github.com/yoonhero/baljaguk_server/wallet"
)

// variable post string
var port string

// new type URL
type url string

// type URL's interface
func (u url) MarshalText() ([]byte, error) {
	// var url is http://localhost + port + URL
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

//////////////////////////////// Http Decode Structure //////////////////////////////////
type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addUserBlockBody struct {
	Address     string `json:"address"`
	PrivateKey  string `json:"privateKey"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

type addStoreBlockBody struct {
	Address     string `json:"address"`
	PrivateKey  string `json:"privateKey"`
	PhoneNumber string `json:"phoneNumber"`
}

type addBaljagukBlockBody struct {
	StoreHash string `json:"storehash"`
	UserHash  string `json:"userhash"`

	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type walletPayload struct {
	Key string `json:"key"`
}

type createKeyAddressPayload struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

type InfoMining struct {
	Block      *blockchain.Block `json:"block"`
	Hash       string            `json:"hash"`
	Difficulty int               `json:"difficulty"`
}

// when url is "/"
func documentation(rw http.ResponseWriter, r *http.Request) {

	// []URLDescription struct slice
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the blockchain (user, store, baljaguk)",
		},
		{
			URL:         url("/userblocks"),
			Method:      "POST",
			Description: "Add A User Block",
			Payload:     "{'username':'','phone_number':'','address':''}",
		},
		{
			URL:         url("/storeblocks"),
			Method:      "POST",
			Description: "Add A Store Block",
			Payload:     "{'storename':'','phone_number':'','address':''}",
		},
		{
			URL:         url("/baljaguk?username='',storename=''"),
			Method:      "Get",
			Description: "See A Baljaguk TimeLine Block",
			Payload:     "query: username, storename",
		},
		{
			URL:         url("/baljaguks"),
			Method:      "POST",
			Description: "Add A Baljaguk Block",
			Payload:     "{'store_id':'','user_id':'','timestamp':''}",
		},
		{
			URL:         url("/createkey"),
			Method:      "GET",
			Description: "Make a Random Private Key and Public Key for User Block",
		},
	}
	// add content json type
	// rw.Header().Add("Content-Type", "application/json")

	// json.NewEncoder(rw).Encode(data)
	// is same
	// b, err := json.Marshal(data)
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data)
}

// when get or post url /userblock
func userBlocks(rw http.ResponseWriter, r *http.Request) {
	var addUserBlockBody addUserBlockBody
	json.NewDecoder(r.Body).Decode(&addUserBlockBody)
	// {"message":"myblockdata"}

	// // new variable struct AddBlockBody
	// var addBlockBody addBlockBody

	// // send pointers and set variable a posted data
	// utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))

	if addUserBlockBody.Address == "" && addUserBlockBody.PrivateKey == "" && addUserBlockBody.PhoneNumber == "" && addUserBlockBody.Email == "" {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	// add block whose data is addBlockBody.Message
	blockchain.UserBlockchain().AddUserBlock(addUserBlockBody.Address, addUserBlockBody.PrivateKey, addUserBlockBody.PhoneNumber, addUserBlockBody.Email)

	// p2p.BroadcastNewBlock(newBlock)

	// send a 201 sign
	rw.WriteHeader(http.StatusCreated)
}

// when get or post url /userblock
func storeBlocks(rw http.ResponseWriter, r *http.Request) {
	var addStoreBlockBody addStoreBlockBody
	json.NewDecoder(r.Body).Decode(&addStoreBlockBody)
	// {"message":"myblockdata"}

	// // new variable struct AddBlockBody
	// var addBlockBody addBlockBody

	// // send pointers and set variable a posted data
	// utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))

	if addStoreBlockBody.Address == "" && addStoreBlockBody.PrivateKey == "" && addStoreBlockBody.PhoneNumber == "" {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	// add block whose data is addBlockBody.Message
	blockchain.StoreBlockchain().AddStoreBlock(addStoreBlockBody.Address, addStoreBlockBody.PrivateKey, addStoreBlockBody.PhoneNumber)

	// p2p.BroadcastNewBlock(newBlock)

	// send a 201 sign
	rw.WriteHeader(http.StatusCreated)

}

// when get or post url /userblock
func baljagukBlocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// when GET
	case "GET":
		// recognize that this content is json
		// rw.Header().Add("Content-Type", "application/json")

		// send all blocks
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.BaljagukBlocks(blockchain.BaljagukBlockchain())))

		// when POST
	case "POST":
		var addBaljagukBlockBody addBaljagukBlockBody
		json.NewDecoder(r.Body).Decode(&addBaljagukBlockBody)
		// {"message":"myblockdata"}

		// // new variable struct AddBlockBody
		// var addBlockBody addBlockBody

		// // send pointers and set variable a posted data
		// utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))

		// add block whose dat`a is addBlockBody.Message
		if addBaljagukBlockBody.UserHash == "" && addBaljagukBlockBody.StoreHash == "" && addBaljagukBlockBody.Latitude == "" && addBaljagukBlockBody.Longitude == "" {
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		blockchain.BaljagukBlockchain().AddBaljagukBlock(addBaljagukBlockBody.StoreHash, addBaljagukBlockBody.UserHash, addBaljagukBlockBody.Latitude, addBaljagukBlockBody.Longitude)

		// p2p.BroadcastNewBlock(newBlock)

		// send a 201 sign
		rw.WriteHeader(http.StatusCreated)
	}

}

func findUser(rw http.ResponseWriter, r *http.Request) {
	// get mux var from http.Request
	// shape looks like
	// map[id:1]
	vars := mux.Vars(r)

	// get only id
	// id := vars["height"]

	// strconv.Atoi convert string to int
	hash := vars["hash"]

	// handle err
	// utils.HandleErr(err)

	// FindBlock by id
	block, err := blockchain.FindUserBlock(hash)

	encoder := json.NewEncoder(rw)

	// if err founded
	if err == blockchain.ErrNotFound {
		// format err to string
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		// send the block
		encoder.Encode(block)
	}

}

func findStore(rw http.ResponseWriter, r *http.Request) {
	// get mux var from http.Request
	// shape looks like
	// map[id:1]
	vars := mux.Vars(r)

	// get only id
	// id := vars["height"]

	// strconv.Atoi convert string to int
	hash := vars["hash"]

	// FindBlock by id
	block, err := blockchain.FindStoreBlock(hash)

	encoder := json.NewEncoder(rw)

	// if err founded
	if err == blockchain.ErrNotFound {
		// format err to string
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		// send the block
		encoder.Encode(block)
	}

}

// func add json content type
func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	// make a type of http.Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// add content json type
		utils.AllowConnection(rw)
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleWare(next http.Handler) http.Handler {
	// make a type of http.Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// add content json type
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.BaljagukBlockchain(), rw)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	var payload walletPayload
	json.NewDecoder(r.Body).Decode(&payload)
	// json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
	bytes, err := hex.DecodeString(payload.Key)
	utils.HandleErr(err)
	json.NewEncoder(rw).Encode(wallet.RestApiWallet(bytes))
	rw.WriteHeader(http.StatusOK)
}

func createKey(rw http.ResponseWriter, r *http.Request) {
	address, key := wallet.RestApiCreatePrivKey()
	utils.HandleErr(json.NewEncoder(rw).Encode(createKeyAddressPayload{address, fmt.Sprintf("%x", key)}))
}

// func peers(rw http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "POST":
// 		var payload addPeerPayload
// 		json.NewDecoder(r.Body).Decode(&payload)
// 		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
// 		rw.WriteHeader(http.StatusOK)
// 	case "GET":
// 		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
// 	}
// }

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()

	// add json content type
	router.Use(jsonContentTypeMiddleWare, loggerMiddleWare)

	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")

	router.HandleFunc("/userblock", userBlocks).Methods("POST")
	router.HandleFunc("/storeblock", storeBlocks).Methods("POST")
	router.HandleFunc("/baljaguks", baljagukBlocks)

	router.HandleFunc("/store/{hash:[a-f0-9]+}", findStore).Methods("GET")
	router.HandleFunc("/user/{hash:[a-f0-9]+}", findUser).Methods("GET")

	router.HandleFunc("/wallet", myWallet).Methods("POST")
	router.HandleFunc("/createkey", createKey).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)

	// print if err exist
	log.Fatal(http.ListenAndServe(port, router))
}
