package dht

const size int = 3 
type Finger_table struct {
	fingers 	[size]*Finger

}

type Finger struct {
	id 		string
}

func init_finger_table(node *DHTNode) [size]*DHTNode {
	var fingerTable [size]*Finger
	for i:=0; i < size; i++ {
		dBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i+1), size)
		fingerSuccessor := dhtNode.lookup(fingerHex)

		fingerTable[i] = fingerSuccessor
	}
	return fingerTable
}

func (dhtNode *DHTNode) update_fingers() {
	for i:=0; i < size; i++ {
		dBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i+1), size)
		
		if fingerHex != dhtNode.fingers.fingers[i].nodeId {
			fingerSuccessor := dhtNode.lookup(fingerHex)
			dhtNode.finger_table.fingers[i] = fingerSuccessor
		}
		}	
}

func (dhtNode *DHTNode) update_finger_table(fNode *DHTNode, i int){
	if (between([]byte(dhtNode.nodeId), []byte(dhtNode.finger_table.fingers[i].nodeId), []byte(fNode.nodeId))){
		dhtNode.finger_table.fingers[i] = new Finger{fNode.nodeId}
		p := dhtNode.predecessor
		p.update_finger_table(fNode, i)
	}

}