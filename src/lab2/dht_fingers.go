
package dht

import (
	"encoding/hex"
	"time"
	"fmt"
	)

const size int = 3 
type FingerTable struct {
	fingers 	[size]*Finger
}

type Finger struct {
	Ip 		string 
	Id 		string
}

func (dhtNode *DHTNode) setStaticFinger(msg *Msg) {
		for i:=0; i < size; i++ {
			dhtNode.fingers.fingers[i] = &Finger{msg.LightNode[0], msg.LightNode[1]}
		}
}

func (dhtNode *DHTNode) updateFingers() {
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	for i:=0; i < size; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i+1), size)
		fmt.Println(fingerHex)

		if fingerHex == " " {
			fingerHex = "00"	
		} else{
			m := createLookupMsg(nodeAddress, fingerHex, nodeAddress, dhtNode.successor[0])
			go func () { dhtNode.transport.send(m)}() 

			waitResponse := time.NewTimer(time.Millisecond*2000)
			for {
				select {
					case r := <- dhtNode.responseQueue:

						newFinger := &Finger{r.LightNode[0], r.LightNode[1]}
						dhtNode.fingers.fingers[i] = newFinger
						fmt.Println(dhtNode.nodeId, " fingers ", dhtNode.fingers)

						return

					case t := <- waitResponse.C: //if timer is greater than 2000ms
						//check if alive
						fmt.Println(t, "finger timeout")
						return
				}
			}
		}

	}
	return 
}

func (dhtNode *DHTNode) fingerTimer() {
	for {
		time.Sleep(time.Millisecond*3000)
		dhtNode.createTask("updateFingers", nil)
	}	
}
