package dht

import (
	"fmt"
	//"encoding/hex"
)

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      	string
	successor   	[2]string
	predecessor 	[2]string
	contact     	Contact
	//finger_table 	*Finger_table
	transport		*Transport
	taskQueue 		chan *Task
}

type LightNode struct {
	address 	string
	nodeId 		string
}

type Task struct {
	taskType 	string
	msg 		*Msg
}



/*** CREATE ***/
func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port

	if nodeId == nil {
		genNodeId := generateNodeId()
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}

	dhtNode.successor = [2]string{dhtNode.contact.ip + ":" + dhtNode.contact.port, dhtNode.nodeId}
	dhtNode.predecessor = [2]string{}
	dhtNode.taskQueue = make(chan *Task)

	//dhtNode.finger_table = &Finger_table{}
	dhtNode.createTransport()
	//dhtNode.transport.listen()
	//dhtNode.finger_table.fingers = nil


	return dhtNode
}

func (dhtNode *DHTNode) createTransport() {
	dhtNode.transport = &Transport{dhtNode.contact.ip + ":" + dhtNode.contact.port, nil, dhtNode}
	dhtNode.transport.msgQueue = make(chan *Msg)
	dhtNode.transport.init_msgQueue()
	
}

func (dhtNode *DHTNode) createTask(taskType string, msg *Msg) {
	task := &Task{taskType, msg}
	dhtNode.taskQueue <- task
}

func (dhtNode *DHTNode) startServer() {
	fmt.Println("starting node ", dhtNode.nodeId)
	go dhtNode.init_taskQueue()
	go dhtNode.transport.listen()	
}

// JOIN
func (dhtNode *DHTNode) addToRing(msg *Msg) {
	//n := dhtNode.successor
	if 	(dhtNode.predecessor[0] == "" && dhtNode.successor[1] == dhtNode.nodeId) {
		fmt.Println("Only two nodes in ring")
		//dhtNode.setPredecessor(msg)
		dhtNode.successor[0] = msg.LightNode[0]
		dhtNode.successor[1] = msg.LightNode[1]
		//newPred := createUpdatePSMsg("updatePred", dhtNode.nodeId, dhtNode.transport.bindAddress, msg.Src)
		//newSucc := createUpdatePSMsg("updateSucc", dhtNode.nodeId, dhtNode.transport.bindAddress, msg.Src)
		newSucc := createUpdatePSMsg("updateSucc", msg.LightNode[0], [2]string{dhtNode.transport.bindAddress, dhtNode.nodeId})
		//go func () { dhtNode.transport.send(newPred)}() 
		go func () { dhtNode.transport.send(newSucc)}() 


	} else if (between([]byte(dhtNode.nodeId), []byte(dhtNode.successor[1]), []byte(msg.LightNode[1]))){
		fmt.Println(msg.LightNode[1], " is between ", dhtNode.nodeId, " and ", dhtNode.successor[1])
		fmt.Println("join with more nodes")
		//changeDHTNodeSuccPred := createUpdatePSMsg("updatePred", msg.Key, msg.Src, dhtNode.successor[0])
		changeNewNodeSucc := createUpdatePSMsg("updateSucc", msg.LightNode[0], [2]string{dhtNode.successor[0], dhtNode.successor[1]})
		//changeNewNodeSucc := createUpdatePSMsg("updateSucc", dhtNode.successor[1], dhtNode.successor[0], msg.Src)
		//changeNewNodePred := createUpdatePSMsg("updatePred", dhtNode.nodeId, dhtNode.transport.bindAddress, msg.Src)
		//go func () { dhtNode.transport.send(changeDHTNodeSuccPred)}() 
		go func () { dhtNode.transport.send(changeNewNodeSucc)}() 
		//go func () { dhtNode.transport.send(changeNewNodeSucc)}() 
		//go func () { dhtNode.transport.send(changeNewNodePred)}() 

		dhtNode.successor[0] = msg.LightNode[0]
		dhtNode.successor[1] = msg.LightNode[1]
		//dhtNode.successor.predecessor = newDHTNode
		//newDHTNode.successor = n

		//dhtNode.successor = newDHTNode
		//newDHTNode.predecessor = dhtNode
		
/**		TODO: FINGERS
		dhtNode.stabilize()
		newDHTNode.finger_table.fingers = init_finger_table(newDHTNode)
		dhtNode.update_fingers()
		
		//fmt.Print(dhtNode.nodeId)
		//fmt.Println(dhtNode.finger_table.fingers)

**/
		
	} else {
		forwardToSucc := createJoinMsg(dhtNode.successor[0], msg.LightNode)
		go func () { dhtNode.transport.send(forwardToSucc)} () 
		//n.addToRing(newDHTNode)
	}
}


