package dht

/** go test -test.run TestDHT1 */
/** Nodernas ordning: 4, 5, 2, 3, 7, 0, 6, 1   */

import (
	"fmt"
	"testing"
	"time"
	//"os"
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

<<<<<<< HEAD
=======
	//Sophias path
	//path := "/Users/Sophia/workshop/go/src/lab3/file/"

	path := "/Users/Zengin/Documents/Coding/D7024E/D7024E_Distributed-systems/src/lab3/file/"
>>>>>>> fd2b172aafcceeeced8b50a3880f136038bfe4d6
	
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
/*
	path2 := "/Users/Zengin/Documents/Coding/D7024E/D7024E_Distributed-systems/src/lab3/file2/"

	files2, error := ioutil.ReadDir(path2)

	if error != nil {
		panic(error)
	}

	for _, f := range files2 {
		file, _ := ioutil.ReadFile(path + f.Name())

		nodeId := generateNodeId(node5.transport.bindAddress)

		sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
		sFileData := b64.StdEncoding.EncodeToString(file)

		fmt.Print("--------------------------------------------" + "\n")
		fmt.Print("Send data: " + string(file) + " and file: " + f.Name() + " to: " + node4.transport.bindAddress + " with NodeId: " + nodeId + "\n")
		fmt.Print("--------------------------------------------" + "\n")
		
		msg := createUploadMsg(node4.transport.bindAddress, node5.transport.bindAddress, sFileName, sFileData)	
		go func () { node4.transport.send(msg)}()
	}
*/


	//node5.alive(node1)
	fmt.Print("")
	time.Sleep(10000*time.Millisecond)
	
	//node3.printMyFingers()
	//fmt.Println("#####################", node3.responsible("bf06670af35ed4abcadd95abe8079568f4df38e6"), "#####################")
	node7.kill()

	time.Sleep(7000*time.Millisecond)

	node0.kill()

	time.Sleep(7000*time.Millisecond)

	node7.alive(node1)

	time.Sleep(6000*time.Millisecond)

	node0.alive(node1)

	time.Sleep(2000*time.Second)

}