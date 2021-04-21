package blockchain

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	db "github.com/HamzaPY/FYP/Database"
	model "github.com/HamzaPY/FYP/Models"
	gomail "gopkg.in/mail.v2"

	//gmail "google.golang.org/api/gmail/v1"
	//"google.golang.org/api/option"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	//"golang.org/x/oauth2"
	//"golang.org/x/oauth2/google"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//"github.com/gorilla/websocket"
)

var chainHead *Block

var unverifiedChain *Block

//Course Structure for course content
type Course struct {
	Code        string `json:"courseCode"`
	Name        string `json:"courseName"`
	CreditHours int    `json:"courseCHrs"`
	Grade       string `json:"courseGrade"`
}

// CourseWeb for detailed course content
type CourseWeb struct {
	Code        string `json:"courseCode"`
	Name        string `json:"courseName"`
	CreditHours int    `json:"courseCHrs"`
	Grade       string `json:"courseGrade"`
	VEmail      string `json:"courseEmail"`
	SEmail      string `json:"userEmail"`
	SPass       string `json:"userPass"`
}

//Project Struct for project contents
type Project struct {
	Name       string `json:"projectName"`
	Details    string `json:"projectDetails"`
	FileName   string `json:"projectFile"`
	CourseName string `json:"projectCourse"`
}

// CourseWeb for detailed project content
type ProjectWeb struct {
	Name       string `json:"projectName"`
	Details    string `json:"projectDetails"`
	FileName   string `json:"projectFile"`
	CourseName string `json:"projectCourse"`
	VEmail     string `json:"projectEmail"`
	SEmail     string `json:"userEmail"`
	SPass      string `json:"userPass"`
}

//Block stores block information which includes hash
type Block struct {
	Course      Course
	Project     Project
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
	BlockNo     int
	Status      string
	Email       string
	Username    string
}

//Skill Not used yet
type Skill struct {
	types string
	name  string
	level string
}

//BlockHandler not used
type BlockHandler struct {
	id int
}

//CV Structure defined for web pages //
type CV struct {
	Email     string
	Firstname string
	Lastname  string
	Course    []Course
	Project   []Project
	Username  string
}

//ListTheBlock for webpages listing content
type ListTheBlock struct {
	Course      []Course
	Project     []Project
	PrevPointer []*Block
	PrevHash    []string
	CurrentHash []string
	BlockNo     []int
	Status      []string
	Email       []string
	Username    []string
}

//UnverifyBlock displaying unverified blocks
type UnverifyBlock struct {
	Course      []Course  `json:"courses"`
	Project     []Project `json:"projects"`
	PrevPointer []*Block  `json:"prevPointer"`
	PrevHash    []string  `json:"prevHash"`
	CurrentHash []string  `json:"currHash"`
	BlockNo     []int     `json:"blockNo"`
	Status      []string  `json:"status"`
	Email       []string  `json:"email"`
	Username    string    `json:"username"`
	UserEmail   string    `json:"userEmail"`
}

//Client stores info of client connected to Satoshi
type Client struct {
	ListeningAddress string
	Types            bool //true for node and false for miner
	Mail             string
}

//Combo Stores blockchain and clientsconnected
type Combo struct {
	ClientsSlice       []Client
	ChainHead          *Block
	UnverifiedCourses  *Block
	UnverifiedProjects *Block
}

//Connected just for connection
type Connected struct {
	Conn net.Conn
}

var count int = 0
var stuff Combo
var localData []Connected
var mutex = &sync.Mutex{}

var tokenString = ""
var urlLogin = ""
var chainHeadArray []*Block

var currUser model.User

