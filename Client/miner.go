//
//
// func main(){
//
// satoshiAddress := os.Args[1]
// myListeningAddress := os.Args[2]
//
// conn, err := net.Dial("tcp", satoshiAddress)
// if err != nil {
//   log.Fatal(err)
// }
// //The function below launches the server, uses different second argument
// //It then starts a routine for each connection request received
// //go StartListening(myListeningAddress, "others")
//
// log.Println("Sending my listening address to Satoshi")
// //Satoshi is there waiting for our address, it stores it somehow
// WriteString(conn, myListeningAddress)
//
// }
