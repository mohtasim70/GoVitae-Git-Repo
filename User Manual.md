Implementation:

Part 0 - Genesis:
It’s a peer to peer system with miners, users as its main system users. A node ( Satoshi Node ) is created at the time of Genesis/Start of everything. Now this Satoshi node has a mongoDB backend and the complete updated blockchain is stored here as a json file.
System itself will run on the Satoshi node

Part 1 - User Side:
User using command prompt will start its TCP server. Now user is taken to the User side of the website. User is authenticated into the system using JWT.
After authentication user is now in the Portal. Here user can select and enter courses, skills, projects that he wants to add to his CV and verified in BlockChain. After adding required details user can then send this block to the respective miner using email functionality of GoVitae which is implemented using GOMAIL package.
User is sending this “Verification Request Email” to miner so that miner can opt to verify or not verify the block and its details

Part 2 - Miner Side:
Miner will have received a “Verification Request Email” by the User already. This email includes
1.	The user details who have sent this mail
2.	The block content ( courses, hash, skills, projects etc)
3.	The HashMine link ( it’s a link which the user will select in order to verify that block. This link is based on the Hash of the block )
Now the miner can select HashMine Link. Upon selecting this link, miner will be redirected and now notified that the block in now verified and part of the blockchain.

Part 3 - User Side:
Once Part 3 is done now the user can see on his dashboard about which of their contents are now verified and part of blockchain. User can also see which of the contents are pending to be verified.

The Blockchain:
The blockchain itself is continuously broadcasted so that all the nodes have the latest untampered copy of blockchain. In this way blockchain is distributed across all the nodes.

User Manual:
Setup:
Install Go Lang
Install Mongo DB
Open command prompt in project folder and run these commands:
go get github.com/dgrijalva/jwt-go
go get github.com/gorilla/mux
go get go.mongodb.org/mongo-driver/bson
go get golang.org/x/crypto/bcrypt
go get gopkg.in/mail.v2

Running:
Open 3 terminals in project folder and run these commands:
•	Terminal 1:
Go run satoshi.go 4001
•	Terminal 2:
Go run miner.go 4001 (any port No) 4002
•	Terminal 3:
Go run nodes.go 4001 (any port No) 4002

You can also start new terminals for more miners and nodes. Run same commands.
•	For new miner:
Go run miner.go 4001 (any port No) 4002
•	For new node:
Go run nodes.go 4001 (any port No) 4002

Now open browser and dial 
localhost:4002 

*Refer to implementation part to better understand how this runs.
