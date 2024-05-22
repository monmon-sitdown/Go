package block

import (
	"Blockchain/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE Blockchain"
	MINING_REWARD     = 1.0
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
	bcAddr          string
	port            uint16
}
type Transaction struct {
	senderAddr   string
	receiverAddr string
	value        float32
}
type TransactionRequest struct {
	SenderBlockchainAddress   *string  `json:"sender_blockchain_address"`
	ReceiverBlockchainAddress *string  `json:"receiver_blockchain_address"`
	SenderPublicKey           *string  `json:"sender_public_key"`
	Value                     *float32 `json:"value"`
	Signature                 *string  `json:"signature"`
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

func NewBlockChain(bcAddr string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.bcAddr = bcAddr
	bc.port = port
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) AddTransaction(sender, receiver string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, receiver, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		if bc.CalculateTotalAmount(sender) < value {
			log.Println("Errror: Not enough Money")
			return false
		}
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR : Verify Transaction")
	}
	return false
}

func (bc *Blockchain) CalculateTotalAmount(bcAddr string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.Transactions {
			value := t.value
			if bcAddr == t.receiverAddr {
				totalAmount += value
			}

			if bcAddr == t.senderAddr {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.senderAddr, t.receiverAddr, t.value))
	}
	return transactions
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.bcAddr, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	preHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, preHash)
	log.Println("action-mining, status=success")
	return true
}

func (bc *Blockchain) ValidProof(nonce int, preHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{nonce, preHash, 0, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))

	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	preHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, preHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
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

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil ||
		tr.ReceiverBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}
	return true
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender   string  `json:"sender_addr"`
		Receiver string  `json:"receiver_addr"`
		Value    float32 `json:"value"`
	}{
		Sender:   t.senderAddr,
		Receiver: t.receiverAddr,
		Value:    t.value,
	})
}
func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender address	%s\n", t.senderAddr)
	fmt.Printf("receiver address	%s\n", t.receiverAddr)
	fmt.Printf("value	%.1f\n", t.value)
}
