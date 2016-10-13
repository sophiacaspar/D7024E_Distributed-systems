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

func TestDHT1(t *testing.T) {
/*
	id0 := "00"
    id1 := "01"
    id2 := "02"
    id3 := "03"
    id4 := "04"
    id5 := "05"
    id6 := "06"
    id7 := "07"


    node0 := startNode(&id0, "1110")
    node1 := startNode(&id1, "1111")
    node2 := startNode(&id2, "1112")
    node3 := startNode(&id3, "1113")
    node4 := startNode(&id4, "1114")
    node5 := startNode(&id5, "1115")
    node6 := startNode(&id6, "1116")
    node7 := startNode(&id7, "1117")
*/

 
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

	fmt.Print("")
	time.Sleep(10000*time.Millisecond)
	
	//node3.printMyFingers()
	//fmt.Println("#####################", node3.responsible("bf06670af35ed4abcadd95abe8079568f4df38e6"), "#####################")
	node3.kill()

	time.Sleep(10000*time.Millisecond)
	fmt.Println("hej")

	
	//node0.kill()
	//time.Sleep(10000*time.Millisecond)
	//msg := createPrintMsg(node7.transport.bindAddress, node1.transport.bindAddress)
	//go func () { node2.createTask("printRing", createPrintMsg(node2.transport.bindAddress, node3.transport.bindAddress))}()

	
/*	fmt.Println("#####################", node3.responsible(bf06670af35ed4abcadd95abe8079568f4df38e6), "#####################")
*/
/*
	time.Sleep(10000*time.Millisecond)
	msg := createPrintFingerMsg(node4.transport.bindAddress, node5.transport.bindAddress)
	go func () { node4.transport.send(msg)}() 

*/



	//time.Sleep(10000*time.Millisecond)
	//node2.lookup("04")
 


/*
	time.Sleep(7000*time.Millisecond)
	node1.lookupReq("lookup", "10", node5)
/*
	//msg := createPrintMsg(node2.transport.bindAddress, node3.transport.bindAddress)
	//go func () { node1.transport.send(msg)}() 
*/

	//node0.transport.listen()
	time.Sleep(2000*time.Second)
}