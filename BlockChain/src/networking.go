package main

import (
	"time"
	//"log"
	"github.com/davecgh/go-spew/spew"
	//"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"io"
	"bufio"
	"strconv"
	"encoding/json"
	"sync"
)


// Block represents each 'item' in the blockchain
// type Block struct


// Blockchain is a series of validated Blocks

//var Blockchain []Block

// func calculateHash(block Block) string

// func generateBlock(oldBlock Block, BPM int) (Block, error)


// make sure block is valid by checking index, and comparing the hash of the previous block
// func isBlockValid(newBlock, oldBlock Block) bool

// func replaceChain(newBlock Block)

// bcServer handles incoming concurrent Blocks
var bcServer chan []Block

var mutex = &sync.Mutex{}

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal(err)
	//}

	bcServer = make(chan []Block)

	// create genesis block
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	spew.Dump(genesisBlock) // ???
	Blockchain = append(Blockchain, genesisBlock)

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept() // 容许一个TCP服务建立多个连接
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

	defer conn.Close()

	io.WriteString(conn, "Enter a new BPM:")

	scanner := bufio.NewScanner(conn)
	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v is not a number %v", scanner.Text(), err)
				continue
			}
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], bpm)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

			bcServer <- Blockchain
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	// simulate receiving broadcast
	go func() {
		for {
			time.Sleep(30 * time.Second)
			mutex.Lock()
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(output))
		}
	}()

	for range bcServer {
		spew.Dump(Blockchain)
	}
}





