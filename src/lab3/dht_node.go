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
	lookupQueue		chan *Msg
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

	dhtNode.online = true
	dhtNode.successor = [2]string{dhtNode.contact.ip + ":" + dhtNode.contact.port, dhtNode.nodeId}
	dhtNode.predecessor = [2]string{dhtNode.contact.ip + ":" + dhtNode.contact.port, dhtNode.nodeId}
	dhtNode.taskQueue = make(chan *Task)
	dhtNode.responseQueue = make(chan *Msg)
	dhtNode.heartbeatQueue = make(chan *Msg)
	dhtNode.fingerMemory = make(chan *Finger)
	dhtNode.lookupQueue = make(chan *Msg)

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
	go dhtNode.init_taskQueue()
	go dhtNode.stabilizeTimer()
	go dhtNode.fingerTimer()
	go dhtNode.heartbeatTimer()
	go dhtNode.transport.listen()

}


func (dhtNode *DHTNode) stabilizeTimer() {
	for {
			if dhtNode.online {
				time.Sleep(time.Millisecond*2000)
				go dhtNode.createTask("stabilize", nil)
			} else {
				return
			}
		}	
}

func (dhtNode *DHTNode) heartbeatTimer() {
	for {
			if dhtNode.online {
				time.Sleep(time.Millisecond*1000)
				go dhtNode.createTask("heartbeat", nil)	
			} else {
				return
			}
		}	
}


// JOIN
func (dhtNode *DHTNode) addToRing(msg *Msg) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if (between([]byte(dhtNode.nodeId), []byte(dhtNode.successor[1]), []byte(msg.LightNode[1]))){
			changeNewNodeSucc := createSetPreSuccMsg("setSucc", msg.LightNode[0], [2]string{dhtNode.successor[0], dhtNode.successor[1]})
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

func (dhtNode *DHTNode) getPredecessor(msg *Msg) {
		m := createResponseMsg(msg.Dst, msg.Origin, dhtNode.predecessor)
		go dhtNode.transport.send(m)

}


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
							fmt.Println(dhtNode.contact.port, "updated fingers")
						case "printRing":
							dhtNode.printRing(t.msg)
						case "printRingFingers":
							dhtNode.printRingFingers(t.msg)
						case "heartbeat":
							dhtNode.heartbeat()
						case "alive":
							//fmt.Println("node", dhtNode.nodeId, "is alive")
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

	// wait for response msg with my successor's predecessor
	waitResponse := time.NewTimer(time.Millisecond*2000)
	fmt.Println(nodeAddress, "Stab. Waiting for", getSuccPred.Dst)
	for {
		select {
			case r := <- dhtNode.responseQueue:
				if ((between([]byte(dhtNode.nodeId), []byte(dhtNode.successor[1]), []byte(r.LightNode[1]))) && r.LightNode[1] != "" && dhtNode.nodeId != r.LightNode[1]){
					dhtNode.successor[0] = r.LightNode[0]
					dhtNode.successor[1] = r.LightNode[1]
					fmt.Println(dhtNode.contact.port, "stabilized successor is", r.LightNode[0])
				}
				// I think I am your predecessor, update!
				notify := createNotifyMsg(nodeAddress, dhtNode.successor[0], [2]string{nodeAddress, dhtNode.nodeId})
				go dhtNode.transport.send(notify)
				
				//fmt.Println(dhtNode.nodeId, dhtNode.successor, dhtNode.predecessor)
				return
		// if we get no response, search for a finger that is alive to put as successor
			case <- waitResponse.C:
				//check if alive
				fmt.Println("xxxxxxxxxxxxxx  stabilize timeout", dhtNode.contact.port,  "xxxxxxxxxxxxxx")
				//dhtNode.updateSuccessor(dhtNode.successor[1])
				dhtNode.successor[0] = dhtNode.fingers.fingers[(size-1)].ip
				dhtNode.successor[1] = dhtNode.fingers.fingers[(size-1)].id
				return
		}
	} 
}

// Update predecessor if node should be between dhtnode and its predecessor
func (dhtNode *DHTNode) notify(msg *Msg){
	if ((dhtNode.predecessor[0] == "") || between([]byte (dhtNode.predecessor[1]), []byte (dhtNode.nodeId), []byte (msg.LightNode[1]))){
		dhtNode.predecessor[0] = msg.LightNode[0]
		dhtNode.predecessor[1] = msg.LightNode[1]
	}
	fmt.Println(dhtNode.predecessor[0], "is predecessor to", dhtNode.contact.port)
}

func (dhtNode *DHTNode) updateSuccessor(id string) {
	k := &Finger{dhtNode.successor[0], dhtNode.successor[1]}
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	// Get next alive finger without having to send 10000000 check-msg
	for _, key := range dhtNode.fingers.fingers {
		if key.id > id {
			getFingerPre := createGetNodeMsg("pred", nodeAddress, key.ip)
			go dhtNode.transport.send(getFingerPre)
			k = key
			break
		}
	}

	waitResponse := time.NewTimer(time.Millisecond * 500)
	for {
		select {
			case <-dhtNode.responseQueue:
				dhtNode.successor[0] = k.ip
				dhtNode.successor[1] = k.id
				fmt.Println(dhtNode.contact.port, "Successor updated to", dhtNode.successor[0])

				notify := createNotifyMsg(nodeAddress, k.ip, [2]string{nodeAddress, dhtNode.nodeId})
				go dhtNode.transport.send(notify)

				//dhtNode.successor = r.LightNode
				return
			case <-waitResponse.C:
				dhtNode.updateSuccessor(k.id)
				return
		}
	}
}

