package main

import (
	"Blockchain/block"
	"Blockchain/wallet"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("bloc:")
}

func main() {
	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	//wallet side
	t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)
	//blockchain side
	blockchain := block.NewBlockChain(walletM.BlockchainAddress())
	isAdded := blockchain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, walletA.PublicKey(), t.GenerateSignature())
	fmt.Println("isAdded? -- ", isAdded)
}
