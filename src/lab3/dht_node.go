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
	heartbeatQueue	chan *Msg
	fingerMemory	chan *Finger
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
		genNodeId := generateNodeId(ip + ":" + port)
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}

	dhtNode.successor = [2]string{dhtNode.contact.ip + ":" + dhtNode.contact.port, dhtNode.nodeId}
	dhtNode.predecessor = [2]string{dhtNode.contact.ip + ":" + dhtNode.contact.port, dhtNode.nodeId}
	dhtNode.taskQueue = make(chan *Task)
	dhtNode.responseQueue = make(chan *Msg)
	dhtNode.heartbeatQueue = make(chan *Msg)
	dhtNode.fingerMemory = make(chan *Finger)

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
	//go dhtNode.fingerTimer()
	//go dhtNode.heartbeatTimer()
	go dhtNode.transport.listen()

}


func (dhtNode *DHTNode) stabilizeTimer() {
	for {
		time.Sleep(time.Millisecond*2000)
		go dhtNode.createTask("stabilize", nil)
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

		dhtNode.initFingerTable(&Msg{"", "", "", "","", dhtNode.successor, nil})
		f := createInitFingerMsg(nodeAddress, dhtNode.successor[0], [2]string{nodeAddress, dhtNode.nodeId})
		go func () { dhtNode.transport.send(f)}() 


	} else if (between([]byte(dhtNode.nodeId), []byte(dhtNode.successor[1]), []byte(msg.LightNode[1]))){
		changeNewNodeSucc := createUpdatePSMsg("updateSucc", msg.LightNode[0], [2]string{dhtNode.successor[0], dhtNode.successor[1]})
		go func () { dhtNode.transport.send(changeNewNodeSucc)}() 

		dhtNode.successor[0] = msg.LightNode[0]
		dhtNode.successor[1] = msg.LightNode[1]
		fmt.Println(dhtNode.contact.port, "successor is", msg.LightNode[0])


		dhtNode.initFingerTable(&Msg{"", "", "", "","", dhtNode.successor, nil})
		f := createInitFingerMsg(nodeAddress, dhtNode.successor[0], [2]string{dhtNode.transport.bindAddress, dhtNode.nodeId})
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
		fmt.Println(dhtNode.contact.port, "successor is", msg.LightNode[0])
}


// WTF
func (dhtNode *DHTNode) getPredecessor(msg *Msg) {
		m := createResponseMsg(msg.Dst, msg.Src, dhtNode.predecessor)
		go dhtNode.transport.send(m)

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
	time.Sleep(time.Second*5)
}


func (dhtNode *DHTNode) init_taskQueue() {
	
	go func() {
		for {
			select {
				case t := <-dhtNode.taskQueue:
					//fmt.Println(dhtNode.contact.port, " task ", t.taskType)
					switch t.taskType {
						case "addToRing":
							dhtNode.addToRing(t.msg)
						case "stabilize":
							dhtNode.stabilize()
						case "notify":
							dhtNode.notify(t.msg)
						case "updateFingers":
							dhtNode.updateFingers()
						case "printRing":
							dhtNode.printRing(t.msg)
						case "printRingFingers":
							dhtNode.printRingFingers(t.msg)
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
	// get successor's predecessor	
	getSuccPred := createGetNodeMsg("pred", nodeAddress, dhtNode.successor[0])
	go dhtNode.transport.send(getSuccPred)

	//send msg
	waitResponse := time.NewTimer(time.Millisecond*1000)
	for {
		select {
			//&& r.LightNode[1] != ""
			case r := <- dhtNode.responseQueue:
				if ((between([]byte(dhtNode.nodeId), []byte(dhtNode.successor[1]), []byte(r.LightNode[1]))) && r.LightNode[1] != "" && dhtNode.nodeId != r.LightNode[1]){
					dhtNode.successor[0] = r.LightNode[0]
					dhtNode.successor[1] = r.LightNode[1]
					fmt.Println(dhtNode.contact.port, "stabilized successor is", r.LightNode[0])
				}
				notify := createNotifyMsg(nodeAddress, dhtNode.successor[0], [2]string{nodeAddress, dhtNode.nodeId})
				go dhtNode.transport.send(notify)
				
				//fmt.Println(dhtNode.nodeId, dhtNode.successor, dhtNode.predecessor)
				return

			case <- waitResponse.C:
				//check if alive
				//changeNewNodeSucc := createUpdatePSMsg("updateSucc", nodeAddress, [2]string{dhtNode.fingers.fingers[1].ip, dhtNode.fingers.fingers[1].id})
				//go func () { dhtNode.transport.send(changeNewNodeSucc)}()
				fmt.Println("xxxxxxxxxxxxxx  stabilize timeout", dhtNode.contact.port,  "xxxxxxxxxxxxxx")
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
	fmt.Println(dhtNode.predecessor[0], "predecessor to", dhtNode.contact.port)
}

func (dhtNode *DHTNode) responsible(key string) bool{
	if dhtNode.nodeId == key {
		return true
	} else {
		return (between([]byte(dhtNode.nodeId),[]byte(dhtNode.successor[1]), []byte(key)))
	}
}


func (dhtNode *DHTNode) lookup(msg *Msg) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	var m *Msg
	if dhtNode.responsible(msg.Key) {
		// return successor
		//fmt.Println(dhtNode.successor[1], "is responsible for", msg.Key)
		m = createLookupFoundMsg(nodeAddress, msg.Origin, dhtNode.successor)
		dhtNode.transport.send(m)
    } else {
    	//fmt.Println(dhtNode.successor[1], "is not responsible for", msg.Key)
    	m = createLookupMsg("lookup", msg.Origin, msg.Key, nodeAddress, dhtNode.successor[0])
    	dhtNode.transport.send(m)
    }
	return
}

func (dhtNode *DHTNode) fingerLookup(msg *Msg) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	fingerTable := dhtNode.fingers.fingers

	for i := len(fingerTable); i > 0; i-- {
		//fmt.Println(i)
		//fmt.Println("Checks if ", msg.Key, " is between ", dhtNode.nodeId, " and ", fingerTable[(i-1)].id)
		if !(between([]byte(dhtNode.nodeId), []byte(fingerTable[(i-1)].id), []byte(msg.Key))){
			m := createLookupMsg(msg.Type, msg.Origin, msg.Key, nodeAddress, fingerTable[(i-1)].ip)
			go func () { dhtNode.transport.send(m)} () 
			return
			//return fingerTable[(i-1)].acceleratedLookupUsingFingers(key)
		} 
	}
	m := createLookupFoundMsg(nodeAddress, msg.Origin, dhtNode.successor)
	go func () { dhtNode.transport.send(m)} ()
	//fmt.Println(dhtNode.successor)
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
			case <- dhtNode.responseQueue:
			//case <- dhtNode.heartbeatQueue:
				//fmt.Println(dhtNode.nodeId, "predecessor is alive: ", dhtNode.predecessor[1])
				fmt.Println(dhtNode.contact.port, "heartbeat respond from", dhtNode.predecessor[0])
				return
			case  <- waitResponse.C:
				fmt.Println("heartbeat timeout")
				dhtNode.predecessor[0] = ""
				dhtNode.predecessor[1] = ""
					// fix data
				return
		}
	}

}


	func (dhtNode *DHTNode) kill() {
		fmt.Println("%!%!%!%!%!%!%!%!%!%!%!%!%!",dhtNode.contact.port, "is dead %!%!%!%!%!%!%!%!%!%!%!%!%!")
		dhtNode.online = false
	}



