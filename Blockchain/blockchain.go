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
	course   Course
}
type Peer struct {
	ListeningAddress string
	Role             string //1 for user 0 for miner
	//	Conn             net.Conn
}
type Connected struct {
	Conn net.Conn
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

//var chainHead *Block
var localData []Connected
var globalData Data
var mutex = &sync.Mutex{}

//256bit
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
		fmt.Println("Genesis Block Inserted")
		return chainHead
	}
	newBlock.PrevPointer = chainHead
	newBlock.PrevHash = chainHead.CurrentHash

	fmt.Println("Later Block  Inserted")
	return newBlock

}
func InsertBlock(course Course, project Project, chainHead *Block) *Block {
	newBlock := &Block{
		//Hash here
		Course:  course,
		project: project,
	}
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
func InsertProject(project Project, chainHead *Block) *Block {
	newBlock := &Block{
		//Hash here
		project: project,
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

func ChangeCourse(oldCourse Course, newCourse Course, chainHead *Block) {
	present := false
	for chainHead != nil {
		if chainHead.Course == oldCourse {
			chainHead.Course = newCourse
			present = true
			break
		}

		//fmt.Printf("%v ->", chainHead.transactions)
		chainHead = chainHead.PrevPointer
	}
	if present == false {
		fmt.Println("Input Course not found")
		return
	}
	fmt.Println("Block Course Changed")

	chainHead.CurrentHash = CalculateHash(chainHead)

}

func ChangeProject(oldProject Project, newProject Project, chainHead *Block) {
	present := false
	for chainHead != nil {
		if chainHead.project == oldProject {
			chainHead.project = newProject
			present = true
			break
		}

		//fmt.Printf("%v ->", chainHead.transactions)
		chainHead = chainHead.PrevPointer
	}
	if present == false {
		fmt.Println("Input Course not found")
		return
	}
	fmt.Println("Block Course Changed")

	chainHead.CurrentHash = CalculateHash(chainHead)

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

func VerifyChain(chainHead *Block) { //What to do?
	for chainHead != nil {
		if chainHead.PrevPointer != nil {
			if chainHead.PrevHash != chainHead.PrevPointer.CurrentHash {
				fmt.Println("Blockchain Compromised")
				return
			}
		}

		chainHead = chainHead.PrevPointer
	}
	fmt.Println("Blockchain Verified")
}
func sendBlockchain(c net.Conn, chainHead *Block) {

	log.Println("A client has connected",

		c.RemoteAddr())
	gobEncoder := gob.NewEncoder(c)
	err := gobEncoder.Encode(chainHead)
	if err != nil {

		log.Println(err)

	}
}

//Read

// func WriteData(conn net.Conn, blockchan chan *Block) {
//
// 	firstCourse := Course{code: "CS50", name: "AI", creditHours: 3, grade: "A+"}
// 	block := &Block{
// 		//Hash here
// 		course: firstCourse,
// 	}
// 	blockchan <- block
// 	gobEncoder := gob.NewEncoder(conn)
// 	err1 := gobEncoder.Encode(block)
// 	if err1 != nil {
// 		//	log.Println(err)
// 	}
//
// }

var addchan = make(chan Peer)
var globe Data
var stopchan = make(chan string)

//var clientsSlice []Verifier
func broadcastAdminData() {

	for i := 0; i < len(localData); i++ {
		gobEncoder := gob.NewEncoder(localData[i].Conn)
		//fmt.Println("BroadCheck: ", localData[i])
		err1 := gobEncoder.Encode(globalData)
		fmt.Println("Broadcasting StreamData:: ")
		if err1 != nil {
			//	log.Println(err)
		}

	}
	<-StepbyChan

}

func StoreClient(conn net.Conn) {

	log.Println("A client has connected",
		conn.RemoteAddr())

	newClient := Peer{}
	//Block or Blockchain from Verifier?
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&newClient)
	if err != nil {
		//handle error
		log.Println(err)
	}

	fmt.Printf("Received : %+v", newClient)

	//	newClient.Conn = conn

	//	fmt.Println("Slice:", globalData.clientsSlice[0].ListeningAddress)
	addchan <- newClient
	//	<-stopchan

	// dec := gob.NewDecoder(conn)
	// p := P{}
	// dec.Decode(&p)

}

var StepbyChan = make(chan string)

var RW3Chan = make(chan string)

func readBlockchain(conn net.Conn) {
	for {
		//var globe Data
		stopchan <- "hello"
		fmt.Println("In read blockchain before decode:")
		var chainhead *Block
		gobEncoder := gob.NewDecoder(conn)
		//Stuck
		err1 := gobEncoder.Decode(&chainhead)
		//Stuck
		//	fmt.Println("In Admindata: ", globe)
		if err1 != nil {
			log.Println(err1, "Error in reading Blockchain ")
		}
		if Length(chainhead) > Length(globalData.ChainHead) {
			globalData.ChainHead = chainhead
		}
		fmt.Println("Blockchain received:: ", Length(globalData.ChainHead))
		//	globalData = globe
		//	<-Globechan
		//<-RW3Chan
	}
}

func broadcastBlock() {
	StepbyChan <- "Hello"
	//	RW3Chan <- "hello"
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

//For User and Miner
func StartListening(listeningAddress string, node string) {

	if node == "server" {
		ln, err := net.Listen("tcp", "localhost:"+listeningAddress)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Faital")
		}
		//	clientsSlice = make([]Verifier, 10)
		fmt.Println("Stream Starts")
		//blockchan := make(chan *Block)
		newBlock := &Block{}
		globalData.ChainHead = InsertOnlyBlock(newBlock, globalData.ChainHead)

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			conns := Connected{
				Conn: conn,
			}
			// go broadcastBlockchaintoPeer(conn)
			// go receiveBlockchainfromPeer(conn)

			go StoreClient(conn)
			//	go readBlockchain(conn)
			go readAdminData(conn)

			globalData.ClientsSlice = append(globalData.ClientsSlice, <-addchan)
			localData = append(localData, conns)

			go broadcastAdminData()
			//	go broadcastBlock()

			//	go WriteData(conn, blockchan)

			//	fmt.Println("Slice:", globalData.clientsSlice[0].ListeningAddress)
			//	<-blockchan
			//	chainHead = <-Blockchan
		}

	} else if node == "miner" {
		ln, err := net.Listen("tcp", listeningAddress)
		if err != nil {
			log.Fatal(err, ln)
		}
		clientsSlice := make([]Peer, 10)
		//	addchan := make(chan Peer)
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			newClient := Peer{
				// Conn: conn,
			}
			clientsSlice = append(clientsSlice, newClient)
			// go broadcastBlockchaintoPeer(conn)
			// go receiveBlockchainfromPeer(conn)

			go MinerverifyBlock(conn)
		}

	}
}

//Sending course to be verified
func UserSendBlock(minerAddress string, block *Block) {
	//Input from me

	//Dialing Miner
	conn, errs := net.Dial("tcp", minerAddress)
	if errs != nil {
		log.Fatal(errs)
	}
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(block)
	if err != nil {
		//	log.Println(err)
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
		InsertOnlyBlock(recvdBlock, globalData.ChainHead)
	}
}

func broadcastBlockchaintoPeer(conn net.Conn) {
	//channel
	gobEncoder := gob.NewEncoder(conn)
	err1 := gobEncoder.Encode(globalData.ChainHead)
	if err1 != nil {
		log.Println(err1)
	}

}
func receiveBlockchainfromPeer(conn net.Conn) {
	//channel
	var newChain *Block
	gobEncoder := gob.NewDecoder(conn)
	err1 := gobEncoder.Decode(newChain)
	if err1 != nil {
		//	log.Println(err)
	}

}
func WriteString(conn net.Conn, details *Peer) {
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(details)
	if err != nil {
		//	log.Println(err)
	}
}
func readAdminData(conn net.Conn) {
	//	RW3Chan <- "hello"

	StepbyChan <- "Hello"
	fmt.Println("In reading blockchain", conn.LocalAddr())
	var globe Data
	gobEncoder := gob.NewDecoder(conn)
	err1 := gobEncoder.Decode(&globe)
	//	fmt.Println("In Admindata: ", globe.ClientsSlice[0])
	if err1 != nil {
		log.Println("admin:: ", err1)
	}
	if Length(globe.ChainHead) < Length(globalData.ChainHead) {
		globe.ChainHead = globalData.ChainHead
	}
	if len(globe.ClientsSlice) < len(globalData.ClientsSlice) {
		globe.ClientsSlice = globalData.ClientsSlice
	}
	fmt.Println("Blockchain read:")
	ListBlocks(globe.ChainHead)

	globalData.ChainHead = globe.ChainHead
}

func ViewMinerData() {
	for i := 0; i < len(globalData.ClientsSlice); i++ {
		if globalData.ClientsSlice[i].Role == "miner" {
			fmt.Println("Miners connected to system:")
			fmt.Print(" Their address: ", globalData.ClientsSlice[i].ListeningAddress)
		}
	}
}

func main() {

	//secondCourse := Course{code: "CS99", name: "DIP", creditHours: 3, grade: "B-"}
	//firstProject := Project{name: "TigerKing", document: "//Hello.cpp", course: secondCourse}
	//var chainHead *Block
	//chainHead = InsertBlock(secondCourse, firstProject, chainHead)
	//chainHead = InsertCourse(firstCourse, chainHead)
	//	ListBlocks(chainHead)

	//The function below launches the server, uses different second argument
	//It then starts a routine for each connection request received
	satoshiAddress := os.Args[1]
	go StartListening(satoshiAddress, "server")
	//log.Println("Sending my course to Verifier")

	// firstCourse := Course{code: "CS50", name: "AI", creditHours: 3, grade: "A+"}
	// minerAddress := ":4502"

	//SendCourseV(minerAddress, firstCourse)

	//Satoshi is there waiting for our address, it stores it somehow

	// ln, err := net.Listen("tcp", ":6003")
	// if err != nil {
	//
	// 	log.Fatal(err)
	//
	// }
	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		log.Println(err)
	// 		continue
	// 	}
	// 	chainHead := &Block{}
	// 	go sendBlockchain(conn, chainHead)
	// }
	select {}
}
