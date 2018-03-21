package main

import (
	"encoding/json"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"log"
	"io"
)

type Block struct {
	Index int // 是这个块在整个链中的位置
	Timestamp string // 显而易见就是块生成时的时间戳
	BPM       int
	Hash      string
	PrevHash  string
}

var Blockchain []Block // 本地的BlockChain

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// BPM相当于数据
func generateBlock(oldBlock Block, BPM int) (Block, error){
	var newBlock Block

	t := time.Now()
	newBlock.Index = oldBlock.Index+1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

// 校验块new block是否
func isBlockValid(nextBlock, prevBlock Block) bool {
	if prevBlock.Index+1 != nextBlock.Index {
		return false
	}
	if prevBlock.Hash != nextBlock.PrevHash {
		return false
	}
	if calculateHash(nextBlock) != nextBlock.Hash {
		return false
	}
	return true
}

// 将获取到的Blockchain进行替换
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func makeMuxRouter() http.Handler{
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

type Message struct {
	BPM int
}