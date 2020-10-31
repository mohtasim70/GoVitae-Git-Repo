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
var mutex = &sync.Mutex{}

//var globalData Data

func CalculateHash(inputBlock *Block) string {

	var temp string
	temp = inputBlock.Course.Name + inputBlock.project.name
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
func InsertCourse(course Course, chainHead *Block) *Block {
	newBlock := &Block{
		//Hash here
		Course: course,
	}
	newBlock.CurrentHash = CalculateHash(newBlock)

	if chainHead == nil {
		chainHead = newBlock
		fmt.Println("Block Inserted")
		return chainHead
	}
	newBlock.PrevPointer = chainHead
	newBlock.PrevHash = chainHead.CurrentHash

	fmt.Println("Project Block Inserted")
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

	} else if node == "miner" {
		ln, err := net.Listen("tcp", ":"+listeningAddress)
		if err != nil {
			log.Fatal(err, ln)
		}
		//	ClientsSlice := make([]Peer, 10)
		//	addchan := make(chan Peer)
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			// newClient := Peer{
			// 	Conn: conn,
			// }
			// ClientsSlice = append(ClientsSlice, newClient)
			// go broadcastBlockchaintoPeer(conn)
			conns := Connected{
				Conn: conn,
			}
			localData = append(localData, conns)
			// go receiveBlockchainfromPeer(conn)

			go MinerverifyBlock(conn)

		}
	}
}
func MinerverifyBlock(conn net.Conn) {
	var recvdBlock Course
	//fmt.Println("block: ", recvdBlock.Course.name)
	dec2 := gob.NewDecoder(conn)
	err2 := dec2.Decode(&recvdBlock)
	if err2 != nil {
		//handle error
		fmt.Println("err")
	} else {
		UpdateChan <- "start mining"
		mutex.Lock()
		//	fmt.Println("Checkinf Course: ", recvdBlock.Course)
		fmt.Println("Verify Data: Yes or No: ")
		var choice string
		fmt.Scanln(&choice)
		if choice == "Yes" {
			fmt.Println("Block Verified")

			globalData.ChainHead = InsertCourse(recvdBlock, globalData.ChainHead)
			broadcastBlock()
		} else {
			fmt.Println("Block Not Verified")
		}
		mutex.Unlock()
		fmt.Println("Length of blockchain", Length(globalData.ChainHead))
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

var UpdateChan = make(chan string)

var NewChan = make(chan string)

var MinerChan = make(chan string)

//var RW2Chan = make(chan string)

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
			//		log.Println(err1, "In readAdmindtat")
		}
		//	fmt.Println("In read admin data:")
		if Length(globe.ChainHead) < Length(globalData.ChainHead) {
			globe.ChainHead = globalData.ChainHead
		}
		globalData = globe
		<-MinerChan
		<-UpdateChan
		<-NewChan

	}
}
func ViewMinerData() {
	for i := 0; i < len(globalData.ClientsSlice); i++ {
		if globalData.ClientsSlice[i].Role == "miner" {
			fmt.Println("Miners connected to system:")
			fmt.Println(" Their address: ", globalData.ClientsSlice[i].ListeningAddress)
		}
	}
}
func broadcastBlock() {
	NewChan <- "start mining"
	//RW2Chan <- "Gee"
	time.Sleep(8 * time.Second)
	for i := 0; i < len(localData); i++ {
		gobEncoder := gob.NewEncoder(localData[i].Conn)
		//fmt.Println("BroadCheck: ", localData[i])
		err1 := gobEncoder.Encode(globalData.ChainHead)
		fmt.Println("Broadcasting Blockchain:: ")
		if err1 != nil {
			//	log.Println(err)
		}

	}

}
func readBlockchain(conn net.Conn) {
	for {
		//var globe Data
		MinerChan <- "Start Reading"

		//	fmt.Println("In read blockchain before decode:")
		var chainhead *Block
		gobEncoder := gob.NewDecoder(conn)
		//Stuck
		err1 := gobEncoder.Decode(&chainhead)
		if err1 != nil {
			//		log.Println(err1, "Error in readBlockchain")
		}

		if Length(chainhead) > Length(globalData.ChainHead) {
			globalData.ChainHead = chainhead
		}
		fmt.Println("Length blockchain:: ", Length(globalData.ChainHead))
		//	globalData = globe
		//	<-Globechan
		//	<-RW2Chan
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
		Role:             "miner",
	}
	go StartListening(myListeningAddress, "miner")

	WriteString(conn, myPeer)
	//	log.Println("I Verifier Sending my listening address to Admin")

	go readAdminData(conn)

	go broadcastBlock()

	go readBlockchain(conn)

	select {}
}
