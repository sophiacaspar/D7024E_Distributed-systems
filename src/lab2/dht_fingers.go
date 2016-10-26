package dht

import (
	"encoding/hex"
	"time"
	"fmt"
	)

/* Runs in 160 bits 
const size int = 160 */

/* Runs in 3 bits */
const size int = 3 

type FingerTable struct {
	fingers 	[size]*Finger
}

type Finger struct {
	ip 		string 
	id 		string
}

func (dhtNode *DHTNode) setStaticFinger(msg *Msg) {
	for i:=0; i < size; i++ {
		dhtNode.fingers.fingers[i] = &Finger{msg.LightNode[0], msg.LightNode[1]}
	}
}

func (dhtNode *DHTNode) updateFingers() {
	
	nodeAddress := dhtNode.contact.ip + ":" + dhtNode.contact.port
	for i:=0; i < size; i++ {
		response := false
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i+1), size)

		if fingerHex == " " {
			fingerHex = "00"	
		} else{

			m := createLookupMsg("lookup", nodeAddress, fingerHex, nodeAddress, dhtNode.successor[0])
			go func () { dhtNode.transport.send(m)}() 

			waitResponse := time.NewTimer(time.Millisecond*2000)
			for response != true{
				select {
					case r := <- dhtNode.responseQueue:
						newFinger := &Finger{r.LightNode[0], r.LightNode[1]}
						dhtNode.fingers.fingers[i] = newFinger
						response = true

					case t := <- waitResponse.C: //if timer is greater than 2000ms
						//check if alive
						fmt.Println(t, "finger timeout")
						response = true
				}
			}
		}
	}
}

func (dhtNode *DHTNode) fingerTimer() {
	for {
		time.Sleep(time.Millisecond*3000)
		dhtNode.createTask("updateFingers", nil)
	}	
}

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
		fmt.Print(f.id, " ")
	}
}