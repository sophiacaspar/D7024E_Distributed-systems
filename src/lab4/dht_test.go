package dht

/** go test -test.run TestDHT1 */
/** Nodernas ordning: 4, 5, 2, 3, 7, 0, 6, 1   */

import (
	"fmt"
	"testing"
	"time"
	//"os"
	//"io/ioutil"
	//b64 "encoding/base64"
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
		//dhtNode.startServer()
		dhtNode.joinReq(master)


		dhtNode.newSuccessor()
		time.Sleep(300*time.Millisecond)
		nextmsg := createGetBackupMsg(dhtNode.transport.bindAddress, dhtNode.successor[0])
		go func () {dhtNode.transport.send(nextmsg)}() 

		fmt.Print("--------------------------------------------" + "\n")
		fmt.Print(dhtNode.transport.bindAddress + " gets data and file back if it exist in successor " + dhtNode.successor[0] + "\n")
		fmt.Print("--------------------------------------------" + "\n")

	} 
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


	//path := "/Users/Sophia/Documents/Coding/D7024E/D7024E_Distributed-systems/src/lab3/file/"
	/*
	path := "/Users/Sophia/workshop/go/src/lab4/file/"
	files, err := ioutil.ReadDir(path)

	if err != nil {
		panic(err)
	}

	for _, f := range files {
		file, _ := ioutil.ReadFile(path + f.Name())

		nodeId := generateNodeId(node4.transport.bindAddress)

		//sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
		//sFileData := b64.StdEncoding.EncodeToString(file)

		fmt.Print("--------------------------------------------" + "\n")
		fmt.Print("Send data: " + string(file) + " and file: " + f.Name() + " to: " + node4.transport.bindAddress + " with NodeId: " + nodeId + "\n")
		fmt.Print("--------------------------------------------" + "\n")
		node3.responsibleForFile([]byte(f.Name), []byte(file))
	}
*/
	fmt.Print("")
	time.Sleep(8000*time.Millisecond)
	fmt.Println("IS NODE 0 RESPONSIBLE FOR HASH?", node0.responsible("ba6bb37f774e7d2045da5a3abbb3d1687927667a"))
	fmt.Println("IS NODE 6 RESPONSIBLE FOR HASH?", node6.responsible("ba6bb37f774e7d2045da5a3abbb3d1687927667a"))
	//node3.printMyFingers()
	//fmt.Println("#####################", node3.responsible("bf06670af35ed4abcadd95abe8079568f4df38e6"), "#####################")
	node0.kill()

	time.Sleep(6000*time.Millisecond)
	//fmt.Println("IS NODE 0 RESPONSIBLE FOR HASH?", node0.responsible("ba6bb37f774e7d2045da5a3abbb3d1687927667a"))
	fmt.Println("IS NODE 6 RESPONSIBLE FOR HASH?", node6.responsible("ba6bb37f774e7d2045da5a3abbb3d1687927667a"))

	//node5.kill()

	//time.Sleep(3000*time.Millisecond)

	node7.alive(node1)

	//time.Sleep(6000*time.Millisecond)

	//node5.alive(node1)
	
	time.Sleep(2000*time.Second)

}