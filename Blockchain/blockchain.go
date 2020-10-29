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

//256bit
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
func InsertBlock(course Course, project Project, chainHead *Block) *Block {
	newBlock := &Block{
		//Hash here
		course:  course,
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
		course: course,
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
		if chainHead.course == oldCourse {
			chainHead.course = newCourse
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

		fmt.Print(" Course: ", chainHead.course.name)
		fmt.Print(" Project: ", chainHead.project.name)
		fmt.Print(" -> ")
		chainHead = chainHead.PrevPointer

	}
	fmt.Println()

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

func WriteData(conn net.Conn, blockchan chan *Block) {

	firstCourse := Course{code: "CS50", name: "AI", creditHours: 3, grade: "A+"}
	block := &Block{
		//Hash here
		course: firstCourse,
	}
	blockchan <- block
	gobEncoder := gob.NewEncoder(conn)
	err1 := gobEncoder.Encode(block)
	if err1 != nil {
		//	log.Println(err)
	}

}

var addchan = make(chan Peer)

//var clientsSlice []Verifier
func broadcastAdminData() {

	<-addchan
	for i := 0; i < len(globalData.clientsSlice); i++ {
		gobEncoder := gob.NewEncoder(globalData.clientsSlice[i].Conn)
		err1 := gobEncoder.Encode(globalData)
		fmt.Println("Broadcasting:: ")
		if err1 != nil {
			//	log.Println(err)
		}
	}

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

	newClient.Conn = conn
	globalData.clientsSlice = append(globalData.clientsSlice, newClient)
	addchan <- newClient
	// dec := gob.NewDecoder(conn)
	// p := P{}
	// dec.Decode(&p)

}

//For User and Miner
func StartListening(listeningAddress string, node string) {

	if node == "server" {
		ln, err := net.Listen("tcp", ":"+listeningAddress)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Faital")
		}
		//	clientsSlice = make([]Verifier, 10)

		//blockchan := make(chan *Block)

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			// go broadcastBlockchaintoPeer(conn)
			// go receiveBlockchainfromPeer(conn)

			go StoreClient(conn)
			go broadcastAdminData()
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
				Conn: conn,
			}
			clientsSlice = append(clientsSlice, newClient)
			// go broadcastBlockchaintoPeer(conn)
			// go receiveBlockchainfromPeer(conn)

			go MinerverifyBlock(conn)
			//	go broadcastAdminData()
			//	go WriteData(conn, blockchan)

			//	fmt.Println("Slice:", globalData.clientsSlice[0].ListeningAddress)
			//	<-blockchan
			//	chainHead = <-Blockchan
		}

	}
}

//Sending course to be verified
func SendCourseV(minerAddress string, course Course) {
	//Input from me

	//Dialing Miner
	conn, errs := net.Dial("tcp", minerAddress)
	if errs != nil {
		log.Fatal(errs)
	}
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(course)
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
		InsertOnlyBlock(recvdBlock, globalData.chainHead)
	}
}

func broadcastBlockchaintoPeer(conn net.Conn) {
	//channel
	gobEncoder := gob.NewEncoder(conn)
	err1 := gobEncoder.Encode(globalData.chainHead)
	if err1 != nil {
		//	log.Println(err)
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
