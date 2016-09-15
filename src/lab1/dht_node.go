package dht

import (
	"fmt"
	"encoding/hex"
)

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *DHTNode
	predecessor *DHTNode
	contact     Contact
	finger_table 	*Finger_table
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

	dhtNode.finger_table = &Finger_table{}
	//dhtNode.finger_table.fingers = nil

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
	n := dhtNode.successor
	if (dhtNode.predecessor == nil && dhtNode.successor == nil) {
		dhtNode.predecessor = newDHTNode
		dhtNode.successor = newDHTNode
		newDHTNode.predecessor = dhtNode
		newDHTNode.successor = dhtNode

		for i:=0; i < size; i++ {
			newDHTNode.finger_table.fingers[i] = &Finger{dhtNode.nodeId}
			dhtNode.finger_table.fingers[i] = &Finger{newDHTNode.nodeId}

		}

	} else if (between([]byte(dhtNode.nodeId), []byte(n.nodeId), []byte(newDHTNode.nodeId))){
		n.predecessor = newDHTNode
		newDHTNode.successor = n
		dhtNode.successor = newDHTNode
		newDHTNode.predecessor = dhtNode
		newDHTNode.finger_table.fingers = init_finger_table(newDHTNode)
		dhtNode.update_fingers()
		
		//fmt.Print(dhtNode.nodeId)
		//fmt.Println(dhtNode.finger_table.fingers)
		
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
		//node := lookup_fingers(dhtNode, key)
		return dhtNode
		}



func (dhtNode *DHTNode) responsible(key string) bool {
	// TODO
	return false
}

func (dhtNode *DHTNode) printRing() {
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
			fmt.Print(f.id, " ")
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