//ReadBlockchainFile for reading FIle
func ReadBlockchainFile() {

	file, err := os.Open("blockchainFile.json")
	if err != nil {
		log.Println("Can't read file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.Token()
	block := Block{}
	// Appends decoded object to dataArr until every object gets parsed
	for decoder.More() {
		decoder.Decode(&block)
		chainHead = InsertCourse(block)
	}
	stuff.ChainHead = chainHead
	ListBlocks(chainHead)
}

//WriteBlockchainFile Writing into file
func WriteBlockchainFile(chainHead []Block) {
	arr := removeDuplicateValues(chainHead)
	file, _ := json.MarshalIndent(arr, "", " ")
	_ = ioutil.WriteFile("blockchainFile.json", file, 0644)
	fmt.Println("file")

}
func SaveUVFile() {

	file, err := os.Open("unverified.json")
	if err != nil {
		log.Println("Can't read file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.Token()
	block := Block{}
	// Appends decoded object to dataArr until every object gets parsed
	for decoder.More() {
		decoder.Decode(&block)
		//	chainHead = InsertCourse(block)
		if block.Project.CourseName != "" {
			chainHead = InsertProjectUnverified(block)
		}
		if block.Course.Name != "" {
			chainHead = InsertCourseUnverified(block)
		}

	}
	stuff.ChainHead = chainHead
	ListBlocks(chainHead)
}

//WriteUVFile Writing into file
func WriteUVFile(chainHead []Block) {

	file, _ := json.MarshalIndent(chainHead, "", " ")
	_ = ioutil.WriteFile("unverified.json", file, 0644)
	fmt.Println("file")

}
func ReadHash(hash string) Block {

	file, err := os.Open("unverified.json")
	if err != nil {
		log.Println("Can't read file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.Token()
	block := Block{}
	// Appends decoded object to dataArr until every object gets parsed
	for decoder.More() {
		decoder.Decode(&block)
		fmt.Println("FileBlock", block)
		//	chainHead = InsertCourse(block)
		filehas := CalculateHash(&block)
		if filehas == hash {
			return block
		}
	}
	//	stuff.ChainHead = chainHead
	return block
	//	ListBlocks(chainHead)
}

//GetBlockhainArray   Getting blockchain Data in Array
func GetBlockhainArray(chainHead *Block) []Block {
	var data []Block
	i := 0
	var block Block
	for chainHead != nil {
		block.Email = chainHead.Email
		if (chainHead.Course != Course{}) {
			block.Course = chainHead.Course
		}
		block.Status = chainHead.Status
		block.Username = chainHead.Username
		if (chainHead.Project != Project{}) {
			block.Project = chainHead.Project
		}
		data = append(data, block)
		chainHead = chainHead.PrevPointer
		i++

	}
	return data

}
func removeDuplicateValues(intSlice []Block) []Block {
	keys := make(map[Block]bool)
	list := []Block{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//CalculateHash Generates 256bit hash of block using content inside
func CalculateHash(inputBlock *Block) string {

	var temp string
	if (inputBlock.Course != Course{}) {
		temp = inputBlock.Course.Code + inputBlock.Course.Name + inputBlock.Course.Grade
	}
	if (inputBlock.Project != Project{}) {
		temp = inputBlock.Project.CourseName + inputBlock.Project.Name + inputBlock.Project.Details + inputBlock.Project.FileName
	}
	h := sha256.New()
	h.Write([]byte(temp))
	sum := hex.EncodeToString(h.Sum(nil))

	// sum := sha256.Sum256([]byte(temp))

	return sum
}

//InsertCourse Insert Verified Course in chain //
func InsertCourse(myBlock Block) *Block {

	myBlock.CurrentHash = CalculateHash(&myBlock)
	fmt.Println("Course Hash, ", CalculateHash(&myBlock))
	if chainHead == nil {
		myBlock.Status = "Verified"
		myBlock.BlockNo = count
		myBlock.PrevHash = "null"
		chainHead = &myBlock
		fmt.Println("Genesis Course Block Inserted!")
		return chainHead
	}
	count = count + 1
	myBlock.Status = "Verified"
	myBlock.PrevPointer = chainHead
	myBlock.PrevHash = chainHead.CurrentHash
	myBlock.BlockNo = count

	fmt.Println("Course Block Inserted!")
	return &myBlock

}

//InsertProject Insert Verified Project //
func InsertProject(myBlock Block) *Block {

	myBlock.CurrentHash = CalculateHash(&myBlock)
	fmt.Println("Course Hash, ", CalculateHash(&myBlock))
	if chainHead == nil {
		myBlock.BlockNo = count
		myBlock.PrevHash = "null"
		chainHead = &myBlock
		fmt.Println("Genesis Project Block Inserted!")
		return chainHead
	}
	count = count + 1
	myBlock.PrevPointer = chainHead
	myBlock.PrevHash = chainHead.CurrentHash
	myBlock.BlockNo = count

	fmt.Println("Project Block Inserted!")
	return &myBlock

}

//InsertCourseUnverified Unverified Course Chain //
func InsertCourseUnverified(myBlock Block) *Block {

	myBlock.CurrentHash = CalculateHash(&myBlock)
	fmt.Println("Course Hash, ", CalculateHash(&myBlock))
	if unverifiedChain == nil {
		myBlock.BlockNo = count
		myBlock.PrevHash = "null"
		unverifiedChain = &myBlock
		fmt.Println("Genesis Course Unverified Block Inserted!")

		return unverifiedChain
	}
	count = count + 1
	myBlock.PrevPointer = unverifiedChain
	myBlock.PrevHash = unverifiedChain.CurrentHash
	myBlock.BlockNo = count

	fmt.Println("Course Unverified Block Inserted!")
	return &myBlock

}

//InsertProjectUnverified  Unverified Project Chain //
func InsertProjectUnverified(myBlock Block) *Block {

	myBlock.CurrentHash = CalculateHash(&myBlock)
	fmt.Println("Course Hash, ", CalculateHash(&myBlock))
	if unverifiedChain == nil {
		myBlock.BlockNo = count
		myBlock.PrevHash = "null"
		unverifiedChain = &myBlock
		fmt.Println("Genesis Project Unverified Block Inserted!")
		return unverifiedChain
	}
	count = count + 1
	myBlock.PrevPointer = unverifiedChain
	myBlock.PrevHash = unverifiedChain.CurrentHash
	myBlock.BlockNo = count

	fmt.Println("Project Unverified Block Inserted!")
	return &myBlock

}

//ChangeCourse for changing details not used yet
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

//ChangeProject for changing details not used yet
func ChangeProject(oldProject Project, newProject Project, chainHead *Block) {
	present := false
	for chainHead != nil {
		if chainHead.Project == oldProject {
			chainHead.Project = newProject
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

//ListBlocks Prints Blocks in Commandline
func ListBlocks(chainHead *Block) {

	for chainHead != nil {
		fmt.Print("Block NO: ", chainHead.BlockNo)
		fmt.Print(" Current Hash: ", chainHead.CurrentHash)
		if chainHead.PrevHash == "" {
			fmt.Print(" Previous Hash: ", "Null")
		} else {
			fmt.Print(" Previous Hash: ", chainHead.PrevHash)
		}

		fmt.Print(" Course: ", chainHead.Course.Name)
		fmt.Print(" Project: ", chainHead.Project.Name)
		fmt.Print(" -> ")
		chainHead = chainHead.PrevPointer

	}
	fmt.Println()

}

//VerifyChain Checks if hash of chain is changed or not
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

//ReceiveChain Recieves blockchain the firsttime
func ReceiveChain(conn net.Conn) *Block {
	fmt.Println("In func")
	var block *Block
	gobEncoder := gob.NewDecoder(conn)
	err := gobEncoder.Decode(&block)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Received chain")
	chainHead = block
	stuff.ChainHead = chainHead
	ListBlocks(chainHead)

	//chainHead = InsertCourse(block)
	return block
}

//Length length of blockchain
func Length(chainHead *Block) int {
	sum := 0
	for chainHead != nil {

		chainHead = chainHead.PrevPointer
		sum++
	}
	return sum

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

//InsertCourse1 not used
func InsertCourse1(course Course, chainHead *Block) *Block {
	newBlock := &Block{
		//Hash here
		Course: course,
	}
	newBlock.CurrentHash = CalculateHash(newBlock)

	if chainHead == nil {
		chainHead = newBlock
		chainHead.BlockNo = count
		fmt.Println("Block Inserted")
		return chainHead
	}
	count = count + 1
	newBlock.PrevPointer = chainHead
	newBlock.PrevHash = chainHead.CurrentHash
	newBlock.BlockNo = count

	fmt.Println("Course Block Inserted")
	return newBlock

}
func getCourse(ChainHead *Block) []Block {
	var courses []Block
	for chainHead != nil {
		courses = append(courses, *chainHead)
		chainHead = chainHead.PrevPointer
	}
	//	fmt.Println("Yo")
	return courses
}

//WriteString for writing info to satoshi
func WriteString(conn net.Conn, myListeningAddress Client) {
	Satoshiconn = conn //Saves Satoshi
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(myListeningAddress)
	if err != nil {
		log.Println("In Write String: ", err)
	}
}

//SendChain not used yet
func SendChain(conn net.Conn) {
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(chainHead)
	if err != nil {
		log.Println("In Write Chain: ", err)
	}
}

//Satoshiconn to store connection for Satoshi
var Satoshiconn net.Conn
var clientsSlice []Client
var rwchan = make(chan string)

// Satoshi handles client using this
func handleConnection(conn net.Conn, addchan chan Client) {
	// newClient := Connected{
	// 	Conn: conn,
	// }
	Clientz := Client{}
	//var ListeningAddress string
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&Clientz)
	if err != nil {
		//handle error
	}

	// newClient.ListeningAddress = ListeningAddress
	fmt.Println("inHandle: ", Clientz.ListeningAddress)
	addchan <- Clientz //Blocks until client is added to its list
	//WaitForQuorum()

}

var nodesSlice []Client
var minechan = make(chan Client)

var blockchan = make(chan Block)

//Minedblock Mined Block in this var
var Minedblock Block

var newchan = make(chan *Block)

//Handles connection of node and miner
func handlePeer(conn net.Conn) {

	//	Clientz := Client{}
	block := Block{}
	//var ListeningAddress string
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&block)
	if err != nil {
		//handle error
		log.Print("Eror in receiveing block", block)
	}

	// newClient.ListeningAddress = ListeningAddress
	fmt.Println("inHandlePeer: ", block)
	blockchan <- block

}

//ReceiveMinerChain    Not Usedd
func ReceiveMinerChain(conn net.Conn) *Block {
	fmt.Println("In func")
	var block *Block
	gobEncoder := gob.NewDecoder(conn)
	err := gobEncoder.Decode(&block)
	if err != nil {
		log.Println(err)
	}
	if Length(chainHead) <= Length(block) {
		fmt.Println("Received new chain")
		chainHead = block
	} else {
		fmt.Println("Received old chain")

	}
	ListBlocks(chainHead)
	gobEncoder2 := gob.NewEncoder(conn)
	err2 := gobEncoder2.Encode(&chainHead)
	if err2 != nil {
		log.Println(err2)
	}

	return block
}

//ReceiveEverything Satoshi receives everything
func ReceiveEverything(conn net.Conn) { //Admin
	for {
		fmt.Println("In Recieved  func Doit", Doit)
		var stuu Combo
		gobEncoder := gob.NewDecoder(conn)
		err := gobEncoder.Decode(&stuu)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Received Stuff chain")

		//ListBlocks(stuu.ChainHead)
		fmt.Println("Received head chain")
		ListBlocks(chainHead)
		if stuu.ChainHead != nil {
			if Length(chainHead) <= Length(stuu.ChainHead) { //Checks and retains the longer blockchain
				fmt.Println("Received new chain")
				chainHead = stuu.ChainHead
				stuff.ChainHead = chainHead
				data := GetBlockhainArray(chainHead)
				data = removeDuplicateValues(data)
				WriteBlockchainFile(data)
			} else {
				fmt.Println("Received old chain")
			}
		}
		ListBlocks(chainHead)

	}
	// if Doit == false {
	// 	log.Println("First Time")
	// 	gobEncoder2 := gob.NewEncoder(conn)
	// 	err2 := gobEncoder2.Encode(&stuff)
	// 	if err2 != nil {
	// 		log.Println(err2)
	// 	}
	// }

}

//ReceiveChain1 not used
func ReceiveChain1(conn net.Conn) *Block {
	//<-check
	for {
		rwchan <- "sss"
		var block *Block
		gobEncoder := gob.NewDecoder(conn)
		err := gobEncoder.Decode(&block)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Received chain")
		chainHead = block
		ListBlocks(chainHead)

		//chainHead = InsertCourse(block)
	}
	//	return block
}

var j int

//not usedddd
func broadcastPeerData() {

	for i := 0; i < len(localData); i++ {
		gobEncoder := gob.NewEncoder(localData[i].Conn)
		//fmt.Println("BroadCheck: ", localData[i])
		err1 := gobEncoder.Encode(clientsSlice)
		fmt.Println("Broadcasting PeerData:: ")
		if err1 != nil {
			log.Println("Errpr in broadcasting", err1)
		}

	}

	//	<-StepbyChan

}

//not used
func broadcastChain() {

	for i := 0; i < len(localData); i++ {
		//		fmt.Println("ss", nodesSlice[i].Types)
		gobEncoder := gob.NewEncoder(localData[i].Conn)
		//fmt.Println("BroadCheck: ", localData[i])
		err1 := gobEncoder.Encode(chainHead)
		fmt.Println("Broadcasting Chain to:: ", localData[i].Conn)
		if err1 != nil {
			log.Println("Errpr in broadcasting Chain", err1)
		}

	}

	//	<-StepbyChan

}

//not used
func broadcastEverything() {
	// stuff.ChainHead = chainHead
	// stuff.ClientsSlice = nodesSlice
	for i := 0; i < len(localData); i++ {
		//		fmt.Println("ss", nodesSlice[i].Types)
		gobEncoder := gob.NewEncoder(localData[i].Conn)
		//fmt.Println("BroadCheck: ", localData[i])
		err1 := gobEncoder.Encode(stuff)
		fmt.Println("Broadcasting Chain to:: ", localData[i].Conn)
		if err1 != nil {
			log.Println("Errpr in broadcasting Chain", err1)
		}

	}

	//	<-StepbyChan

}

//ReadPeers Reading peers    (Not Usedd yett)
func ReadPeers(conn net.Conn) []Client {
	//	for {
	//	mutex.Lock()
	var slice []Client
	gobEncoder := gob.NewDecoder(conn)
	err := gobEncoder.Decode(&slice)
	if err != nil {
		log.Println(err)
	}
	nodesSlice = slice
	fmt.Println("Read Peers: ", nodesSlice, len(nodesSlice))
	//	mutex.Unlock()
	//		check <- "d"
	//	}
	return nodesSlice
}

//ReadPeers1    not useddd
func ReadPeers1(conn net.Conn) []Client {
	for {
		//	mutex.Lock()

		var slice []Client
		gobEncoder := gob.NewDecoder(conn)
		err := gobEncoder.Decode(&slice)
		if err != nil {
			log.Println(err)
		}
		nodesSlice = slice
		fmt.Println("Read Peers: ", nodesSlice)

		//		<-rwchan
		//	mutex.Unlock()
		//		check <- "d"
	}
	//	return nodesSlice
}

//ReadPeersMinerChain       not useddddd
func ReadPeersMinerChain(conn net.Conn) []Client {
	for {
		//	mutex.Lock()
		if Doit != false {
			var slice []Client
			fmt.Println("In Read Peers ggg")
			gobEncoder := gob.NewDecoder(conn)
			err := gobEncoder.Decode(&slice)
			if err != nil {
				log.Println(err, "FFF")
			}
			nodesSlice = slice
			fmt.Println("Read Peers: ", nodesSlice)
		}
		//	ReceiveChain(conn)

		//		<-rwchan
		//	mutex.Unlock()
		//		check <- "d"
	}
	//	return nodesSlice
}

//ReadPeersMinerChainEverything miner recieves chain and everything
func ReadPeersMinerChainEverything(conn net.Conn) { //Miner
	for {
		//	mutex.Lock()
		var stuu Combo
		fmt.Println("In Read Peers ggg")
		gobEncoder := gob.NewDecoder(conn)
		err := gobEncoder.Decode(&stuu)
		if err != nil {
			log.Println(err, "FFF")
		}
		fmt.Println("Read StuuPeers: ", stuu.ClientsSlice)
		if len(stuu.ClientsSlice) >= len(nodesSlice) {
			nodesSlice = stuu.ClientsSlice
			stuff.ClientsSlice = nodesSlice
			fmt.Println("Read Peers: ", nodesSlice)

		}
		if Length(stuu.ChainHead) >= Length(chainHead) {
			chainHead = stuu.ChainHead
			stuff.ChainHead = chainHead
			fmt.Println("Read Chain: ")
			ListBlocks(chainHead)
		}

		//	ReceiveChain(conn)

		//		<-rwchan
		//	mutex.Unlock()
		//		check <- "d"
	}
	//	return nodesSlice
}

//ReadBlockPeers not used
func ReadBlockPeers(conn net.Conn) Block {
	var block Block
	gobEncoder := gob.NewDecoder(conn)
	err := gobEncoder.Decode(&block)
	if err != nil {
		log.Println(err)
	}
	return block
}

//StartListening Server for Satoshi node and miner
func StartListening(ListeningAddress string, node string) {
	if node == "satoshi" {
		ln, err := net.Listen("tcp", ":"+ListeningAddress)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Faital")
		}
		j = 0
		addchan := make(chan Client)
		block := Block{}
		chainHead = InsertCourse(block) //Genesis Block
		ReadBlockchainFile()
		stuff.ChainHead = chainHead
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err, "Yooooo")
				continue
			}
			sendBlockchain(conn, chainHead) //sends chain the first time to each user that connects
			conns := Connected{
				Conn: conn,
			}

			go handleConnection(conn, addchan)             //Reads Clients Satoshi
			clientsSlice = append(clientsSlice, <-addchan) //adds client to list
			stuff.ClientsSlice = clientsSlice
			fmt.Println("stuffCl: ", stuff.ClientsSlice)
			fmt.Println("clS: ", clientsSlice)
			localData = append(localData, conns) //adds connection
			//fmt.Println("BroadCheck: ", localData[i])
			//		broadcastPeerData()
			//		broadcastEverything()

			go func() { //Broadcasting data after every x seconds to each node/miner connected
				for {
					time.Sleep(8 * time.Second)
					mutex.Lock()
					for i := 0; i < len(localData); i++ {
						//		fmt.Println("ss", nodesSlice[i].Types)
						gobEncoder := gob.NewEncoder(localData[i].Conn)
						//fmt.Println("BroadCheck: ", localData[i])
						err1 := gobEncoder.Encode(stuff)
						fmt.Println("Broadcasting Chain to:: ", localData[i].Conn)
						if err1 != nil {
							log.Println("Errpr in broadcasting Chain", err1)
						}

					}
					mutex.Unlock()

				}
			}()

			go ReceiveEverything(conn) //Receives chaina nd anything fromn nodes/miner connected

			ListBlocks(chainHead)

		}
	} else if node == "others" { //For nodes own server
		ln, err := net.Listen("tcp", ":"+ListeningAddress)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Faital")
		}

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err, "Yooooo")
				continue
			}
			go handlePeer(conn)
			nodesSlice = append(nodesSlice, <-minechan)

		}

	} else { //miner's own server
		ln, err := net.Listen("tcp", "localhost:"+ListeningAddress)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Faital")
		}

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err, "Yooooo")
				continue
			}
			fmt.Println("COnnedted")
			testConn = conn
			conns := Connected{
				Conn: conn,
			}
			localData = append(localData, conns)

			go handlePeer(conn) //Handles nodes connected

			Minedblock = <-blockchan

		}
	}

}

