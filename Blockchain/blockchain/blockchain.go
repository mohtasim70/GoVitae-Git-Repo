package blockchain

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	//"github.com/gorilla/websocket"
)

var chainHead *Block

type Skill struct {
}
type Course struct {
	Code        string
	Name        string
	CreditHours int
	Grade       string
}
type Project struct {
	Name     string
	Document string
	Course   Course
}

type Block struct {
	Course      Course
	Project     Project
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
	BlockNo     int
	Status      string
	Email       string
}

type ListTheBlock struct {
	Course      []Course
	Project     []Project
	PrevPointer []*Block
	PrevHash    []string
	CurrentHash []string
	BlockNo     []int
	Status      []string
	Email       []string
}

type Client struct {
	ListeningAddress string
	Types            bool //true for node and false for miner
	Mail             string
}
type Connected struct {
	Conn net.Conn
}

var count int = 0

var localData []Connected

//var nodes = make(map[*websocket.Conn]bool) // connected clients
//var upgrader = websocket.Upgrader{
//CheckOrigin: func(r *http.Request) bool {
//return true
//},
//}

//256bit
func CalculateHash(inputBlock *Block) string {

	var temp string
	temp = inputBlock.Course.Code + inputBlock.Project.Name
	h := sha256.New()
	h.Write([]byte(temp))
	sum := hex.EncodeToString(h.Sum(nil))

	// sum := sha256.Sum256([]byte(temp))

	return sum
}
func InsertBlock(course Course, project Project, chainHead *Block) *Block {
	newBlock := &Block{
		//Hash here
		Course:  course,
		Project: project,
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
		Project: project,
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

// Changing InsertCourse Code //
func InsertCourse(myBlock Block) *Block {

	myBlock.CurrentHash = CalculateHash(&myBlock)

	if chainHead == nil {
		myBlock.BlockNo = count
		myBlock.PrevHash = "null"
		myBlock.Status = "Pending"
		chainHead = &myBlock
		fmt.Println("Genesis Block Inserted")
		return chainHead
	}
	count = count + 1
	myBlock.PrevPointer = chainHead
	myBlock.PrevHash = chainHead.CurrentHash
	myBlock.BlockNo = count
	myBlock.Status = "Pending"

	fmt.Println("Course Block Inserted")
	return &myBlock

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
func WriteString(conn net.Conn, myListeningAddress Client) {
	Satoshiconn = conn
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(myListeningAddress)
	if err != nil {
		log.Println("In Write String: ", err)
	}
}

func SendChain(conn net.Conn) {
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(chainHead)
	if err != nil {
		log.Println("In Write Chain: ", err)
	}
}

var Satoshiconn net.Conn
var clientsSlice []Client
var rwchan = make(chan string)

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
	addchan <- Clientz
	//WaitForQuorum()

}

var nodesSlice []Client
var minechan = make(chan Client)

var blockchan = make(chan Block)
var Minedblock Block

func handlePeer(conn net.Conn) {

	//	Clientz := Client{}
	block := Block{}
	//var ListeningAddress string
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&block)
	if err != nil {
		//handle error
	}

	// newClient.ListeningAddress = ListeningAddress
	fmt.Println("inHandlePeer: ", block.Email)
	blockchan <- block

}
func ReceiveChain(conn net.Conn) Block {
	var block Block
	gobEncoder := gob.NewDecoder(conn)
	err := gobEncoder.Decode(&block)
	if err != nil {
		log.Println(err)
	}
	chainHead = InsertCourse(block)
	return block
}

func broadcastPeerData() {

	for i := 0; i < len(localData); i++ {
		gobEncoder := gob.NewEncoder(localData[i].Conn)
		//fmt.Println("BroadCheck: ", localData[i])
		err1 := gobEncoder.Encode(clientsSlice)
		fmt.Println("Broadcasting StreamData:: ")
		if err1 != nil {
			log.Println("Errpr in broadcasting", err1)
		}

	}
	//	<-StepbyChan

}

func ReadPeers(conn net.Conn) []Client {
	var slice []Client
	gobEncoder := gob.NewDecoder(conn)
	err := gobEncoder.Decode(&slice)
	if err != nil {
		log.Println(err)
	}
	nodesSlice = slice
	return nodesSlice
}
func ReadBlockPeers(conn net.Conn) Block {
	var block Block
	gobEncoder := gob.NewDecoder(conn)
	err := gobEncoder.Decode(&block)
	if err != nil {
		log.Println(err)
	}
	return block
}

func StartListening(ListeningAddress string, node string) {
	if node == "satoshi" {
		ln, err := net.Listen("tcp", "localhost:"+ListeningAddress)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Faital")
		}

		addchan := make(chan Client)
		block := Block{}
		chainHead = InsertCourse(block) //Genesis Block
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err, "Yooooo")
				continue
			}
			sendBlockchain(conn, chainHead)
			conns := Connected{
				Conn: conn,
			}

			go handleConnection(conn, addchan)
			clientsSlice = append(clientsSlice, <-addchan)
			localData = append(localData, conns)
			broadcastPeerData()
			//	chainHead = a2.InsertBlock("", "", "Satoshi", 0, chainHead)
			// var block Block
			// gobEncoder := gob.NewDecoder(conn)
			// err2 := gobEncoder.Decode(&block)
			// if err2 != nil {
			// 	log.Println(err2)
			// }
			go ReceiveChain(conn)
			//chainHead = InsertCourse(block)
			ListBlocks(chainHead)

		}
	} else if node == "others" {
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
			go handlePeer(conn)
			nodesSlice = append(nodesSlice, <-minechan)

		}

	} else { //miner
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
			conns := Connected{
				Conn: conn,
			}
			localData = append(localData, conns)

			go handlePeer(conn)

			Minedblock = <-blockchan

		}
	}

}

