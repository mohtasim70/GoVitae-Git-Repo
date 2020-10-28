package main

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
}

type ListTheBlock struct {
	Course      []Course
	Project     []Project
	PrevPointer []*Block
	PrevHash    []string
	CurrentHash []string
}

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
		chainHead = &myBlock
		fmt.Println("Block Inserted")
		return chainHead
	}
	myBlock.PrevPointer = chainHead
	myBlock.PrevHash = chainHead.CurrentHash

	fmt.Println("Project Block Inserted")
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
		fmt.Print("Block-- ")
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
	t, err := template.ParseFiles("../Website/Blockchain.html") //parse the html file homepage.html
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
	}

	chainHead = InsertCourse(MyBlock)
	ListBlocks(chainHead)

	tempHead := chainHead
	viewTheBlock := new(ListTheBlock)
	tempCourse := []Course{}
	for tempHead != nil {
		tempCourse = append(tempCourse, tempHead.Course)
		viewTheBlock = &ListTheBlock{
			Course: tempCourse,
		}
		tempHead = tempHead.PrevPointer
		fmt.Println(viewTheBlock.Course)
	}
	// generate page by passing page variables into template
	t, err := template.ParseFiles("../Website/viewBlock.html") //parse the html file homepage.html
	if err != nil {                                            // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, viewTheBlock) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
	return nil
}

func runWebServer() {
	r := chi.NewRouter()
	r.Method("GET", "/", Handler(setHandler))
	r.Method("POST", "/blockInsert", Handler(getHandler))

	http.ListenAndServe(":3333", r)
}

// ---- //

func main() {
	ln, err := net.Listen("tcp", ":6003")
	if err != nil {

		log.Fatal(err)

	}
	go runWebServer()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go sendBlockchain(conn, chainHead)
	}

}