var testConn net.Conn //Users connection stored

/// Mux Router HTTP Services ///

//MinerConn connection stored for user
var MinerConn net.Conn

//Mined ; for checking if block is mined or not
var Mined bool

/// Web Handler to show all blocks of blockchain in satoshi server ///
func showBlocksHandler(w http.ResponseWriter, r *http.Request) {
	tempHead := chainHead
	viewTheBlock := new(ListTheBlock)
	tempCourse := []Course{}
	tempBlockNo := []int{}
	tempCurrHash := []string{}
	tempPrevHash := []string{}
	tempStatus := []string{}
	for tempHead != nil {
		tempCourse = append(tempCourse, tempHead.Course)
		tempBlockNo = append(tempBlockNo, tempHead.BlockNo)
		tempCurrHash = append(tempCurrHash, tempHead.CurrentHash)
		tempPrevHash = append(tempPrevHash, tempHead.PrevHash)
		tempStatus = append(tempStatus, tempHead.Status)

		viewTheBlock = &ListTheBlock{
			Course:      tempCourse,
			BlockNo:     tempBlockNo,
			CurrentHash: tempCurrHash,
			PrevHash:    tempPrevHash,
			Status:      tempStatus,
		}
		tempHead = tempHead.PrevPointer
		fmt.Println(viewTheBlock.Course)
		fmt.Println(viewTheBlock.BlockNo)
		fmt.Println(viewTheBlock.CurrentHash)
		fmt.Println(viewTheBlock.PrevHash)
		fmt.Println(viewTheBlock.Status)
	}
	// generate page by passing page variables into template
	t, err := template.ParseFiles("../Website/viewBlocks.html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, viewTheBlock) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

var check = make(chan string)

//Doit check if miner has clicked the link
var Doit bool

//Mineblock  Web Handler to verify and mine the block in miner server ///
func Mineblock(w http.ResponseWriter, r *http.Request) {
	Satoshiconn, err := net.Dial("tcp", ":"+"2500")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("In Mine Block")
	Doit = true
	params := mux.Vars(r)
	mineHash := params["hash"]
	fmt.Println(mineHash)

	ReceiveChain(Satoshiconn) //Receives chain the first time from Satoshi

	Peers := Client{
		Types: false,
	}
	WriteString(Satoshiconn, Peers) //Sends its info including his mail to Satoshi

	// go ReadPeersMinerChainEverything(Satoshiconn) //Reads info from Satoshi
	//	blockHash := CalculateHash(&Minedblock)
	//	ListBlocks(unverifiedChain)
	//	blockHash := CalculateHash(unverifiedChain)
	block := ReadHash(mineHash)
	blockHash := CalculateHash(&block)

	fmt.Println("Yoo", blockHash)
	if blockHash == mineHash { //Checks if hash is same as the block

		Minedblock.Status = "Verified"
		chainHead = InsertCourse(block)
		stuff.ChainHead = chainHead
		fmt.Println("In Mining")
		fmt.Println(chainHead)
		ListBlocks(chainHead)
		fmt.Println("Trrr")

		gobEncoder := gob.NewEncoder(Satoshiconn)
		err2 := gobEncoder.Encode(stuff)
		if err2 != nil {
			log.Println("In Error Write Chain: ", err2)
		}
		log.Println("Sent to Satoshi: ")

		// gobEncoder2 := gob.NewEncoder(testConn)
		// //fmt.Println("BroadCheck: ", localData[i])
		// err1 := gobEncoder2.Encode(stuff)
		// fmt.Println("Bro Chain sent to peer:: ", testConn)
		// if err1 != nil {
		// 	log.Println("Errpr in brosti Chain", err1)
		// }
	} else {
		fmt.Println("Wrong Hash")
	}
	//	broadcastChain()

}

// Clients Web Server //

//RegisterHandler   Web Handler to register user into DB ///
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var newUser model.User
		_ = json.NewDecoder(r.Body).Decode(&newUser)
		fmt.Println(newUser)

		collection, err := db.GetDBCollection()

		var result model.User
		err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: newUser.Username}}).Decode(&result)

		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 5)
				if err != nil { // if there is an error
					log.Print("Error ", err) //log it
				}
				newUser.Password = string(hash)

				_, err = collection.InsertOne(context.TODO(), newUser)
				if err != nil { // if there is an error
					log.Print("Error ", err) //log it
				}
			}
		}
	}

}