func (dhtNode *DHTNode) setPredecessor(msg *Msg){
        dhtNode.predecessor[0] = msg.LightNode[0]
		dhtNode.predecessor[1] = msg.LightNode[1]
		//fmt.Println(dhtNode.nodeId, " predecessor: ", dhtNode.predecessor)
}

func (dhtNode *DHTNode) setSuccessor(msg *Msg) {
		dhtNode.successor[0] = msg.LightNode[0]
		dhtNode.successor[1] = msg.LightNode[1]
		fmt.Println(dhtNode.nodeId, " successor is ", dhtNode.successor)
}


// WTF
func (dhtNode *DHTNode) getPredecessor(msg *Msg) {
		m := createPredMsg(msg.Origin, msg.Src, dhtNode.predecessor)
		//myPred := &LightNode{dhtNode.predecessor[0], dhtNode.predecessor[1]}
		//m := createPredMsg(msg.Origin, msg.Src, myPred)

	
}


// responskö
// periodically verify nodes immediate successor and tell the successor about node
/*
func (dhtNode *DHTNode) stabilize(){
	n := 
	//n := dhtNode.successor.predecessor
	if (between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(n.nodeId))){
		dhtNode.successor = n
	}
	dhtNode.successor.notify(dhtNode)
}

/*
func (dhtNode *DHTNode) notify(node *DHTNode){
	if ((dhtNode.predecessor == nil) || between([]byte (dhtNode.predecessor.nodeId), []byte (dhtNode.nodeId), []byte (node.nodeId))){
		dhtNode.predecessor = node
	}
}


func (dhtNode *DHTNode) findSuccessor(msg *Msg) string, string {
		if (between([]byte(dhtNode.nodeId), byte[](msg.Key)) {
			
		}

        if (between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(key))){
		return dhtNode.successor
        } else {
                return dhtNode.successor.lookup(key)
        }
}


func (dhtNode *DHTNode) acceleratedLookupUsingFingers(key string) *DHTNode {

		fingerTable := dhtNode.finger_table.fingers

		for i := len(fingerTable); i > 0; i-- {
			fmt.Println(i)
			fmt.Println("Checks if ", key, " is between ", dhtNode.nodeId, " and ", fingerTable[(i-1)].nodeId)
			if !(between([]byte(dhtNode.nodeId), []byte(fingerTable[(i-1)].nodeId), []byte(key))){
				return fingerTable[(i-1)].acceleratedLookupUsingFingers(key)
			} 
		}
		return dhtNode.successor
		}



func (dhtNode *DHTNode) responsible(key string) bool {
	// TODO
	return false
}


*/

/******* OUTPUTS ********/
func (dhtNode *DHTNode) printRing(msg *Msg) {
	if msg.Origin != msg.Dst {
		fmt.Println(msg.Src)
		msg := createPrintMsg(msg.Origin, dhtNode.successor[0])
		go func () { dhtNode.transport.send(msg)}() 
	}	
}

/**
func (dhtNode *DHTNode) printRingFingers() {
	fmt.Print(dhtNode.nodeId, " [ ")
	dhtNode.printFingers()
	fmt.Println("]")
	for i := dhtNode.successor; i != dhtNode; i = i.successor {
		fmt.Print(i.nodeId, " [ ")
		i.printFingers()
		fmt.Println("]")
	}
}

func (dhtNode *DHTNode) printFingers() {
		for _, f := range dhtNode.finger_table.fingers {
			fmt.Print(f.nodeId, " ")
		}
}


func (dhtNode *DHTNode) testCalcFingers(m int, bits int) {
	idBytes, _ := hex.DecodeString(dhtNode.nodeId)
	fingerHex, _ := calcFinger(idBytes, m, bits)
	fingerSuccessor := dhtNode.lookup(fingerHex)
	fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
	fmt.Println("successor    " + fingerSuccessor.nodeId)

	dist := distance(idBytes, fingerSuccessorBytes, bits)
	fmt.Println("distance     " + dist.String())
}

*/

func (dhtNode *DHTNode) init_taskQueue() {
	
	go func() {
		for {
			select {
				case t := <-dhtNode.taskQueue:
					switch t.taskType {
						case "addToRing":
							fmt.Println(dhtNode.nodeId, " adds ", t.msg.LightNode[1])
							dhtNode.addToRing(t.msg)
							fmt.Println("exit", dhtNode.nodeId)
					}
				}	
			}		
		} ()
}
