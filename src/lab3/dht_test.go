package dht

/* Make sure you are located in the same folder as the file before running command */
/* Command to run a test-case: "go test -test.run TestDHTX" where X i testnumber 1-2 */
/* Nodes arrangement: 4, 5, 2, 3, 7, 0, 6, 1 */

import (
	"fmt"
	"testing"
	"time"
	"io/ioutil"
	b64 "encoding/base64"
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
	} 
}

/* Testcase 1: 
Uploades a file in the folder "dataFolder/" located in the same place as in the folder the "dht_test.go" file is.
Kills the node holding the file, then killes the node that backsup the file. Then brings the first node back to life, waits 6 sec, then brings the second node back to life.
Stabilizes the file structure in "dataFolder/" as it was before any nodes died  */

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


	path := "/Users/Sophia/workshop/go/src/lab3/file/"
	time.Sleep(6000*time.Millisecond)


	/* path := "/Path/to/the/folder/with/file/ */

	//My path
	path := "file/" /* <- file.txt */

	/* Reads file in path and sends it to a function that finds the responsible node */

	files, err := ioutil.ReadDir(path)

	if err != nil {
		panic(err)
	}

	for _, f := range files {
		file, _ := ioutil.ReadFile(path + f.Name())

		sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
		sFileData := b64.StdEncoding.EncodeToString(file)

		node4.responsibleForFile(sFileName, sFileData)
	}

	fmt.Print("")
	time.Sleep(10000*time.Millisecond)

	node7.kill()

	time.Sleep(7000*time.Millisecond)

	node0.kill()

	time.Sleep(7000*time.Millisecond)

	node7.alive(node1)

	time.Sleep(6000*time.Millisecond)

	node0.alive(node1)

	time.Sleep(2000*time.Second)
}

/* Testcase 2: 
Uploades a file in the folder "dataFolder/" located in the same place as in the folder the "dht_test.go" file is. 
Kills the node holding the backup file. Waits 8 sec, then brings the node back to life 
Stabilizes the file structure in "dataFolder/" as it was before any nodes died  */

func TestDHT2(t *testing.T) {

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

	time.Sleep(6000*time.Millisecond)

	/* If you name a file "file.txt" it will be be assign to the node with portnumber 1117 when you upload it */

	/* Your path */
	/* path := "/Path/to/the/folder/with/file/ */

	//My path
	path := "file/" /* <- file.txt */

	/* Reads file in path and sends it to a function that finds the responsible node */
	files, err := ioutil.ReadDir(path)

	if err != nil {
		panic(err)
	}

	for _, f := range files {
		file, _ := ioutil.ReadFile(path + f.Name())

		sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
		sFileData := b64.StdEncoding.EncodeToString(file)

		node4.responsibleForFile(sFileName, sFileData)
	}

	fmt.Print("")
	time.Sleep(10000*time.Millisecond)

	node0.kill()

	time.Sleep(8000*time.Millisecond)

	node0.alive(node1)

	time.Sleep(2000*time.Second)
}