//LoginHandler   Web Handler to login user using JWT Authentication ///
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var newUser model.User
		_ = json.NewDecoder(r.Body).Decode(&newUser)
		fmt.Println(newUser)

		collection, err := db.GetDBCollection()

		if err != nil {
			log.Fatal(err)
		}
		var result model.User
		var res model.ResponseResult

		err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: newUser.Username}}).Decode(&result)

		if err != nil {
			res.Error = "Invalid username"
			fmt.Println("Invalid username")
			json.NewEncoder(w).Encode(res)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(newUser.Password))

		if err != nil {
			res.Error = "Invalid password"
			fmt.Println("Invalid password")
			json.NewEncoder(w).Encode(res)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username":  result.Username,
			"firstname": result.FirstName,
			"lastname":  result.LastName,
			"email":     result.Email,
		})

		tokenString, err = token.SignedString([]byte("secret"))

		if err != nil {
			fmt.Println(err)
		}

		json, err := json.Marshal(struct {
			Token string `json:"token"`
		}{
			tokenString,
		})

		if err != nil {
			fmt.Println(err)
		}

		//http.Redirect(w, r, urlLogin+"/login", http.StatusSeeOther)
		satoshiAddress := "2500"
		myListeningAddress := "6002"

		conn, err := net.Dial("tcp", ":"+satoshiAddress)
		if err != nil {
			log.Fatal(err)
		}

		go StartListening(myListeningAddress, "others") //Starts own server

		log.Println("Sending my listening address to Satoshis")
		chainHead := ReceiveChain(conn)
		ListBlocks(chainHead)

		Peers := Client{
			ListeningAddress: myListeningAddress,
			Types:            true,
		}
		WriteString(conn, Peers) //Writes its address

		//go b.ReceiveChain(conn)

		go ReadPeersMinerChainEverything(conn) // Reads information from Satoshi every second

		// go func() { //Go routine for reading the chain that miner sends
		// 	for {
		// 		if Mined == true { // checks if the block sent is mined or not
		// 			fmt.Println("trueue")
		// 			var stuu Combo
		// 			fmt.Println("In Read Peers fffwd")
		// 			gobEncoder := gob.NewDecoder(MinerConn)
		// 			err := gobEncoder.Decode(&stuu)
		// 			if err != nil {
		// 				log.Println(err, "FFF")
		// 			}
		// 			fmt.Println("Read StuuPeers: ", stuu.ClientsSlice)
		// 			ListBlocks(stuu.ChainHead)
		// 			// if Length(stuu.ChainHead) >= Length(chainHead) {
		// 			// 	chainHead = stuu.ChainHead
		// 			// 	stuff.ChainHead = chainHead
		// 			// 	fmt.Println("Read Chain: ")
		// 			// 	ListBlocks(chainHead)
		// 			// }
		// 			Mined = false
		// 		}
		//
		// 	}
		//
		// }()

		w.Write(json)
	}
}

