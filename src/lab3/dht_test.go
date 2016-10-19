package dht

/** go test -test.run TestDHT1 */

import (
	"fmt"
	"testing"
	"time"
)

// dhtNode sends request to master of ring: please add me somewhere
func (dhtNode *DHTNode) joinReq(master *DHTNode) {
	msg := createJoinMsg(master.transport.bindAddress, [2]string{dhtNode.transport.bindAddress, dhtNode.nodeId})
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
		dhtNode.startServer()
		dhtNode.joinReq(master)

		msg := createGetBackupMsg(dhtNode.transport.bindAddress, dhtNode.successor[0])
		go func () { dhtNode.transport.send(msg)}() 	

	} 
}

func TestDHT1(t *testing.T) {
 
	node0 := startNode(nil, "1110")
    node1 := startNode(nil, "1111")
    node2 := startNode(nil, "1112")
    node3 := startNode(nil, "1113")
    node4 := startNode(nil, "1114")
    node5 := startNode(nil, "1115")
    node6 := startNode(nil, "1116")
    node7 := startNode(nil, "1117")

	node6.joinReq(node1)
	node5.joinReq(node1)
	node7.joinReq(node1)
	node3.joinReq(node1)
	node2.joinReq(node1)
	node0.joinReq(node1)
	node4.joinReq(node1)

	time.Sleep(10000*time.Millisecond)

	/*
	msg := createUploadMsg(node4.transport.bindAddress, node3.transport.bindAddress, data)
	go func () { node4.transport.send(msg)}() 
	*/

	//node3.printMyFingers()
	//fmt.Println("#####################", node3.responsible("bf06670af35ed4abcadd95abe8079568f4df38e6"), "#####################")

	//node3.kill()

	time.Sleep(10000*time.Millisecond)

    //node3.alive()
	
	time.Sleep(2000*time.Second)

}