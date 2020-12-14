package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"./router"
	"./server"
)

func StartListening(listeningAddress string, node string) {

	if node == "server" {
		ln, err := net.Listen("tcp", "localhost:"+listeningAddress)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Faital")
		}
		fmt.Println("Stream Starts")
		newBlock := server.Block{} // Removed & //
		server.GlobalData.ChainHead = server.InsertOnlyBlock(&newBlock, server.GlobalData.ChainHead)

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			conns := server.Connected{
				Conn: conn,
			}
			// go broadcastBlockchaintoPeer(conn)
			// go receiveBlockchainfromPeer(conn)

			go server.StoreClient(conn)
			//	go readBlockchain(conn)
			go server.ReadAdminData(conn)

			server.GlobalData.ClientsSlice = append(server.GlobalData.ClientsSlice, <-server.Addchan)
			server.LocalData = append(server.LocalData, conns)

			go server.BroadcastAdminData()
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
		clientsSlice := make([]server.Peer, 10)
		//	addchan := make(chan Peer)
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			newClient := server.Peer{
				// Conn: conn,
			}
			clientsSlice = append(clientsSlice, newClient)
			// go broadcastBlockchaintoPeer(conn)
			// go receiveBlockchainfromPeer(conn)

			go server.MinerverifyBlock(conn)
		}

	}
}

func runWebServer(port string) {
	r := router.Router()
	fmt.Println("Starting server on the port: ",port)
	log.Fatal(http.ListenAndServe("localhost"+":"+port, r))
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
	//go StartListening(satoshiAddress, "server")
	go runWebServer(satoshiAddress)
	

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
