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
	"time"
)

type Skill struct {
}
type Course struct {
	Code        string
	Name        string
	CreditHours int
	Grade       string
}
type Project struct {
	name     string
	document string
	Course   Course
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
	Course      Course
	project     Project
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}
type Connected struct {
	Conn net.Conn
}

//var chainHead *Block
var globalData Data
var localData []Connected

//var chainHead *Block
var mutex = &sync.Mutex{}

//var globalData Data

func CalculateHash(inputBlock *Block) string {

	var temp string
	temp = inputBlock.Course.Code + inputBlock.project.name
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
		fmt.Print(" Course: ", chainHead.Course.Name)
		if (chainHead.Course != Course{}) {
			fmt.Print(" Course: ", chainHead.Course.Name)
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
		ln, err := net.Listen("tcp", "localhost:"+listeningAddress)
		if err != nil {
			log.Fatal(err, ln)
		}

	}
}
func MinerverifyBlock(conn net.Conn) {
	var recvdBlock *Block
	//fmt.Println("block: ", recvdBlock.Course.name)
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
var checkchan = make(chan string)

var NewChan = make(chan string)

//var RWChan = make(chan string)
var MinerChan = make(chan string)

//var WaitChan = make(chan string)

func readAdminData(conn net.Conn) {
	for {
		//var globe Data

		var globe Data
		gobEncoder := gob.NewDecoder(conn)
		//Stuck
		err1 := gobEncoder.Decode(&globe)
		//Stuck
		fmt.Println("In Admindata before: ", globe)
		if err1 != nil {
			//	log.Println(err1)
		}
		ListBlocks(globalData.ChainHead)
		fmt.Println("In read admin data after:")
		if Length(globe.ChainHead) < Length(globalData.ChainHead) {
			globe.ChainHead = globalData.ChainHead
		}
		if len(globe.ClientsSlice) < len(globalData.ClientsSlice) {
			globe.ClientsSlice = globalData.ClientsSlice
		}
		globalData = globe
		ListBlocks(globalData.ChainHead)
		<-Globechan
		//	<-checkchan

		// <-MinerChan
		//
		// //	<-RWChan
		// <-NewChan

	}
}

func ViewMinerData() {

	for {
		Globechan <- "hello"

		for i := 0; i < len(globalData.ClientsSlice); i++ {
			if globalData.ClientsSlice[i].Role == "miner" {
				fmt.Println("Miners connected to system:")
				fmt.Println(" Their address: ", globalData.ClientsSlice[i].ListeningAddress)
			}
		}
	}
}
func readBlockchain(conn net.Conn) {
	for {
		//var globe Data
		MinerChan <- "Start Reading"

		fmt.Println("In read blockchain before decode:")
		var chainhead *Block
		gobEncoder := gob.NewDecoder(conn)
		//Stuck
		err1 := gobEncoder.Decode(&chainhead)
		if err1 != nil {
			//	log.Println(err1, "Errorzz in readBlock")
		}
		if Length(chainhead) > Length(globalData.ChainHead) {
			globalData.ChainHead = chainhead
		}
		fmt.Println("Blockchain updated :: ", Length(globalData.ChainHead))
		ListBlocks(globalData.ChainHead)

		//	globalData = globe
		//	<-Globechan
	}
}
func UserSendBlock(minerAddress string, block *Block) {
	//Input from me

	//Dialing Miner
	conn, errs := net.Dial("tcp", ":"+minerAddress)
	if errs != nil {
		log.Fatal(errs)
	}
	conns := Connected{
		Conn: conn,
	}
	localData = append(localData, conns)
	fmt.Println("Sending Block CONTENT to be verified to miner")
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(block)
	if err != nil {
		//	log.Println(err)
	}
}
func UserSendCourse(minerAddress string, block Course) {
	conn, errs := net.Dial("tcp", ":"+minerAddress)
	if errs != nil {
		log.Fatal(errs)
	}
	conns := Connected{
		Conn: conn,
	}
	localData = append(localData, conns)
	fmt.Println("Sending Block CONTENT to be verified to miner")
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(block)
	if err != nil {
		//	log.Println(err)
	}
}
func UserReceiveBlock(conn net.Conn) {

	fmt.Println("Receiving Block CONTENT verified by miner")
	gobEncoder := gob.NewDecoder(conn)
	var chainHead *Block
	err := gobEncoder.Decode(chainHead)
	if err != nil {
		log.Println(err)
	}
	globalData.ChainHead = chainHead
}
func broadcastBlock() {
	NewChan <- "Hello"
	//	RWChan <- "Yoo"
	time.Sleep(5 * time.Second)
	for i := 0; i < len(localData); i++ {
		gobEncoder := gob.NewEncoder(localData[i].Conn)
		//fmt.Println("BroadCheck: ", localData[i])
		err1 := gobEncoder.Encode(globalData.ChainHead)
		fmt.Println("Broadcasting Blockchain:: ")
		if err1 != nil {
			log.Println(err1, "in broadcasting block")
		}

	}

}
func broadcastAdminData() {
	checkchan <- "hello"

	for i := 0; i < len(localData); i++ {
		gobEncoder := gob.NewEncoder(localData[i].Conn)
		//fmt.Println("BroadCheck: ", localData[i])
		err1 := gobEncoder.Encode(globalData)
		fmt.Println("Broadcasting StreamData:: ")
		if err1 != nil {
			//	log.Println(err)
		}

	}
	//	<-StepbyChan

}

func input() {
	for {
		fmt.Println("Enter Verifier port number from the list: ")
		fmt.Println()
		var minerAddress string
		fmt.Scanln(&minerAddress)

		fmt.Println("Enter 1 for Course or 2 for Project details to verify: ")
		var numb int
		fmt.Scanln(&numb)
		var cour Course
		//var proj Project
		block := &Block{}
		if numb == 1 {
			fmt.Println("Enter name for Course: ")
			var names string
			fmt.Scanln(&names)
			fmt.Println("Names: ", names)
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
				Code:        string(code),
				Name:        string(names),
				Grade:       string(grade),
				CreditHours: int(creditHours),
			}
			block.Course = cour

		} else if numb == 2 {

		}

		//minerAddress := "1200"
		//fmt.Println("In Course:: ", block.Course)
		conn, errs := net.Dial("tcp", ":"+minerAddress)
		if errs != nil {
			log.Fatal(errs)
		}
		//	UserSendBlock(minerAddress, block)
		UserSendCourse(minerAddress, cour)
		readAdminData(conn)
		//	mutex.Unlock()
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
	conns := Connected{
		Conn: conn,
	}
	localData = append(localData, conns)
	//The function below launches the server, uses different second argument
	//It then starts a routine for each connection request received
	//	role := "user"

	myPeer := Peer{
		ListeningAddress: string(myListeningAddress),
		Role:             "user",
	}
	go StartListening(myListeningAddress, "user")

	WriteString(conn, myPeer)
	//	go broadcastBlock()
	//go broadcastAdminData()
	//log.Println("Sending my listening address to Admin")

	go readAdminData(conn)

	go ViewMinerData()

	input()
	//	mutex.Lock()

	//	go readBlockchain(conn)

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
	// fmt.Println("block: ", recvdBlock.Course.name)
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
