package main

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
)

type Skill struct {
}
type Course struct {
	code        string
	name        string
	creditHours int
	grade       string
}
type Project struct {
	name     string
	document string
	course   Course
}

type Block struct {
	course      Course
	project     Project
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

var chainHead *Block
var mutex = &sync.Mutex{}

func CalculateHash(inputBlock *Block) string {

	var temp string
	temp = inputBlock.course.code + inputBlock.project.name
	h := sha256.New()
	h.Write([]byte(temp))
	sum := hex.EncodeToString(h.Sum(nil))

	// sum := sha256.Sum256([]byte(temp))

	return sum
}
func InsertOnlyBlock(newBlock *Block, chainHead *Block) *Block {
	newBlock.CurrentHash = CalculateHash(newBlock)

	if chainHead == nil {
		chainHead = newBlock
		fmt.Println("Block Inserted")
		return chainHead
	}
	newBlock.PrevPointer = chainHead
	newBlock.PrevHash = chainHead.CurrentHash

	fmt.Println("Block Course and Project Inserted")
	return newBlock

}

func StartListening(listeningAddress string, node string) {
	//var chainHead *Block
	if node == "user" {

	} else if node == "miner" {
		ln, err := net.Listen("tcp", listeningAddress)
		if err != nil {
			log.Fatal(err, ln)
		}

	}
}
func MinerverifyBlock(conn net.Conn) {
	var recvdBlock *Block
	//fmt.Println("block: ", recvdBlock.course.name)
	dec2 := gob.NewDecoder(conn)
	err2 := dec2.Decode(&recvdBlock)
	if err2 != nil {
		//handle error
	} else {
		fmt.Println("Block Verified")
		InsertOnlyBlock(recvdBlock, chainHead)
	}
}
func main() {

	conn, err := net.Dial("tcp", "localhost:1403")
	if err != nil {
		//handle error
	}

	address := "2000"
	fmt.Println("Coursecc")
	dec := gob.NewEncoder(conn)
	err = dec.Encode(&address)
	if err != nil {
		//handle error
	}

	//Decode
	MinerverifyBlock(conn)
	// var recvdBlock Block
	// fmt.Println("block: ", recvdBlock.course.name)
	// dec2 := gob.NewDecoder(conn)
	// err2 := dec2.Decode(&recvdBlock)
	// if err2 != nil {
	// 	//handle error
	// }
	//go StartListening(":4502", "miner")

	// var chainHead *Block
	// gobEncoder := gob.NewEncoder(conn)
	// err := gobEncoder.Encode(chainhead)
	// if err != nil {
	// 	//	log.Println(err)
	// }

	// fmt.Println(recvdBlock.CurrentHash)
	select {}
}
