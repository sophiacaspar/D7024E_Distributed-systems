package dht

import (
	"fmt"
	"time"
	//"encoding/hex"
)

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      	string
	successor   	[2]string // 0: address, 1: nodeID
	predecessor 	[2]string // 0: address, 1: nodeID
	contact     	Contact
	fingers 		*FingerTable
	transport		*Transport
	taskQueue 		chan *Task
	responseQueue	chan *Msg
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
	dhtNode.responseQueue = make(chan *Msg)

	dhtNode.fingers = &FingerTable{}
	dhtNode.createTransport()
	//dhtNode.transport.listen()

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
	go dhtNode.stabilizeTimer()
	go dhtNode.fingerTimer()
	go dhtNode.transport.listen()

}

func (dhtNode *DHTNode) stabilizeTimer() {
	for {
		time.Sleep(time.Millisecond*2000)
		dhtNode.createTask("stabilize", nil)
		}	
}

// JOIN
func (dhtNode *DHTNode) addToRing(msg *Msg) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	//n := dhtNode.successor
	if 	(dhtNode.predecessor[0] == "" && dhtNode.successor[1] == dhtNode.nodeId) {
		dhtNode.successor[0] = msg.LightNode[0]
		dhtNode.successor[1] = msg.LightNode[1]
		newSucc := createUpdatePSMsg("updateSucc", msg.LightNode[0], [2]string{nodeAddress, dhtNode.nodeId})
		go func () { dhtNode.transport.send(newSucc)}() 

		dhtNode.setStaticFinger(&Msg{"", "", "", "","", dhtNode.successor, nil})
		f := createStatFingerMsg(nodeAddress, dhtNode.successor[0], [2]string{nodeAddress, dhtNode.nodeId})
		go func () { dhtNode.transport.send(f)}() 
	/*
		for i:=0; i < size; i++ {
			dhtNode.fingers.fingers[i] = &Finger{dhtNode.successor[0], dhtNode.successor[1]}
		}
*/
		//fmt.Println(dhtNode.nodeId, " fingers: ", dhtNode.fingers)

		//f1 := createFingerMsg(nodeAddress, msg.LightNode[0])
		//go func () { dhtNode.transport.send(f1)}() 


	} else if (between([]byte(dhtNode.nodeId), []byte(dhtNode.successor[1]), []byte(msg.LightNode[1]))){
		changeNewNodeSucc := createUpdatePSMsg("updateSucc", msg.LightNode[0], [2]string{dhtNode.successor[0], dhtNode.successor[1]})
		go func () { dhtNode.transport.send(changeNewNodeSucc)}() 

		dhtNode.successor[0] = msg.LightNode[0]
		dhtNode.successor[1] = msg.LightNode[1]

		dhtNode.setStaticFinger(&Msg{"", "", "", "","", dhtNode.successor, nil})
		f := createStatFingerMsg(nodeAddress, dhtNode.successor[0], [2]string{dhtNode.transport.bindAddress, dhtNode.nodeId})
		go func () { dhtNode.transport.send(f)}() 

/*
		for i:=0; i < size; i++ {
			dhtNode.fingers.fingers[i] = &Finger{dhtNode.successor[0], dhtNode.successor[1]}
		}
*/

		//f1 := createFingerMsg(nodeAddress, msg.LightNode[0])
		//go func () { dhtNode.transport.send(f1)}() 

		//dhtNode.successor.predecessor = newDHTNode
		//newDHTNode.successor = n

		//dhtNode.successor = newDHTNode
		//newDHTNode.predecessor = dhtNode
		
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
}


// WTF
func (dhtNode *DHTNode) getPredecessor(msg *Msg) {

		//m := createPredMsg(msg.Origin, msg.Src, dhtNode.predecessor)
		m := createResponseMsg(msg.Dst, msg.Src, dhtNode.predecessor)
		go func () { dhtNode.transport.send(m)} () 
		
		//myPred := &LightNode{dhtNode.predecessor[0], dhtNode.predecessor[1]}
		//m := createPredMsg(msg.Origin, msg.Src, myPred)
}

/******* OUTPUTS ********/
func (dhtNode *DHTNode) printRing(msg *Msg) {
	if msg.Origin != msg.Dst {
		fmt.Println("Pos in ring:", msg.Dst)
		msg := createPrintMsg(msg.Origin, dhtNode.successor[0])
		go func () { dhtNode.transport.send(msg)}() 
	} else {
		fmt.Println("Pos origin:", msg.Origin)
	}
	
}


func (dhtNode *DHTNode) init_taskQueue() {
	
	go func() {
		for {
			select {
				case t := <-dhtNode.taskQueue:
					switch t.taskType {
						case "addToRing":
							dhtNode.addToRing(t.msg)
						case "stabilize":
							dhtNode.stabilize()
						case "updateFingers":
							dhtNode.updateFingers()
						case "printRing":
							dhtNode.printRing(t.msg)
					}
				}	
			}		
		} ()
}


// periodically verify nodes immediate successor and tell the successor about node
func (dhtNode *DHTNode) stabilize(){
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	getSuccPred := createGetNodeMsg("pred", nodeAddress, dhtNode.successor[0])
	go func () { dhtNode.transport.send(getSuccPred)} () 
	//send msg
	waitResponse := time.NewTimer(time.Millisecond*5000)
	for {
		select {
			case r := <- dhtNode.responseQueue:
				if ((between([]byte(dhtNode.nodeId), []byte(dhtNode.successor[1]), []byte(r.LightNode[1]))) && r.LightNode[1] != "" ){
					dhtNode.successor[0] = r.LightNode[0]
					dhtNode.successor[1] = r.LightNode[1]
					return
				}
				
				notify := createNotifyMsg(nodeAddress, dhtNode.successor[0], [2]string{nodeAddress, dhtNode.nodeId})

				go func () { dhtNode.transport.send(notify)} () 
				fmt.Println(dhtNode.nodeId, dhtNode.successor, dhtNode.predecessor)
				return

			case t := <- waitResponse.C: //if timer is greater than 2000ms
				//check if alive
				fmt.Println(t, "successor timeout")
				return
		}
	}

}


func (dhtNode *DHTNode) notify(msg *Msg){
	if ((dhtNode.predecessor[0] == "") || between([]byte (dhtNode.predecessor[1]), []byte (dhtNode.nodeId), []byte (msg.LightNode[1]))){
		dhtNode.predecessor[0] = msg.LightNode[0]
		dhtNode.predecessor[1] = msg.LightNode[1]
	}
}


func (dhtNode *DHTNode) lookup(msg *Msg) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if (between([]byte(dhtNode.nodeId),[]byte(dhtNode.successor[1]), []byte(msg.Key))) {
		m := createResponseMsg(nodeAddress, msg.Origin, dhtNode.successor)
		go func () { dhtNode.transport.send(m)} () 
		//fmt.Println(dhtNode.successor, "is responsible for ", msg.Key)

    } else {
    	m := createLookupMsg(msg.Origin, msg.Key, nodeAddress, dhtNode.successor[0])
 		go func () { dhtNode.transport.send(m)} () 
    }
}

/*
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


