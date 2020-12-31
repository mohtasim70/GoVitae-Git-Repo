package main

import (
	"encoding/gob"
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

	log.Println("Sending my listening address to Satoshis")
	chainHead := b.ReceiveChain(conn)
	b.ListBlocks(chainHead)

	Peers := b.Client{
		ListeningAddress: myListeningAddress,
		Types:            true,
	}
	b.WriteString(conn, Peers)

	//go b.ReceiveChain(conn)

	go b.ReadPeersMinerChainEverything(conn)

	go func() {
		for {
			if b.Mined == true {
				fmt.Println("trueue")
				var stuu b.Combo
				fmt.Println("In Read Peers fffwd")
				gobEncoder := gob.NewDecoder(b.MinerConn)
				err := gobEncoder.Decode(&stuu)
				if err != nil {
					log.Println(err, "FFF")
				}
				fmt.Println("Read StuuPeers: ", stuu.ClientsSlice)
				b.ListBlocks(stuu.ChainHead)
				// if Length(stuu.ChainHead) >= Length(chainHead) {
				// 	chainHead = stuu.ChainHead
				// 	stuff.ChainHead = chainHead
				// 	fmt.Println("Read Chain: ")
				// 	ListBlocks(chainHead)
				// }
				b.Mined = false
			}

		}

	}()

	select {}

}
