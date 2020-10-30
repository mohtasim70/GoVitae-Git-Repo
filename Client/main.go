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
}
type Data struct {
	MinerList    []Peer
	ClientsSlice []Peer
	ChainHead    *Block
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
func ListBlocks(chainHead *Block) {

	for chainHead != nil {
		fmt.Print("Block-- ")
		fmt.Print(" Current Hash: ", chainHead.CurrentHash)
		if chainHead.PrevHash == "" {
			fmt.Print(" Previous Hash: ", "Null")
		} else {
			fmt.Print(" Previous Hash: ", chainHead.PrevHash)
		}
		if (chainHead.course != Course{}) {
			fmt.Print(" Course: ", chainHead.course.name)
		}
		if (chainHead.project != Project{}) {
			fmt.Print(" Project: ", chainHead.project.name)
		}

		fmt.Print(" -> ")
		chainHead = chainHead.PrevPointer

	}
	fmt.Println()

}
func Length(chainHead *Block) int {
	sum := 0
	for chainHead != nil {

		chainHead = chainHead.PrevPointer
		sum++
	}
	return sum

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
		InsertOnlyBlock(recvdBlock, globalData.ChainHead)
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

var Globechan = make(chan string)

func readAdminData(conn net.Conn) {
	for {
		//var globe Data
		var globe Data
		gobEncoder := gob.NewDecoder(conn)
		//Stuck
		err1 := gobEncoder.Decode(&globe)
		//Stuck
		//	fmt.Println("In Admindata: ", globe)
		if err1 != nil {
			log.Println(err1)
		}
		fmt.Println("In read admin data:")
		if Length(globe.ChainHead) < Length(globalData.ChainHead) {
			globe.ChainHead = globalData.ChainHead
		}
		globalData = globe
		<-Globechan
	}
}

func ViewMinerData() {

	for {
		Globechan <- "hello"

		for i := 0; i < len(globalData.ClientsSlice); i++ {
			if globalData.ClientsSlice[i].Role == "miner" {
				fmt.Println("Miners connected to system:")
				fmt.Print(" Their address: ", globalData.ClientsSlice[i].ListeningAddress)
			}
		}
	}
}
func UserSendBlock(minerAddress string, block *Block) {
	//Input from me

	//Dialing Miner
	conn, errs := net.Dial("tcp", ":"+minerAddress)
	if errs != nil {
		log.Fatal(errs)
	}
	fmt.Println("Sending Block CONTENT to be verified to miner")
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(block)
	if err != nil {
		//	log.Println(err)
	}
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
	//log.Println("Sending my listening address to Admin")

	go readAdminData(conn)

	go ViewMinerData()

	fmt.Println("Enter Verifier port number from the list: ")
	var minerAddress string
	fmt.Scanln(&minerAddress)

	fmt.Println("Enter 1 for Course or 2 for Project details to verify: ")
	var numb int
	fmt.Scanln(&numb)
	var cour Course
	//var proj Project
	if numb == 1 {
		fmt.Println("Enter name for Course: ")
		var names string
		fmt.Scanln(&names)
		fmt.Println("Enter code for Course: ")
		var code string
		fmt.Scanln(&code)
		fmt.Println("Enter grade for Course: ")
		var grade string
		fmt.Scanln(&grade)
		fmt.Println("Enter credit hours for Course: ")
		var creditHours int
		fmt.Scanln(&creditHours)
		cour = Course{
			code:        code,
			name:        names,
			grade:       grade,
			creditHours: creditHours,
		}

	} else if numb == 2 {

	}

	//minerAddress := "1200"
	block := &Block{
		course: cour,
	}

	UserSendBlock(minerAddress, block)

	//DIaling

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