//ProfileHandler   Web Handler to show dashboard to the user ///
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println(" ---- Access Denied ----")
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if token == nil {
		fmt.Println(" ---- Access Denied ----")
		//http.Redirect(w, r, urlLogin+"/login", http.StatusSeeOther)
		//return
	}
	if err != nil {
		fmt.Println(err)
	}
	var result model.User
	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)
		result.FirstName = claims["firstname"].(string)
		result.LastName = claims["lastname"].(string)
		result.Email = claims["email"].(string)
		currUser = model.User{
			Username:  result.Username,
			FirstName: result.FirstName,
			LastName:  result.LastName,
			Email:     result.Email,
			Password:  "",
		}

		json, err := json.Marshal(struct {
			Result model.User `json:"result"`
		}{
			result,
		})

		if err != nil {
			fmt.Println(err)
		}

		w.Write(json)
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}

//AllUsers  Web Handler to get all the users registered in the system ///
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var newUser model.User
	_ = json.NewDecoder(r.Body).Decode(&newUser)
	fmt.Println(newUser)

	collection, err := db.GetDBCollection()

	if err != nil {
		log.Fatal(err)
	}

	cursor, err := collection.Find(context.TODO(), bson.M{})

	var episodes []bson.M

	if err = cursor.All(context.TODO(), &episodes); err != nil {
		log.Fatal(err)
	}

	json, err := json.Marshal(episodes)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(json)
}

