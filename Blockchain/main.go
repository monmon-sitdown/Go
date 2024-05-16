package main

import (
	"Blockchain/wallet"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("bloc:")
}

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
}
