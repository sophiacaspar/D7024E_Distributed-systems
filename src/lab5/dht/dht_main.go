package dht

/** go test -test.run TestDHT1 */
/** Nodernas ordning: 4, 5, 2, 3, 7, 0, 6, 1   */

import (
	"fmt"
	"time"
	//"os"
	//"io/ioutil"
	//b64 "encoding/base64"
)

// starts up node network (used for docker)
func StartUpNetwork(port string) {

	bootstrapnode := "localhost:1110"
	//ip := "localhost:" + port
	node := startNode(nil, port)
	fmt.Println("hello")
	if port != "1110" {
		node.joinReq(bootstrapnode)
	}
	time.Sleep(300*time.Second)
}

// dhtNode sends request to master of ring: please add me somewhere
func (dhtNode *DHTNode) joinReq(masterIP string) {
	fmt.Print("creates join msg for", dhtNode.transport.bindAddress, "and", masterIP)
	msg := createJoinMsg(masterIP, [2]string{dhtNode.transport.bindAddress, dhtNode.nodeId})
	go func () { dhtNode.transport.send(msg)}() 
	time.Sleep(time.Millisecond*500)
}

// Master tells dhtNode to check who's responsible for key
func (master *DHTNode) lookupReq(t, key string, dhtNode *DHTNode) {
	m := createLookupMsg(t, master.transport.bindAddress, key, master.transport.bindAddress, dhtNode.transport.bindAddress)
 	go func () { dhtNode.transport.send(m)} () 	
}

func startNode(nodeId *string, port string) (dhtNode *DHTNode) {
	node := makeDHTNode(nodeId, "localhost", port)
	go node.startServer()
	return node
}

func (dhtNode *DHTNode) alive(master *DHTNode) {
	if dhtNode.online == false {
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<<",dhtNode.contact.port, "IS ALIVE <<<<<<<<<<<<<<<<<<<<<<")
		dhtNode.online = true
		go dhtNode.init_taskQueue()
		go dhtNode.stabilizeTimer()
		go dhtNode.fingerTimer()
		go dhtNode.heartbeatTimer()

		go dhtNode.transport.listen()
		//dhtNode.startServer()
		dhtNode.joinReq(master.transport.bindAddress)
	} 
}