//UnverifiedBlocksHandler Web Handler to show UnVerified Blocks in blockchain to the user ///
func UnverifiedBlocksHandler(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println(" ---- Access Denied ----")
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if token == nil {
		fmt.Println(" ---- Access Denied ----")
		http.Redirect(w, r, urlLogin+"/login", http.StatusSeeOther)
		return
	}
	var result model.User
	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)
		result.FirstName = claims["firstname"].(string)
		result.LastName = claims["lastname"].(string)
		result.Email = claims["email"].(string)
		currUser = model.User{
			Username:  result.Username,
			FirstName: result.FirstName,
			LastName:  result.LastName,
			Email:     result.Email,
			Password:  "",
		}

		fmt.Println(unverifiedChain)
		fmt.Println(chainHead)
		tempHead2 := chainHead
		tempCurrHash2 := []string{}
		for tempHead2 != nil {
			if tempHead2.Username == result.Username {
				tempCurrHash2 = append(tempCurrHash2, tempHead2.CurrentHash)
			}
			tempHead2 = tempHead2.PrevPointer
		}
		fmt.Println("In Unverif")
		ListBlocks(unverifiedChain)
		tempHead3 := unverifiedChain
		for tempHead3 != nil {
			if tempHead3.Username == result.Username {
				for i := 0; i < len(tempCurrHash2); i++ {
					fmt.Println(tempHead3.CurrentHash, tempCurrHash2[i])
					if tempHead3.Status == "Pending" && tempHead3.CurrentHash == tempCurrHash2[i] {
						tempHead3.Status = "Verified"
					}
				}
			}
			tempHead3 = tempHead3.PrevPointer
		}

		tempHead := unverifiedChain
		viewTheBlock := new(UnverifyBlock)
		tempProject := []Project{}
		tempCourse := []Course{}
		tempBlockNo := []int{}
		tempCurrHash := []string{}
		tempPrevHash := []string{}
		tempEmail := []string{}
		tempStatus := []string{}

		extCourse := new(Course)
		extProject := new(Project)

		ListBlocks(tempHead)

		for tempHead != nil {
			fmt.Println("1st If")
			if tempHead.Username == result.Username {
				fmt.Println("2nd If")

				if tempHead.Status == "Pending" {
					fmt.Println("3rd If")

					if tempHead.Course.Name == "" {
						fmt.Println("4th If Course")

						tempStatus = append(tempStatus, tempHead.Status)
						tempProject = append(tempProject, tempHead.Project)
						tempCourse = append(tempCourse, *extCourse)
						tempBlockNo = append(tempBlockNo, tempHead.BlockNo)
						tempCurrHash = append(tempCurrHash, tempHead.CurrentHash)
						tempPrevHash = append(tempPrevHash, tempHead.PrevHash)
						tempEmail = append(tempEmail, tempHead.Email)
						viewTheBlock = &UnverifyBlock{
							Course:      tempCourse,
							Project:     tempProject,
							BlockNo:     tempBlockNo,
							CurrentHash: tempCurrHash,
							PrevHash:    tempPrevHash,
							Email:       tempEmail,
							Status:      tempStatus,
							Username:    result.Username,
							UserEmail:   result.Email,
						}
						fmt.Println(viewTheBlock.Project)
						fmt.Println(viewTheBlock.BlockNo)
						fmt.Println(viewTheBlock.CurrentHash)
						fmt.Println(viewTheBlock.PrevHash)
						fmt.Println(viewTheBlock.Email)
						fmt.Println(viewTheBlock.Status)
					}
					if tempHead.Project.Name == "" {
						fmt.Println("4th If Project")
						tempStatus = append(tempStatus, tempHead.Status)
						tempProject = append(tempProject, *extProject)
						tempCourse = append(tempCourse, tempHead.Course)
						tempBlockNo = append(tempBlockNo, tempHead.BlockNo)
						tempCurrHash = append(tempCurrHash, tempHead.CurrentHash)
						tempPrevHash = append(tempPrevHash, tempHead.PrevHash)
						tempEmail = append(tempEmail, tempHead.Email)
						viewTheBlock = &UnverifyBlock{
							Course:      tempCourse,
							Project:     tempProject,
							BlockNo:     tempBlockNo,
							CurrentHash: tempCurrHash,
							PrevHash:    tempPrevHash,
							Email:       tempEmail,
							Status:      tempStatus,
							Username:    result.Username,
							UserEmail:   result.Email,
						}
						fmt.Println(viewTheBlock.Course)
						fmt.Println(viewTheBlock.BlockNo)
						fmt.Println(viewTheBlock.CurrentHash)
						fmt.Println(viewTheBlock.PrevHash)
						fmt.Println(viewTheBlock.Email)
						fmt.Println(viewTheBlock.Status)
					}
					fmt.Println(viewTheBlock)
				}
			}
			tempHead = tempHead.PrevPointer
		}

		json, err := json.Marshal(struct {
			Result UnverifyBlock `json:"unVerifyBlock"`
		}{
			*viewTheBlock,
		})

		if err != nil {
			fmt.Println(err)
		}

		w.Write(json)
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

//GenerateCVHandler  Web Handler to generate CV for the user ///
func GenerateCVHandler(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println(" ---- Access Denied ----")
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if token == nil {
		fmt.Println(" ---- Access Denied ----")
		http.Redirect(w, r, urlLogin+"/login", http.StatusSeeOther)
		return
	}
	var result model.User
	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)
		result.FirstName = claims["firstname"].(string)
		result.LastName = claims["lastname"].(string)
		result.Email = claims["email"].(string)
		currUser = model.User{
			Username:  result.Username,
			FirstName: result.FirstName,
			LastName:  result.LastName,
			Email:     result.Email,
			Password:  "",
		}
		tempHead := chainHead
		tempCourse := []Course{}
		tempProject := []Project{}

		fmt.Println("Blockchain", tempHead)
		for tempHead != nil {
			if tempHead.Username == result.Username {
				if tempHead.Course.Code != "" {
					tempCourse = append(tempCourse, tempHead.Course)
				}
				if tempHead.Project.Name != "" {
					tempProject = append(tempProject, tempHead.Project)
				}
			}
			tempHead = tempHead.PrevPointer
		}

		cv := CV{
			Email:     result.Email,
			Firstname: result.FirstName,
			Lastname:  result.LastName,
			Course:    tempCourse,
			Project:   tempProject,
			Username:  result.Username,
		}

		json, err := json.Marshal(struct {
			Result CV `json:"cv"`
		}{
			cv,
		})

		if err != nil {
			fmt.Println(err)
		}

		w.Write(json)
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

/*var GmailService *gmail.Service

func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     "275469437806-nqp90b6739i86oiupk45236jinc1h2eh.apps.googleusercontent.com",
		ClientSecret: "GEULiFVDqfZkdQswdLlERnR3",
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  "your_access_token",
		RefreshToken: "your_refresh_token",
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		fmt.Println("Email service is initialized ")
	}
}

func SendEmailOAUTH2(to string, data Block) (bool, error) {

	// emailBody, err := parseTemplate(template, data)
	// if err != nil {
	// 	return false, errors.New("unable to parse email template")
	// }

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"

	subject := "Subject: " + "Test Email form Gmail API using OAuth" + "\n"

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"

	msg := []byte(emailTo + subject + mime + "\n" + data.Project.CourseName)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	fmt.Println("In handler7")

	return true, nil
}*/

