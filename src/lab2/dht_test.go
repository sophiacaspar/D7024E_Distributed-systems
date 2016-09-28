package dht

/** go test -test.run TestDHT1 */

import (
	//"fmt"
	"testing"
)

// dhtNode sends request to master of ring: please add me somewhere
func (dhtNode *DHTNode) joinReq(master *DHTNode) {
	msg := createJoinMsg(master.transport.bindAddress, [2]string{dhtNode.transport.bindAddress, dhtNode.nodeId})
	go func () { dhtNode.transport.send(msg)}() 
}

func TestDHT4(t *testing.T) {
	id0 := "00"
    id1 := "01"
    id2 := "02"
    id3 := "03"
    id4 := "04"
    id5 := "05"
/*
    id6 := "06"
    id7 := "07"
*/
    node0 := makeDHTNode(&id0, "localhost", "1110")
    node1 := makeDHTNode(&id1, "localhost", "1111")
    node2 := makeDHTNode(&id2, "localhost", "1112")
    node3 := makeDHTNode(&id3, "localhost", "1113")

    node4 := makeDHTNode(&id4, "localhost", "1114")
    node5 := makeDHTNode(&id5, "localhost", "1115")
 /*
    node6 := makeDHTNode(&id6, "localhost", "1116")
    node7 := makeDHTNode(&id7, "localhost", "1117")
*/

	node1.startServer()
	node2.startServer()
	node3.startServer()
	

	node4.startServer()
	node5.startServer()
		/*
	node6.startServer()
	node7.startServer()
	*/


	node5.joinReq(node1)
	node4.joinReq(node1)
	node3.joinReq(node1)
	node2.joinReq(node1)

	node0.transport.listen()

	}
