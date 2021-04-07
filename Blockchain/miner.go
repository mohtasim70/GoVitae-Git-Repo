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
	mail := os.Args[4]

	conn, err := net.Dial("tcp", ":"+satoshiAddress)
	if err != nil {
		log.Fatal(err)
	}

	b.Doit = true
	go b.RunWebServerMiner(webAddress)               //Starts own web server
	go b.StartListening(myListeningAddress, "miner") //Starts its own server

	log.Println("Sending my listening address to Satoshi")

	b.ReceiveChain(conn) //Receives chain the first time from Satoshi

	Peers := b.Client{
		ListeningAddress: myListeningAddress,
		Types:            false,
		Mail:             mail,
	}
	b.WriteString(conn, Peers) //Sends its info including his mail to Satoshi

	go b.ReadPeersMinerChainEverything(conn) //Reads info from Satoshi

	//go b.ReadPeersMinerChain(conn)

	// go func() {
	// 	for {
	// 		fmt.Println("Innn Miner Sent:: ")
	// 		b.ReceiveChain(conn)
	//
	// 	}
	// }()

	// slice := b.ReadPeers(conn)
	// fmt.Println("Slice:: ", slice)
	select {}

}
