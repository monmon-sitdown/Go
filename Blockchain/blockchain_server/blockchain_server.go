package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func HelloWorld(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "HelloWorld")
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", HelloWorld)
	address := "0.0.0.0:" + strconv.Itoa(int(bcs.port))
	fmt.Println("Server is running on", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
