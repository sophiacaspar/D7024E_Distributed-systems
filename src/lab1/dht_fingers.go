package dht

const size int = 3 
type Fingers struct {
	fingers 	[size]*DHTNode

}

func init_finger_table(node *DHTNode) [size]*DHTNode {
	var fingerTable [size]*DHTNode
	for i:=0; i < size; i++ {
	dBytes, _ := hex.DecodeString(dhtNode.nodeId)
	fingerHex, _ := calcFinger(idBytes, m, size)
	fingerSuccessor := dhtNode.lookup(fingerHex)

	fingerTable[i] = fingerSuccessor
	}
	return fingerTable
}