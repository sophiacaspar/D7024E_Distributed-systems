package main

import (
	dht "lab5/dht"
	"os"
)

func main() {
	port := os.Args[1]
	dht.StartUp(port)

/*
	//ip := "localhost:" + port
	//fmt.Println("Dockernode : " + ip + " is online")
	//bootstrapnode := "localhost:1110"

	node := dht.StartNode(nil, port)

	if ip != bootstrapnode {
		node.dht.JoinReq(bootstrapnode)
	}
	*/

}

/*
func (dhtNode *dht.DHTNode) joinReq(masterIP string) {
	msg := dht.createJoinMsg(masterIP, [2]string{dhtNode.transport.bindAddress, dhtNode.nodeId})
	go func () { dhtNode.transport.send(msg)}() 
	time.Sleep(time.Millisecond*500)
}

func startNode(nodeId *string, port string) (dhtNode *dht.DHTNode) {
	node := dht.makeDHTNode(nodeId, "localhost", port)
	dht.node.startServer()
	time.Sleep(300*time.Millisecond)
	return node
}
*/