//AddProjectHandler Web Handler to add projects into the blockchain ///
func AddProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tempProject ProjectWeb
		_ = json.NewDecoder(r.Body).Decode(&tempProject)
		fmt.Println(tempProject)
		currUsername := currUser.Username

		var newProject Project
		newProject.Name = tempProject.Name
		newProject.Details = tempProject.Details
		newProject.FileName = tempProject.FileName
		newProject.CourseName = tempProject.CourseName

		// Use pFile for sending files to mailer //

		//fmt.Println(pFile, pErr)

		/////////////////////////////////

		MyBlock := Block{
			Project:  newProject,
			Email:    tempProject.VEmail,
			Username: currUsername,
			Status:   "Pending",
		}

		//chainHead = InsertCourse(MyBlock)

		// gobEncoder := gob.NewEncoder(Satoshiconn)
		// err2 := gobEncoder.Encode(MyBlock)
		// if err2 != nil {
		// 	log.Println("In Write Chain: ", err2)
		// }

		fmt.Println(MyBlock)
		unverifiedChain = InsertProjectUnverified(MyBlock)
		ListBlocks(unverifiedChain)

		uv := GetBlockhainArray(unverifiedChain)
		WriteUVFile(uv)
		//	fmt.Println("FFFFFFFFFF", len(nodesSlice))
		//	for i := 0; i < len(nodesSlice); i++ {
		//	fmt.Println("dddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
		//		if nodesSlice[i].Mail == MyBlock.Email {
		// conn, err := net.Dial("tcp", "localhost:"+nodesSlice[i].ListeningAddress)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// MinerConn = conn
		// gobEncoder := gob.NewEncoder(conn)
		// fmt.Println("blok:ahsh: ", CalculateHash(&MyBlock))
		// err2 := gobEncoder.Encode(MyBlock)
		// if err2 != nil {
		// 	log.Println("In Write Chain: ", err2)
		// }
		//	SendEmailOAUTH2(MyBlock.Email, MyBlock)

		//SMTPPPPPPPPPPPPPPPPPPPPPPPP
		// from := tempProject.SEmail
		// password := tempProject.SPass
		//
		// // Receiver email address.
		// to := []string{
		// 	MyBlock.Email,
		// }
		//
		// // smtp server configuration.
		// smtpHost := "smtp.gmail.com"
		// smtpPort := "587"
		//
		// // Message.
		// message := []byte("Project Name: " + MyBlock.Project.Name + "  Project Details: " + MyBlock.Project.Details + "  Course Grade: " + MyBlock.Course.Grade + "\n" + "Click here to verify this content: " + "localhost:" + "4000" + "/mineBlock/" + CalculateHash(&MyBlock))
		//
		// // Authentication.
		// auth := smtp.PlainAuth("", from, password, smtpHost)
		//
		// // Sending email.
		// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
		//
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// fmt.Println("Email Sent Successfully!")
		//SMTPPPPPPPPPPPPPPPPPPPPPPPP

		m := gomail.NewMessage()

		// Set E-Mail sender
		m.SetHeader("From", tempProject.SEmail)

		// Set E-Mail receivers
		m.SetHeader("To", MyBlock.Email)

		// Set E-Mail subject
		m.SetHeader("Subject", "Verification Content")

		// Set E-Mail body. You can set plain text or html with text/html

		///////////// Add files to send to the mailer /////////////////
		m.SetBody("text/plain", "Project Name: "+MyBlock.Project.Name+"  Project Details: "+MyBlock.Project.Details+"  Course Grade: "+MyBlock.Course.Grade+"\n"+"Click here to verify this content: "+"localhost:"+"4000"+"/mineBlock/"+CalculateHash(&MyBlock))

		// Settings for SMTP server
		d := gomail.NewDialer("smtp.gmail.com", 587, tempProject.SEmail, tempProject.SPass)

		// This is only needed when SSL/TLS certificate is not valid on server.
		// In production this should be set to false.
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// Now send E-Mail
		if err := d.DialAndSend(m); err != nil {
			fmt.Println(err, "mailerr")
			panic(err)
		}
		Mined = true
		//	fmt.Println("Email Sent", Mined, nodesSlice[i].ListeningAddress)

		//break
		//		}
		//		}
		http.Redirect(w, r, urlLogin+"/dashboard", http.StatusSeeOther)
	}

}

//Search dede
type Search struct {
	CourseCode        string `json:"CourseCode"`
	CourseGrade       string `json:"CourseGrade"`
	ProjectCourseName string `json:"ProjectCourseName"`
	CourseName        string `json:"courseName"`
}

//Result dede
type Result struct {
	Username string `json:"Username"`
	Course   Course `json:"Course"`
}

func SearchVerifyContent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// currUsername := currUser.Username
		file, err := os.Open("blockchainFile.json")
		if err != nil {
			log.Println("Can't read file")
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		decoder.Token()
		block := Block{}
		// Appends decoded object to dataArr until every object gets parsed
		var allCourses []Course
		for decoder.More() {
			decoder.Decode(&block)
			if block.Username != "" {
				allCourses = append(allCourses, block.Course)
			}
		}
		// ReadBlockchainFile()
		// tempHead:=chainHead
		// for tempHead != nil {
		// 	if tempHead.Username == result.Username {
		// 		if tempHead.Course.Code != "" {
		// 			tempCourse = append(tempCourse, tempHead.Course)
		// 		}
		// 		if tempHead.Project.Name != "" {
		// 			tempProject = append(tempProject, tempHead.Project)
		// 		}
		// 	}
		// 	tempHead = tempHead.PrevPointer
		// }
		// cv := CV{
		// 	Email:     result.Email,
		// 	Firstname: result.FirstName,
		// 	Lastname:  result.LastName,
		// 	Course:    tempCourse,
		// 	Project:   tempProject,
		// 	Username:  result.Username,
		// }

		json, err := json.Marshal(struct {
			Result []Course `json:"courses"`
		}{
			allCourses,
		})

		if err != nil {
			fmt.Println(err)
		}

		w.Write(json)
	} else {
		// // res.Error = err.Error()
		// json.NewEncoder(w).Encode(res)
		return
	}

}

func SearchRequiredUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var search Search
		_ = json.NewDecoder(r.Body).Decode(&search)
		fmt.Println(search)
		// currUsername := currUser.Username
		file, err := os.Open("blockchainFile.json")
		if err != nil {
			log.Println("Can't read file")
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		decoder.Token()
		block := Block{}
		// Appends decoded object to dataArr until every object gets parsed
		var users []Result
		for decoder.More() {
			decoder.Decode(&block)
			if block.Course.Name == search.CourseName || block.Course.Grade == search.CourseName {
				if block.Username != "" {
					resul := Result{
						Username: block.Username,
						Course:   block.Course,
					}
					users = append(users, resul)
				}
			}
		}
		// ReadBlockchainFile()
		// tempHead:=chainHead
		// for tempHead != nil {
		// 	if tempHead.Username == result.Username {
		// 		if tempHead.Course.Code != "" {
		// 			tempCourse = append(tempCourse, tempHead.Course)
		// 		}
		// 		if tempHead.Project.Name != "" {
		// 			tempProject = append(tempProject, tempHead.Project)
		// 		}
		// 	}
		// 	tempHead = tempHead.PrevPointer
		// }
		// cv := CV{
		// 	Email:     result.Email,
		// 	Firstname: result.FirstName,
		// 	Lastname:  result.LastName,
		// 	Course:    tempCourse,
		// 	Project:   tempProject,
		// 	Username:  result.Username,
		// }

		json, err := json.Marshal(struct {
			Result []Result `json:"users"`
		}{
			users,
		})

		if err != nil {
			fmt.Println(err)
		}

		w.Write(json)
	} else {
		// // res.Error = err.Error()
		// json.NewEncoder(w).Encode(res)
		return
	}

}

