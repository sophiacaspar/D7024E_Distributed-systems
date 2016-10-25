package dht

/** Nodernas ordning: 4, 5, 2, 3, 7, 0, 6, 1   */

import (
	"fmt"
	"testing"
	"time"
)

// dhtNode sends request to master of ring: please add me somewhere
func (dhtNode *DHTNode) joinReq(masterIP string) {
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
	node.startServer()
	time.Sleep(300*time.Millisecond)
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
		dhtNode.joinReq(master)
		//dhtNode.startServer()
	} 
}

func StartUp(port string) {
	bootstrapnode := "localhost:1110"
	ip := "localhost:" + port
	node := startNode(nil, port)

	if ip != bootstrapnode {
		node.dht.JoinReq(bootstrapnode)
	}
}