// Am I responsible for input-key?
func (dhtNode *DHTNode) responsible(key string) bool{
	if dhtNode.predecessor[1] == key {
		return false
	}
	if dhtNode.nodeId == key {
		return true
	} 
	return (between([]byte(dhtNode.predecessor[1]),[]byte(dhtNode.nodeId), []byte(key)))
}


func (dhtNode *DHTNode) lookup(key string) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	var m *Msg
	//fmt.Println("is", dhtNode.contact.port, "responsible for", msg.Key)
	
	if dhtNode.responsible(key) {
		dhtNode.fingerMemory <- &Finger{nodeAddress, dhtNode.nodeId}
	} else {
    	//lookupNextIp := dhtNode.getNextAlive(&Finger{dhtNode.successor[0], dhtNode.successor[1]})
    	m = createLookupMsg("lookup", nodeAddress, key, nodeAddress, dhtNode.successor[0])
		go dhtNode.transport.send(m)
    }
}


func (dhtNode *DHTNode) lookupNext(msg *Msg) {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	var m *Msg
	//fmt.Println("is", dhtNode.contact.port, "responsible for", msg.Key)
	waitResponse := time.NewTimer(time.Millisecond * 500)
	if dhtNode.responsible(msg.Key) {
		//fmt.Println(dhtNode.contact.port, "is responsible")
		m = createLookupFoundMsg(nodeAddress, msg.Origin, [2]string{dhtNode.contact.ip + ":" + dhtNode.contact.port, dhtNode.nodeId})
		//fmt.Println(m)
		go dhtNode.transport.send(m)
		waitResponse.Stop()
	} else {
		
    	//lookupNextIp := dhtNode.getNextAlive(&Finger{dhtNode.successor[0], dhtNode.successor[1]})
    	m = createLookupMsg("lookup", msg.Origin, msg.Key, nodeAddress, dhtNode.successor[0])
		go dhtNode.transport.send(m)
    	//fmt.Println(dhtNode.contact.port, "is not responsible, sending to", m.Dst)
    	waitResponse.Reset(time.Millisecond * 500)
    	for {
			select {
				case <-dhtNode.lookupQueue:
					//fmt.Println(nodeAddress, "got ack")
					//fmt.Println("msg from", dhtNode.contact.port, "reached", msg.Origin)
					return 
				case <-waitResponse.C:
					fmt.Println("==================", nodeAddress, "Lookup timeout ======================")
					/*
					for _, key := range dhtNode.fingers.fingers {
						if key.id > dhtNode.successor[1] {
							m = createLookupMsg("lookup", msg.Origin, msg.Key, nodeAddress, key.ip)
							go dhtNode.transport.send(m)
							break
						} 
					}
					*/
					return
				}
			}
    }
		
}

// If responsible node is'nt alive, then check next finger instead
func (dhtNode *DHTNode) getNextAlive(finger *Finger) string{
		k := finger
		nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
		// Get next alive finger without having to send 10000000 check-msg
		for _, key := range dhtNode.fingers.fingers {
		if key.id > finger.id || finger.id != dhtNode.successor[1]{
			isAlive := createGetNodeMsg("isAlive",nodeAddress, key.ip)
			go dhtNode.transport.send(isAlive)
			k = key
			break
		} 
	}

	waitResponse := time.NewTimer(time.Millisecond * 500)
	for {
		select {
			case <-dhtNode.responseQueue:
				if dhtNode.successor[0] != finger.id {
					dhtNode.successor[0] = k.ip
					dhtNode.successor[1] = k.id
				}
				return k.ip
			case <-waitResponse.C:
				fmt.Println(dhtNode.contact.port, "no respond from", k.ip,"-----------------------------------")
				return dhtNode.getNextAlive(k)
		}
	}
}


// uses fingers to lookup key, not used though
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


// Checks if my predecessor is alive, if not then reset my predecessors
func (dhtNode *DHTNode) heartbeat() {
	timeout := time.Millisecond*1000
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port

	if dhtNode.predecessor[0] != "" {
		m := createHeartbeatMsg(nodeAddress, dhtNode.predecessor[0])
		go func () { dhtNode.transport.send(m)} () 
		waitResponse := time.NewTimer(timeout)

		for {
			select {
				case <- dhtNode.heartbeatQueue:
				//case <- dhtNode.heartbeatQueue:
					//fmt.Println(dhtNode.contact.port, "predecessor is alive: ", dhtNode.predecessor[0])
					//fmt.Println(dhtNode.contact.port, "heartbeat respond from", dhtNode.predecessor[0])
					return
				case  <- waitResponse.C:
					fmt.Println("heartbeat timeout", dhtNode.contact.port)
					dhtNode.predecessor[0] = ""
					dhtNode.predecessor[1] = ""
					dhtNode.createTask("stabilize", nil)
						// fix data
					return
			}
		}
	}

}

// KILLS NODE
	func (dhtNode *DHTNode) kill() {
		fmt.Println("%!%!%!%!%!%!%!%!%!%!%!%!%!",dhtNode.contact.port, "is dead %!%!%!%!%!%!%!%!%!%!%!%!%!")
		dhtNode.online = false

	}