// Chi HTTP Services //

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// handle returned error here.
		w.WriteHeader(503)
		w.Write([]byte("bad"))
	}
}

func setHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("../Website/blockchain.html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, nil) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {         // if there is an error
		log.Print("template executing error: ", err) //log it
	}
	return nil
}

func getHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	cCode := r.Form.Get("courseCode")
	cName := r.Form.Get("courseName")
	cGrade := r.Form.Get("courseGrade")
	cEmail := r.Form.Get("courseEmail")

	a, err := strconv.Atoi(r.FormValue("courseCHrs"))
	if err != nil {
	}
	cCHrs := a

	AddCourse := Course{
		Code:        cCode,
		Name:        cName,
		CreditHours: cCHrs,
		Grade:       cGrade,
	}

	MyBlock := Block{
		Course: AddCourse,
		Email:  cEmail,
	}

	//chainHead = InsertCourse(MyBlock)

	// gobEncoder := gob.NewEncoder(Satoshiconn)
	// err2 := gobEncoder.Encode(MyBlock)
	// if err2 != nil {
	// 	log.Println("In Write Chain: ", err2)
	// }

	ListBlocks(chainHead)

	tempHead := chainHead
	viewTheBlock := new(ListTheBlock)
	tempCourse := []Course{}
	tempBlockNo := []int{}
	tempCurrHash := []string{}
	tempPrevHash := []string{}
	tempEmail := []string{}
	for tempHead != nil {
		tempCourse = append(tempCourse, tempHead.Course)
		tempBlockNo = append(tempBlockNo, tempHead.BlockNo)
		tempCurrHash = append(tempCurrHash, tempHead.CurrentHash)
		tempPrevHash = append(tempPrevHash, tempHead.PrevHash)
		tempEmail = append(tempEmail, tempHead.Email)
		viewTheBlock = &ListTheBlock{
			Course:      tempCourse,
			BlockNo:     tempBlockNo,
			CurrentHash: tempCurrHash,
			PrevHash:    tempPrevHash,
			Email:       tempEmail,
		}
		tempHead = tempHead.PrevPointer
		fmt.Println(viewTheBlock.Course)
		fmt.Println(viewTheBlock.BlockNo)
		fmt.Println(viewTheBlock.CurrentHash)
		fmt.Println(viewTheBlock.PrevHash)
		fmt.Println(viewTheBlock.Email)
	}
	// generate page by passing page variables into template
	t, err := template.ParseFiles("../Website/blockchain.html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, viewTheBlock) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}

	for i := 0; i < len(nodesSlice); i++ {
		if nodesSlice[i].Mail == MyBlock.Email {
			conn, err := net.Dial("tcp", "localhost:"+nodesSlice[i].ListeningAddress)
			if err != nil {
				log.Fatal(err)
			}
			gobEncoder := gob.NewEncoder(conn)
			err2 := gobEncoder.Encode(MyBlock)
			if err2 != nil {
				log.Println("In Write Chain: ", err2)
			}

			break
		}
	}

	// gobEncoder := gob.NewEncoder(Satoshiconn)
	// err2 := gobEncoder.Encode(MyBlock)
	// if err2 != nil {
	// 	log.Println("In Write Chain: ", err2)
	// }

	return nil
}

