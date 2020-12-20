package main

import (
	"os"

	b "./blockchain"
)

func main() {

	satoshiAddress := os.Args[1]

	go b.RunWebServerSatoshi()
	go b.StartListening(satoshiAddress, "satoshi")
	//
	// course := b.Course{
	// 	Code:        "Ccc",
	// 	Name:        "sss",
	// 	CreditHours: 8,
	// }
	//
	// newBlock := b.Block{
	// 	//Hash here
	// 	Course: course,
	// }
	// course2 := b.Course{
	// 	Code:        "wdx",
	// 	Name:        "wqw2",
	// 	CreditHours: 8,
	// }
	// course3 := b.Course{
	// 	Code:        "999wsd",
	// 	Name:        "Kalb",
	// 	CreditHours: 8,
	// }
	// newBlock2 := b.Block{
	// 	//Hash here
	// 	Course: course2,
	// }
	// newBlock3 := b.Block{
	// 	//Hash here
	// 	Course: course3,
	// }
	//
	// chainHead := b.InsertCourse(newBlock)
	// chainHead = b.InsertCourse(newBlock2)
	// chainHead = b.InsertCourse(newBlock3)
	// b.ListBlocks(chainHead)
	// data := b.GetBlockhainArray(chainHead)
	//
	// b.WriteBlockchainFile(data)

	//	b.SendChainandConnInfo()

	select {}

}
