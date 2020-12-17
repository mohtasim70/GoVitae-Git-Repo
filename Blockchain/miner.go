package main

import (
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

	go b.RunWebServerMiner(webAddress)
	go b.StartListening(myListeningAddress, "miner")

	log.Println("Sending my listening address to Satoshi")
	//Satoshi is there waiting for our address, it stores it somehow
	chainHead := b.ReceiveChain(conn)
	b.ListBlocks(chainHead)

	Peers := b.Client{
		ListeningAddress: myListeningAddress,
		Types:            false,
		Mail:             "mohtasimasad@gmail.com",
	}
	b.WriteString(conn, Peers)

	go b.ReadPeers1(conn)

	// slice := b.ReadPeers(conn)
	// fmt.Println("Slice:: ", slice)
	select {}

}
