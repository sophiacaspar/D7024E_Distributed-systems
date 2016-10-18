package dht

/** go test -test.run TestDHT1 */

import (
	"fmt"
	"testing"
	"time"
	"io/ioutil"
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

    node0 := makeDHTNode(&id0, "localhost", "1110")
    node1 := makeDHTNode(&id1, "localhost", "1111")
    node2 := makeDHTNode(&id2, "localhost", "1112")
    node3 := makeDHTNode(&id3, "localhost", "1113")
    node4 := makeDHTNode(&id4, "localhost", "1114")
    node5 := makeDHTNode(&id5, "localhost", "1115")
    node6 := makeDHTNode(&id6, "localhost", "1116")
    node7 := makeDHTNode(&id7, "localhost", "1117")
 

    node0 := makeDHTNode(nil, "localhost", "1110")
    node1 := makeDHTNode(nil, "localhost", "1111")
    node2 := makeDHTNode(nil, "localhost", "1112")
    node3 := makeDHTNode(nil, "localhost", "1113")
    node4 := makeDHTNode(nil, "localhost", "1114")
    node5 := makeDHTNode(nil, "localhost", "1115")
    node6 := makeDHTNode(nil, "localhost", "1116")
    node7 := makeDHTNode(nil, "localhost", "1117")


	node1.startServer()
	node2.startServer()
	node3.startServer()
	node4.startServer()
	node5.startServer()
	node6.startServer()
	node7.startServer()

	node7.joinReq(node1)
	node6.joinReq(node1)
	node5.joinReq(node1)
	node4.joinReq(node1)
	node3.joinReq(node1)
	node2.joinReq(node1)

*/

	dat, err := ioutil.ReadFile("file.txt")
    if err != nil {
		panic(err)
	}
	
    fmt.Print(string(dat))
    fmt.Print("\n")

/*
	time.Sleep(7000*time.Millisecond)
	node3.kill()
*/

	//time.Sleep(10000*time.Millisecond)
	//msg := createPrintFingerMsg(node4.transport.bindAddress, node5.transport.bindAddress)
	//go func () { node4.transport.send(msg)}() 





	//time.Sleep(10000*time.Millisecond)
	//node2.lookup("04")
	//msg := createPrintMsg(node7.transport.bindAddress, node1.transport.bindAddress)
	//go func () { node1.createTask("printRing", createPrintMsg(node1.transport.bindAddress, node2.transport.bindAddress))}() 


/*
	time.Sleep(7000*time.Millisecond)
	node1.lookupReq("fingerLookup", "10", node5)
/*
	//msg := createPrintMsg(node2.transport.bindAddress, node3.transport.bindAddress)
	//go func () { node1.transport.send(msg)}() 


	node0.transport.listen()
*/
}