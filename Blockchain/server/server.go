package server

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)
type ToDoList struct {
	Task   string             `json:"task,omitempty"`
	Status bool               `json:"status,omitempty"`
}
type Skill struct {
}
type Course struct {
	Code        string `json:"cCode"`
	Name        string `json:"cName"`
	CreditHours string `json:"cHrs"`
	Grade       string `json:"cGrade"`
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
	BlockNo 		int `json:"blockno"`
	Course      Course `json:"course"`
	project     Project
	PrevPointer *Block `json:"prevPoint"`
	PrevHash    string `json:"prevHash"`
	CurrentHash string `json:"currHash"`
}

//var chainHead *Block
var LocalData []Connected
var GlobalData Data
var mutex = &sync.Mutex{}

// For Creating and Displaying for now ///

var chainHead *Block
var count int = 0

/////////////////

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

var Addchan = make(chan Peer)
var globe Data
var stopchan = make(chan string)

//var clientsSlice []Verifier
func BroadcastAdminData() {

	for i := 0; i < len(LocalData); i++ {
		gobEncoder := gob.NewEncoder(LocalData[i].Conn)
		//fmt.Println("BroadCheck: ", LocalData[i])
		err1 := gobEncoder.Encode(GlobalData)
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

	//	fmt.Println("Slice:", GlobalData.clientsSlice[0].ListeningAddress)
	Addchan <- newClient
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
		if Length(chainhead) > Length(GlobalData.ChainHead) {
			GlobalData.ChainHead = chainhead
		}
		fmt.Println("Blockchain received:: ", Length(GlobalData.ChainHead))
		//	GlobalData = globe
		//	<-Globechan
		//<-RW3Chan
	}
}

func broadcastBlock() {
	StepbyChan <- "Hello"
	//	RW3Chan <- "hello"
	time.Sleep(5 * time.Second)
	for i := 0; i < len(LocalData); i++ {
		gobEncoder := gob.NewEncoder(LocalData[i].Conn)
		//fmt.Println("BroadCheck: ", LocalData[i])
		err1 := gobEncoder.Encode(GlobalData.ChainHead)
		fmt.Println("Broadcasting Blockchain:: ")
		if err1 != nil {
			log.Println(err1, "in broadcasting block")
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
		InsertOnlyBlock(recvdBlock, GlobalData.ChainHead)
	}
}

func broadcastBlockchaintoPeer(conn net.Conn) {
	//channel
	gobEncoder := gob.NewEncoder(conn)
	err1 := gobEncoder.Encode(GlobalData.ChainHead)
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
func ReadAdminData(conn net.Conn) {
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
	if Length(globe.ChainHead) < Length(GlobalData.ChainHead) {
		globe.ChainHead = GlobalData.ChainHead
	}
	if len(globe.ClientsSlice) < len(GlobalData.ClientsSlice) {
		globe.ClientsSlice = GlobalData.ClientsSlice
	}
	fmt.Println("Blockchain read:")
	ListBlocks(globe.ChainHead)

	GlobalData.ChainHead = globe.ChainHead
}

func ViewMinerData() {
	for i := 0; i < len(GlobalData.ClientsSlice); i++ {
		if GlobalData.ClientsSlice[i].Role == "miner" {
			fmt.Println("Miners connected to system:")
			fmt.Print(" Their address: ", GlobalData.ClientsSlice[i].ListeningAddress)
		}
	}
}

// CreateTask create task route
func CreateBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)
	fmt.Println(course, r.Body)
	json.NewEncoder(w).Encode(course)

	chainHead = InsertCourse(course, chainHead)

	ListBlocks(chainHead)

}

// CreateTask create task route
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var task2 ToDoList
	_ = json.NewDecoder(r.Body).Decode(&task2)
	fmt.Println(task2, r.Body)
	json.NewEncoder(w).Encode(task2)
}


func GetAllBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getCourse(chainHead)
	json.NewEncoder(w).Encode(payload)
}

func getCourse (chainHead *Block) []Block {
	var courses []Block
		for chainHead != nil {
			courses = append(courses,*chainHead)
			chainHead = chainHead.PrevPointer
		}
		fmt.Println()
		return courses
}
