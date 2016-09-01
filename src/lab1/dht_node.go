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

	dhtNode.successor = dhtNode
	dhtNode.predecessor = nil

	return dhtNode
}

/**** ADD SPECIAL CASE FOR ONLY ONE NODE ****/
func (dhtNode *DHTNode) join(node *DHTNode){
	predecessor = nil
	successor = node.findSuccessor(dhtNode)
}



func (dhtNode *DHTNode) findSuccessor(node *DHTNode){
	predNode := findPredecessor(node)
	return predNode.successor
}

func (dhtNode *DHTNode) findPredecessor(node *DHTNode){
	succsNode := dhtNode
	return succsNode
}



func (dhtNode *DHTNode) addToRing(newDHTNode *DHTNode) {
	// TODO
}

func (dhtNode *DHTNode) lookup(key string) *DHTNode {
	// TODO
	return dhtNode // XXX This is not correct obviously
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
