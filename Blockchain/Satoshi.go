package main

import (
	"os"

	b "./blockchain"
)

func main() {

	satoshiAddress := os.Args[1]

	go b.RunWebServerSatoshi()                     //Own web server
	go b.StartListening(satoshiAddress, "satoshi") //Listens to Clients(Nodes and Miners)

	select {}

}
