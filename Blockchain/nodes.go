package main

import (
	b "./blockchain"
)

func main() {

	// satoshiAddress := os.Args[1]
	// myListeningAddress := os.Args[2]
	//webAddress := "3334"
	//
	// conn, err := net.Dial("tcp", "localhost:"+satoshiAddress)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	go b.RunWebServer() //Runs webserver
	// go b.StartListening(myListeningAddress, "others") //Starts own server
	//
	// log.Println("Sending my listening address to Satoshis")
	// chainHead := b.ReceiveChain(conn)
	// b.ListBlocks(chainHead)
	//
	// Peers := b.Client{
	// 	ListeningAddress: myListeningAddress,
	// 	Types:            true,
	// }
	// b.WriteString(conn, Peers) //Writes its address
	//
	// //go b.ReceiveChain(conn)
	//
	// go b.ReadPeersMinerChainEverything(conn) // Reads information from Satoshi every second

	// go func() { //Go routine for reading the chain that miner sends
	// 	for {
	// 		if b.Mined == true { // checks if the block sent is mined or not
	// 			fmt.Println("trueue")
	// 			var stuu b.Combo
	// 			fmt.Println("In Read Peers fffwd")
	// 			gobEncoder := gob.NewDecoder(b.MinerConn)
	// 			err := gobEncoder.Decode(&stuu)
	// 			if err != nil {
	// 				log.Println(err, "FFF")
	// 			}
	// 			fmt.Println("Read StuuPeers: ", stuu.ClientsSlice)
	// 			b.ListBlocks(stuu.ChainHead)
	// 			// if Length(stuu.ChainHead) >= Length(chainHead) {
	// 			// 	chainHead = stuu.ChainHead
	// 			// 	stuff.ChainHead = chainHead
	// 			// 	fmt.Println("Read Chain: ")
	// 			// 	ListBlocks(chainHead)
	// 			// }
	// 			b.Mined = false
	// 		}
	//
	// 	}
	//
	// }()

	select {}

}
