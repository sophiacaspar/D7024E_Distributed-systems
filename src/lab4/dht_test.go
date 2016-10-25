package dht

/* Make sure you are located in the same folder as the file before running command */
/* Command to run a test-case: "go test -test.run TestDHT1" */
/* Nodes arrangement: 4, 5, 2, 3, 7, 0, 6, 1 */

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
		go dhtNode.init_taskQueue()
		go dhtNode.stabilizeTimer()
		go dhtNode.fingerTimer()
		go dhtNode.heartbeatTimer()
		go dhtNode.transport.listen()
		dhtNode.joinReq(master)
	} 
}

/* Testcase 1: Here you can design your own test. When starting the test, open a browser and type in localhost:port (port: 1110 <-> 1117) 
You can in the webpage upload, delete and change data in files. The uploaded file will be available in all adresses in portrange 1110 <-> 1117

If you name a file "file.txt" it will be be assign to the node with portnumber 1117 when you upload it. 
Now you can create a scenario were a node is killed and then alive after a specific time. You could for example kill the node holding the file or the node holding the backup. 

If you kill the node 1117, you can check this by writing localhost:1117 in the browser. This will not work obviously. The page will load, but will not present any avalible files 
But if you write in some other adress, for example localhost:1110, the avalible files will load and show. */

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

	fmt.Print("")
	time.Sleep(10000*time.Millisecond)
	
	/* Killes the node with the file (if filenamen is "file.txt") */
	node7.kill()

	time.Sleep(10000*time.Millisecond)

	node7.alive(node1)

	/* Killes the node with the backup folder with the file "file.txt" in it
	node0.kill()

	time.Sleep(10000*time.Millisecond)

	node0.alive(node1) */

	time.Sleep(2000*time.Second)
}