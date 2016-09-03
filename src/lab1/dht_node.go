package dht

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *DHTNode
	predecessor *DHTNode
	contact     Contact
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

	dhtNode.successor = nil
	dhtNode.predecessor = nil

	return dhtNode
}


func (dhtNode *DHTNode) findPredecessor(node *DHTNode) *DHTNode{
        succsNode := dhtNode
        return succsNode
}


func (dhtNode *DHTNode) findSuccessor(node *DHTNode) *DHTNode{
	predNode := dhtNode.findPredecessor(node)
	return predNode.successor
}


// JOIN
func (dhtNode *DHTNode) addToRing(newDHTNode *DHTNode) {
	//dhtNode.predecessor = nil
        //dhtNode.successor = newDHTNode.findSuccessor(dhtNode)
	n := dhtNode.successor
	if (dhtNode.predecessor == nil && dhtNode.successor == nil) {
		dhtNode.predecessor = newDHTNode
		dhtNode.successor = newDHTNode
		newDHTNode.predecessor = dhtNode
		newDHTNode.successor = dhtNode
	} else if (between([]byte(dhtNode.nodeId), []byte(n.nodeId), []byte(newDHTNode.nodeId))){
		n.predecessor = newDHTNode
		newDHTNode.successor = n
		dhtNode.successor = newDHTNode
		newDHTNode.predecessor = dhtNode
		
	} else {
		n.addToRing(newDHTNode)
	}


}


// periodically verify nodes immediate successor and tell the successor about node
func (dhtNode *DHTNode) stabilize(){
	n := dhtNode.successor.predecessor
	if (between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(n.nodeId))){
		dhtNode.successor = n
	}
	dhtNode.successor.notify(dhtNode)
}

func (dhtNode *DHTNode) notify(node *DHTNode){
	if ((dhtNode.predecessor == nil) || between([]byte (dhtNode.predecessor.nodeId), []byte (dhtNode.nodeId), []byte (node.nodeId))){
		dhtNode.predecessor = node
	}
}

func (dhtNode *DHTNode) lookup(key string) *DHTNode {
	if (between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(key))){
		return dhtNode.successor
	} else {
		return dhtNode.successor.lookup(key)
	}
}

func (dhtNode *DHTNode) acceleratedLookupUsingFingers(key string) *DHTNode {
	// TODO
	return dhtNode // XXX This is not correct obviously
}

func (dhtNode *DHTNode) responsible(key string) bool {
	// TODO
	return false
}

func (dhtNode *DHTNode) printRing() {
	// TODO
}

func (dhtNode *DHTNode) testCalcFingers(m int, bits int) {
	/* idBytes, _ := hex.DecodeString(dhtNode.nodeId)
	fingerHex, _ := calcFinger(idBytes, m, bits)
	fingerSuccessor := dhtNode.lookup(fingerHex)
	fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
	fmt.Println("successor    " + fingerSuccessor.nodeId)

	dist := distance(idBytes, fingerSuccessorBytes, bits)
	fmt.Println("distance     " + dist.String()) */
}
