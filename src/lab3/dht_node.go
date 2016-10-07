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
	online 			bool
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

	return dhtNode
}

func (dhtNode *DHTNode) createTransport() {
	dhtNode.transport = &Transport{dhtNode.contact.ip + ":" + dhtNode.contact.port, nil, dhtNode}
	dhtNode.transport.msgQueue = make(chan *Msg)
	dhtNode.transport.init_msgQueue()
}

func (dhtNode *DHTNode) createTask(taskType string, msg *Msg) {
	if dhtNode.online {
		task := &Task{taskType, msg}
		dhtNode.taskQueue <- task
	}
}

func (dhtNode *DHTNode) startServer() {
	fmt.Println("starting node ", dhtNode.nodeId)
	dhtNode.online = true
	go dhtNode.init_taskQueue()
	go dhtNode.stabilizeTimer()
	go dhtNode.fingerTimer()
	go dhtNode.heartbeatTimer()
	//go dhtNode.iAmAlive()
	go dhtNode.transport.listen()

}

func (dhtNode *DHTNode) stabilizeTimer() {
	for {
		time.Sleep(time.Millisecond*1000)
		dhtNode.createTask("stabilize", nil)
		}	
}

func (dhtNode *DHTNode) heartbeatTimer() {
	for {
		time.Sleep(time.Millisecond*2000)
		dhtNode.createTask("heartbeat", nil)
		}	
}


// JOIN
func (dhtNode *DHTNode) addToRing(msg *Msg) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port

	if 	(dhtNode.predecessor[0] == "" && dhtNode.successor[1] == dhtNode.nodeId) {
		dhtNode.successor[0] = msg.LightNode[0]
		dhtNode.successor[1] = msg.LightNode[1]
		newSucc := createUpdatePSMsg("updateSucc", msg.LightNode[0], [2]string{nodeAddress, dhtNode.nodeId})
		go func () { dhtNode.transport.send(newSucc)}() 

		dhtNode.setStaticFinger(&Msg{"", "", "", "","", dhtNode.successor, nil})
		f := createStatFingerMsg(nodeAddress, dhtNode.successor[0], [2]string{nodeAddress, dhtNode.nodeId})
		go func () { dhtNode.transport.send(f)}() 


	} else if (between([]byte(dhtNode.nodeId), []byte(dhtNode.successor[1]), []byte(msg.LightNode[1]))){
		changeNewNodeSucc := createUpdatePSMsg("updateSucc", msg.LightNode[0], [2]string{dhtNode.successor[0], dhtNode.successor[1]})
		go func () { dhtNode.transport.send(changeNewNodeSucc)}() 

		dhtNode.successor[0] = msg.LightNode[0]
		dhtNode.successor[1] = msg.LightNode[1]

		dhtNode.setStaticFinger(&Msg{"", "", "", "","", dhtNode.successor, nil})
		f := createStatFingerMsg(nodeAddress, dhtNode.successor[0], [2]string{dhtNode.transport.bindAddress, dhtNode.nodeId})
		go func () { dhtNode.transport.send(f)}() 

		
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
		fmt.Println(dhtNode.nodeId, "new successor", dhtNode.successor[1])
}


// WTF
func (dhtNode *DHTNode) getPredecessor(msg *Msg) {
		m := createResponseMsg(msg.Dst, msg.Src, dhtNode.predecessor)
		go func () { dhtNode.transport.send(m)} () 

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
					//fmt.Println(dhtNode.nodeId, " task ", t.taskType)
					switch t.taskType {
						case "addToRing":
							dhtNode.addToRing(t.msg)
						case "stabilize":
							dhtNode.stabilize()
						case "updateFingers":
							dhtNode.updateFingers()
						case "printRing":
							dhtNode.printRing(t.msg)
						case "heartbeat":
							dhtNode.heartbeat()
						case "alive":
							fmt.Println("node", dhtNode.nodeId, "is alive")
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
				//fmt.Println(dhtNode.nodeId, dhtNode.successor, dhtNode.predecessor)
				return

			case <- waitResponse.C: //if timer is greater than 2000ms
				//check if alive
				//changeNewNodeSucc := createUpdatePSMsg("updateSucc", nodeAddress, [2]string{dhtNode.fingers.fingers[1].ip, dhtNode.fingers.fingers[1].id})
				//go func () { dhtNode.transport.send(changeNewNodeSucc)}()
				dhtNode.successor[0] = dhtNode.fingers.fingers[1].ip
				dhtNode.successor[1] = dhtNode.fingers.fingers[1].id
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
	if (between([]byte(dhtNode.nodeId),[]byte(dhtNode.successor[1]), []byte(msg.Key)) && dhtNode.online) {
		m := createResponseMsg(nodeAddress, msg.Origin, dhtNode.successor)
		go func () { dhtNode.transport.send(m)} () 

    } else {
    	m := createLookupMsg(msg.Type, msg.Origin, msg.Key, nodeAddress, dhtNode.successor[0])
 		go func () { dhtNode.transport.send(m)} () 
    }
}


func (dhtNode *DHTNode) fingerLookup(msg *Msg) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	fingerTable := dhtNode.fingers.fingers

	for i := len(fingerTable); i > 0; i-- {
		//fmt.Println(i)
		fmt.Println("Checks if ", msg.Key, " is between ", dhtNode.nodeId, " and ", fingerTable[(i-1)].id)
		if !(between([]byte(dhtNode.nodeId), []byte(fingerTable[(i-1)].id), []byte(msg.Key))){
			m := createLookupMsg(msg.Type, msg.Origin, msg.Key, nodeAddress, fingerTable[(i-1)].ip)
			go func () { dhtNode.transport.send(m)} () 
			return
			//return fingerTable[(i-1)].acceleratedLookupUsingFingers(key)
		} 
	}
	fmt.Println(dhtNode.successor)
	return
	}

// if message arrives to this function, then node is online
/*
func (dhtNode DHTNode) checkOnline(msg *Msg) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	m := createResponseMsg(nodeAddress, msg.Origin)
	go func () { dhtNode.transport.send(m)} () 
}
*/


func (dhtNode *DHTNode) heartbeat() {
	timeout := time.Millisecond*300
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.predecessor[0] == "" {
		return
	}

	m := createHeartbeatMsg(nodeAddress, dhtNode.predecessor[0])
	go func () { dhtNode.transport.send(m)} () 
	waitResponse := time.NewTimer(timeout)

	for {
		select {
			case r := <- dhtNode.responseQueue:
				fmt.Println(r, dhtNode.nodeId, "predecessor is alive: ", dhtNode.predecessor[1])
				return
			case t := <- waitResponse.C: //if timer is greater than 2000ms
				dhtNode.predecessor[0] = ""
				dhtNode.predecessor[1] = ""
					// fix data
				return
				fmt.Println("node is dead", t)
		}
	}

}


	func (dhtNode *DHTNode) kill() {
		fmt.Println(dhtNode.nodeId, "is dead")
		dhtNode.online = false
	}



