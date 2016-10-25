package dht

import "encoding/hex"

const size int = 160 
type Finger_table struct {
	fingers 	[size]*DHTNode

}

func init_finger_table(node *DHTNode) [size]*DHTNode {
	var fingerTable [size]*DHTNode
	for i:=0; i < size; i++ {
		idBytes, _ := hex.DecodeString(node.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i+1), size)

		if fingerHex == " " {
			fingerHex = "00"	
		} else{
			fingerSuccessor := node.lookup(fingerHex)
			fingerTable[i] = fingerSuccessor
		}
	}
	return fingerTable
}

func (dhtNode *DHTNode) update_fingers() {
	for i := dhtNode; i != dhtNode.successor; i = i.predecessor {
		i.finger_table.fingers = init_finger_table(i)
	}
}
