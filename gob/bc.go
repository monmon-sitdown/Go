package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	Nonce        int
	PreHash      [32]byte
	Timestamp    int64
	Transactions []*Transaction
}
type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}
type Transaction struct {
	senderAddr   string
	receiverAddr string
	value        float32
}

func NewBlock(nonce int, preHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.Timestamp = time.Now().UnixNano()
	b.Nonce = nonce
	b.PreHash = preHash
	b.Transactions = transactions
	return b
	/*return &Block{
		timestamp: time.Now().UnixNano(),
	}*/
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	//fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

func (b *Block) Print() {
	fmt.Printf("timestamp	%d\n", b.Timestamp)
	fmt.Printf("nonce	%d\n", b.Nonce)
	fmt.Printf("preHash	%x\n", b.PreHash)
	for _, t := range b.Transactions {
		t.Print()
	}
}

func NewBlockChain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) AddTransaction(sender, receiver string, value float32) {
	t := NewTransaction(sender, receiver, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CreateBlock(nonce int, preHash [32]byte) *Block {
	b := NewBlock(nonce, preHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain	%d %s \n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func NewTransaction(sender, receiver string, value float32) *Transaction {
	return &Transaction{sender, receiver, value}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddr   string  `json:"sender_addr"`
		ReceiverAddr string  `json:"receiver_addr"`
		Value        float32 `json:"value"`
	}{
		SenderAddr:   t.senderAddr,
		ReceiverAddr: t.receiverAddr,
		Value:        t.value,
	})
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender sddress	%s\n", t.senderAddr)
	fmt.Printf("receiver sddress	%s\n", t.receiverAddr)
	fmt.Printf("value	%.1f\n", t.value)
}

func init() {
	log.SetPrefix("bloc:")
}

func main() {
	bc := NewBlockChain()
	bc.Print()

	bc.AddTransaction("A", "B", 1.0)

	preHash := bc.LastBlock().Hash()
	bc.CreateBlock(5, preHash)
	bc.Print()

	bc.AddTransaction("C", "D", 2.0)
	bc.AddTransaction("X", "Y", 3.0)
	preHash = bc.LastBlock().Hash()
	bc.CreateBlock(2, preHash)
	bc.Print()

	log.Println("2")
}
