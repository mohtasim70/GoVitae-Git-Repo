package main

import (
	"fmt"
	"log"
	"net"
	"os"

	b "./blockchain"
)

func main() {

	satoshiAddress := os.Args[1]
	myListeningAddress := os.Args[2]
	webAddress := os.Args[3]

	conn, err := net.Dial("tcp", "localhost:"+satoshiAddress)
	if err != nil {
		log.Fatal(err)
	}

	go b.RunWebServer(webAddress)
	go b.StartListening(myListeningAddress, "others")

	log.Println("Sending my listening address to Satoshi")
	//Satoshi is there waiting for our address, it stores it somehow
	chainHead := b.ReceiveChain(conn)
	b.ListBlocks(chainHead)

	Peers := b.Client{
		ListeningAddress: myListeningAddress,
		Types:            true,
	}
	b.WriteString(conn, Peers)

	//go b.ReadPeers(conn)
	slice := b.ReadPeers(conn)
	fmt.Println("Slice:: ", slice)

	// chainHead = b.ReceiveChain(conn)

	//once the satoshi unblocks on Quorum completion it sends peer to connect to
	// log.Println("receiving peer to connect to ... ")
	// receivedString := b.ReadString(conn)
	// log.Println(receivedString)
	select {}

}