func showBlocksHandler(w http.ResponseWriter, r *http.Request) error {
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

	return nil
}

func Mineblock(w http.ResponseWriter, r *http.Request) error {
	chainHead = InsertCourse(Minedblock)
	fmt.Println("In Mining")
	ListBlocks(chainHead)

	return nil
}

var broadcast = make(chan []Block) // broadcast channel

//func HandleConnections(w http.ResponseWriter, r *http.Request) {

//	ws, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Fatal(err)
//		fmt.Println("Error in ebss")
//	}

// make sure we close the connection when the function returns
//	defer ws.Close()

// register our new client
//	nodes[ws] = true

//	for {
// Read in a new message as JSON and map it to a Message object
//		var course Course
//		err := json.NewDecoder(r.Body).Decode(&course)
//		if err != nil {
//			panic(err)
//		}
// err := ws.ReadJSON(&course)
//		chainHead = InsertCourse1(course, chainHead)
// if err != nil {
// 	log.Printf("error: %v", err)
// 	//	delete(nodes, ws)
// 	break
// }

// Send the newly received message to the broadcast channel
//		broadcast <- getCourse(chainHead)
//	}

//}

// Clients Web Server //

func RunWebServer(port string) {
	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/ws", server.HandleConnections)
	// router.HandleFunc("/api/block", server.GetAllBlock).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/block", server.CreateBlock).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/task", server.CreateTask).Methods("POST", "OPTIONS")

	r := chi.NewRouter()
	r.Method("GET", "/", Handler(setHandler))
	r.Method("POST", "/blockInsert", Handler(getHandler))
	//r.HandleFunc("/ws", HandleConnections)

	http.ListenAndServe("localhost:"+port, r)

}

func RunWebServerMiner(port string) {
	r := chi.NewRouter()
	r.Method("GET", "/mine", Handler(Mineblock))
	// r.Method("POST", "/blockInsert", Handler(getHandler))
	//r.HandleFunc("/ws", HandleConnections)

	http.ListenAndServe("localhost:"+port, r)

}

// Satoshi Web Server //

func RunWebServerSatoshi() {
	r := chi.NewRouter()
	r.Method("GET", "/showBlocks", Handler(showBlocksHandler))
	//r.HandleFunc("/ws", HandleConnections)

	http.ListenAndServe("localhost"+":3333", r)

}

//func BroadcastMessages() {
//	for {
//	// Grab the next message from the broadcast channel
//	msg := <-broadcast
//	fmt.Println("In broadcast: ", msg)
// Send it out to every client that is currently connected
//	for client := range nodes {
//		err := client.WriteJSON(msg)
//		if err != nil {
//		log.Printf("error: %v", err)
//		client.Close()
//		delete(nodes, client)
//		}
//	}
//	}
//}

// ---- //

func main() {
	// ln, err := net.Listen("tcp", "localhost:6003")
	// if err != nil {
	//
	// 	log.Fatal(err, ln)
	//
	// }
	//go RunWebServer()

	//go BroadcastMessages()

	//select {}

	// conn, err := net.Dial("tcp", "localhost:3333")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("ss", conn)
	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		log.Println(err)
	// 		continue
	// 	}
	// 	go sendBlockchain(conn, chainHead)
	// }

}
