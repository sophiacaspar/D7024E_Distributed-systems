package dht

import "encoding/hex"

const size int = 3 
type Finger_table struct {
	fingers 	[size]*Finger

}

type Finger struct {
	id 		string
}

func init_finger_table(node *DHTNode) [size]*Finger {
	var fingerTable [size]*Finger
	for i:=0; i < size; i++ {
		idBytes, _ := hex.DecodeString(node.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i+1), size)

		if fingerHex == " " {
			fingerHex = "00"	
		} else{
			fingerSuccessor := node.lookup(fingerHex)

			fingerTable[i] = &Finger{fingerSuccessor.nodeId}
		}
	}
	return fingerTable
}

func (dhtNode *DHTNode) update_fingers() {
	for i := dhtNode; i != dhtNode.successor; i = i.predecessor {
		dhtNode.finger_table.fingers = init_finger_table(i)
	}

/**
func (dhtNode *DHTNode) update_finger_table(fNode *DHTNode, i int){
	if (between([]byte(dhtNode.nodeId), []byte(dhtNode.finger_table.fingers[i].id), []byte(fNode.nodeId))){
		dhtNode.finger_table.fingers[i] = &Finger{fNode.nodeId}
		p := dhtNode.predecessor
		p.update_finger_table(fNode, i)
	}
*/
}