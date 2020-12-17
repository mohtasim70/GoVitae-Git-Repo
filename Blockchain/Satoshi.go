package main

import (
	"os"

	b "./blockchain"
)

func main() {

	satoshiAddress := os.Args[1]

	go b.RunWebServerSatoshi()
	go b.StartListening(satoshiAddress, "satoshi")

	//	b.SendChainandConnInfo()

	select {}

}
