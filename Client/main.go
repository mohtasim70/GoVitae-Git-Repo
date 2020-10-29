package main

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
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
type Peer struct {
	ListeningAddress string
	Role             string //1 for user 0 for miner
	Conn             net.Conn
}
type Data struct {
	minerList    []Peer
	clientsSlice []Peer
	chainHead    *Block
}

type Block struct {
	course      Course
	project     Project
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

//var chainHead *Block
var globalData Data
var mutex = &sync.Mutex{}

//var globalData Data

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
	if node == "admin" {

	} else if node == "user" {
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
		InsertOnlyBlock(recvdBlock, globalData.chainHead)
	}
}
func WriteString(conn net.Conn, details Peer) {
	fmt.Println("Peer: ", details)
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(details)

	if err != nil {
		//	log.Println(err)
	}
}

func readAdminData(conn net.Conn) {
	var globe Data
	gobEncoder := gob.NewDecoder(conn)
	err1 := gobEncoder.Decode(&globe)
	if err1 != nil {
		//	log.Println(err)
	}
	fmt.Println("In read admin data:")
	globalData = globe
}

func main() {

	satoshiAddress := os.Args[1]
	myListeningAddress := os.Args[2]
	log.Println(satoshiAddress, myListeningAddress)

	conn, err := net.Dial("tcp", ":"+satoshiAddress)
	if err != nil {
		log.Fatal(err)
	}
	//The function below launches the server, uses different second argument
	//It then starts a routine for each connection request received
	//	role := "user"

	myPeer := Peer{
		ListeningAddress: string(myListeningAddress),
		Role:             "user",
	}
	//go StartListening(myListeningAddress, "user")

	WriteString(conn, myPeer)
	log.Println("Sending my listening address to Admin")

	go readAdminData(conn)
	//Satoshi is there waiting for our address, it stores it somehow
	// encoder := gob.NewEncoder(conn)
	// p := P{
	// 	Appleeeee: myListeningAddress,
	// 	Cursorrrr: "userr",
	// }
	// encoder.Encode(p)

	//Decode
	//MinerverifyBlock(conn)
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
