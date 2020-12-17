package main

import (
	"encoding/gob"
	"fmt"
	"net"
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

type Block struct {
	course      Course
	project     Project
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

func main() {

	conn, err := net.Dial("tcp", "localhost:6003")
	if err != nil {
		//handle error
	}
	var recvdBlock Block
	dec := gob.NewDecoder(conn)
	err = dec.Decode(&recvdBlock)
	if err != nil {
		//handle error
	}
	fmt.Println(recvdBlock.CurrentHash)

}
