package dht

import (
	"encoding/hex"
	"time"
	"fmt"
	)

const size int = 160

type FingerTable struct {
	fingers 	[size]*Finger
}

type Finger struct {
	ip 		string 
	id 		string
}

func (dhtNode *DHTNode) initFingerTable(msg *Msg) {
	for i:=0; i < size; i++ {
		dhtNode.fingers.fingers[i] = &Finger{msg.LightNode[0], msg.LightNode[1]}
	}
}

func (dhtNode *DHTNode) updateFingers() {
	if dhtNode.successor[0] != "" {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	var response = false
		for i := 0; i < size; i++ {
			if dhtNode.fingers.fingers[i] != nil {
				idBytes, _ := hex.DecodeString(dhtNode.nodeId)
				fingerHex, _ := calcFinger(idBytes, (i + 1), size)
				
				if fingerHex == " " {
					fingerHex = "00"
				}
				go dhtNode.transport.send(createLookupMsg("lookup", nodeAddress, fingerHex, nodeAddress, dhtNode.successor[0]))

				waitResponse := time.NewTimer(time.Millisecond * 2000)
				for response != true {
					select {
					case s := <-dhtNode.fingerMemory:
						dhtNode.fingers.fingers[i] = s
						response = true
					case <-waitResponse.C:
						fmt.Println("finger timeout,", dhtNode.contact.port,"is searching for ",fingerHex)
						response = true

					}
				}
				response = false
			}
		}
	}
}

func (dhtNode *DHTNode) fingerTimer() {
	for {
		if dhtNode.online {
			time.Sleep(time.Millisecond*6500)
			dhtNode.createTask("updateFingers", nil)
		} else {
			return
		}
	}	
}

/*********************************************
**** PRINTS FINGERS WITH DIFFERENT OUTPUTS ***
*********************************************/
func (dhtNode *DHTNode) printRingFingers(msg *Msg) {
	if msg.Origin != msg.Dst {
		fmt.Print(dhtNode.nodeId, " [ ")
		dhtNode.printFingers()
		fmt.Println("]")
		msg := createPrintFingerMsg(msg.Origin, dhtNode.successor[0])
		go func () { dhtNode.transport.send(msg)}() 
	} else {
		fmt.Print(dhtNode.nodeId, " [ ")
		dhtNode.printFingers()
		fmt.Println("]")
	}
}

func (dhtNode *DHTNode) printFingers() {
		for _, f := range dhtNode.fingers.fingers {
			fmt.Print(f.ip, " ")
		}
}

func (dhtNode *DHTNode) printMyFingers() {
	fmt.Println("")
	fmt.Print(dhtNode.contact.ip, ": [")
		for _, f := range dhtNode.fingers.fingers {
			fmt.Print(f.ip, " ")
		}
	fmt.Println("]")
}