//AddCourseHandler Web Handler to add courses into the blockchain ///
func AddCourseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tempCourse CourseWeb
		_ = json.NewDecoder(r.Body).Decode(&tempCourse)
		fmt.Println(tempCourse)
		currUsername := currUser.Username

		var newCourse Course
		newCourse.Code = tempCourse.Code
		newCourse.Name = tempCourse.Name
		newCourse.CreditHours = tempCourse.CreditHours
		newCourse.Grade = tempCourse.Grade

		MyBlock := Block{
			Course:   newCourse,
			Email:    tempCourse.VEmail,
			Username: currUsername,
			Status:   "Pending",
		}

		//chainHead = InsertCourse(MyBlock)

		// gobEncoder := gob.NewEncoder(Satoshiconn)
		// err2 := gobEncoder.Encode(MyBlock)
		// if err2 != nil {
		// 	log.Println("In Write Chain: ", err2)
		// }

		fmt.Println(MyBlock)
		unverifiedChain = InsertCourseUnverified(MyBlock)
		ListBlocks(chainHead)

		uv := GetBlockhainArray(unverifiedChain)
		WriteUVFile(uv)

		//	fmt.Println("FFFFFFFFFF", len(nodesSlice))
		//		for i := 0; i < len(nodesSlice); i++ {
		//	fmt.Println("dddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
		//		if nodesSlice[i].Mail == MyBlock.Email {
		// conn, err := net.Dial("tcp", "localhost:"+nodesSlice[i].ListeningAddress)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// MinerConn = conn
		// gobEncoder := gob.NewEncoder(conn)
		// fmt.Println("blok:ahsh: ", CalculateHash(&MyBlock))
		// err2 := gobEncoder.Encode(MyBlock)
		// if err2 != nil {
		// 	log.Println("In Write Chain: ", err2)
		// }
		//SMPTPP
		// from := tempCourse.SEmail
		// password := tempCourse.SPass
		//
		// // Receiver email address.
		// to := []string{
		// 	MyBlock.Email,
		// }
		//
		// // smtp server configuration.
		// smtpHost := "smtp.gmail.com"
		// smtpPort := "587"
		//
		// // Message.
		// message := []byte("Course Name: " + MyBlock.Course.Name + "  Course Code: " + MyBlock.Course.Code + "  Course Grade: " + MyBlock.Course.Grade + "\n" + "Click here to verify this content: " + "localhost:" + "4000" + "/mineBlock/" + CalculateHash(&MyBlock))
		//
		// // Authentication.
		// auth := smtp.PlainAuth("", from, password, smtpHost)
		//
		// // Sending email.
		// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
		//
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// fmt.Println("Email Sent Successfully!")
		//SMPTPP

		m := gomail.NewMessage()

		// Set E-Mail sender
		m.SetHeader("From", tempCourse.SEmail)

		// Set E-Mail receivers
		m.SetHeader("To", MyBlock.Email)

		// Set E-Mail subject
		m.SetHeader("Subject", "Verification Content")

		// Set E-Mail body. You can set plain text or html with text/html
		m.SetBody("text/plain", "Course Name: "+MyBlock.Course.Name+"  Course Code: "+MyBlock.Course.Code+"  Course Grade: "+MyBlock.Course.Grade+"\n"+"Click here to verify this content: "+"localhost:"+"4000"+"/mineBlock/"+CalculateHash(&MyBlock))

		// Settings for SMTP server
		d := gomail.NewDialer("smtp.gmail.com", 587, tempCourse.SEmail, tempCourse.SPass)

		// This is only needed when SSL/TLS certificate is not valid on server.
		// In production this should be set to false.
		d.TLSConfig = &tls.Config{InsecureSkipVerify: false}

		// Now send E-Mail
		if err := d.DialAndSend(m); err != nil {
			fmt.Println(err, "mailerr")
			panic(err)
		}
		Mined = true
		//	fmt.Println("Email Sent", Mined, nodesSlice[i].ListeningAddress)

		//		break
		//		}
		//		}
		//http.Redirect(w, r, urlLogin+"/dashboard", http.StatusSeeOther)
	}

}

//LogoutHandler   Web Handler to logout the user ///
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	tokenString = ""
	http.Redirect(w, r, urlLogin+"/login", http.StatusSeeOther)
}

//Index  Web Handler to show homepage ///
func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../Website/homepage.html") //parse the html file homepage.html
	if err != nil {                                           // if there is an error
		log.Print("template parsing error: ", err, t) // log it
	}
	err = t.Execute(w, nil) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {         // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

//Details   Web Handler to show details ///
func Details(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../Website/discover.html") //parse the html file homepage.html
	if err != nil {                                           // if there is an error
		log.Print("template parsing error: ", err, t) // log it
	}
	err = t.Execute(w, nil) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {         // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

//--- serverHandler (Serving Angular files as static for deployment) ---//
var folderDist = "./public"

func serverHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(folderDist + r.URL.Path); err != nil {
		http.ServeFile(w, r, folderDist+"/index.html")
		return
	}
	fmt.Println((r.URL.Path))
	http.ServeFile(w, r, folderDist+r.URL.Path)
}

//----------------------------------------------------------------------//

//--- RunWebServer (Running Web Server) ---//
func RunWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/addProjectUser", AddProjectHandler)
	r.HandleFunc("/getBlocksUser", UnverifiedBlocksHandler)
	r.HandleFunc("/generateCVUser", GenerateCVHandler)
	r.HandleFunc("/addCourseUser", AddCourseHandler)
	r.HandleFunc("/registerUser", RegisterHandler)
	r.HandleFunc("/loginUser", LoginHandler)
	r.HandleFunc("/logout", LogoutHandler)
	r.HandleFunc("/getUser", ProfileHandler)
	r.HandleFunc("/getAllUsers", GetAllUsers)
	r.HandleFunc("/getVerifyContent", SearchVerifyContent)
	r.HandleFunc("/getVerifiedCVs", SearchRequiredUsers)
	r.HandleFunc("/mineBlockMiner/{hash}", Mineblock)

	r.NotFoundHandler = r.NewRoute().HandlerFunc(serverHandler).GetHandler()
	webPort := os.Getenv("PORT")
	fmt.Println(webPort)
	http.Handle("/", r)
	http.ListenAndServe(":"+webPort, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))
}
