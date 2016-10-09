
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

func (dhtNode *DHTNode) updateFingers () {
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
/*
			if i == 0 {
				//dhtNode.transport.send(createMsg("lookup", fingerHex, src, dhtNode.fingers.fingerList[i].address, src))
				//go dhtNode.lookup(fingerHex)
				dhtNode.transport.send(createLookupMsg("lookup", nodeAddress, fingerHex, nodeAddress, dhtNode.successor[0]))
				//return
			} else {
				//dhtNode.transport.send(createLookupMsg("lookup", fingerHex, src, dhtNode.fingers.fingers[i-1].address, src))
				dhtNode.transport.send(createLookupMsg("lookup", nodeAddress, fingerHex, nodeAddress, dhtNode.fingers.fingers[(i-1)].ip))
			}
			*/
			waitRespons := time.NewTimer(time.Millisecond * 1000)
			for response != true {
				select {
				case s := <-dhtNode.fingerMemory:
					dhtNode.fingers.fingers[i] = s
					fmt.Println(dhtNode.contact.port, "added finger", (i+1), s.ip)
					response = true
				case <-waitRespons.C:
					fmt.Println("finger timeout,", dhtNode.contact.ip,"is searching for ",fingerHex)
					//fmt.Print("(CALLED FROM FINGERS) waiting respons from: ")
					//fmt.Println(dhtNode.fingers.fingerList[i-0].address)
					response = true
				//default:
				//	fmt.Println("when you try your best but don't succeed")
				}
			}
			response = false

		}
	}
}

func (dhtNode *DHTNode) fingerTimer() {
	for {
		time.Sleep(time.Millisecond*7300)
		fmt.Println("\n############ STABILIZING FINGERS FOR", dhtNode.nodeId,"###############")
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
			fmt.Print(f.ip, " ")
		}
}

