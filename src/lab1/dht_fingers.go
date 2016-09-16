package dht

import "encoding/hex"

const size int = 3 
type Finger_table struct {
	fingers 	[size]*DHTNode

}

/**
type Finger struct {
	id 		string
}
*/

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

/**
func lookup_fingers(node *DHTNode, key string) *DHTNode{
	fingerTable := node.finger_table
	var length = len(fingerTable)
	for x := 0; i := length; i > 1; i--; x++ {
		if (between([]byte(fingerTable[x].nodeId), []byte(fingerTable[i].nodeId), []byte(key))) {
			if fingerTable[x] == key {
				return fingerTable[x].nodeId
			} 
		} else {
			return lookup_fingers(fingerTable[i], key)
		}
	}
}
*/

/**
func (dhtNode *DHTNode) update_finger_table(fNode *DHTNode, i int){
	if (between([]byte(dhtNode.nodeId), []byte(dhtNode.finger_table.fingers[i].id), []byte(fNode.nodeId))){
		dhtNode.finger_table.fingers[i] = &Finger{fNode.nodeId}
		p := dhtNode.predecessor
		p.update_finger_table(fNode, i)
	}

}
*/