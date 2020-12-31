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

	conn, err := net.Dial("tcp", "localhost:"+satoshiAddress)
	if err != nil {
		log.Fatal(err)
	}
	b.Doit = true
	go b.RunWebServerMiner(webAddress)
	go b.StartListening(myListeningAddress, "miner")

	log.Println("Sending my listening address to Satoshi")

	b.ReceiveChain(conn)

	Peers := b.Client{
		ListeningAddress: myListeningAddress,
		Types:            false,
		Mail:             mail,
	}
	b.WriteString(conn, Peers)

	go b.ReadPeersMinerChainEverything(conn)

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
