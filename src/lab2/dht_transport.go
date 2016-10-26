package dht

import (
	"net"
	"encoding/json"
	"fmt"
)

type Transport struct {
	bindAddress		string
	msgQueue		chan *Msg
	dhtNode			*DHTNode
}

func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress)
	conn, err := net.ListenUDP("udp", udpAddr)
	conn.SetReadBuffer(10000)
	conn.SetWriteBuffer(10000)
	if err != nil { 
		fmt.Println(err) 
	}

	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg := Msg{}
		err = dec.Decode(&msg)
		go func() {
			transport.msgQueue <- &msg
		} ()

	}
} 

func (transport *Transport) init_msgQueue() {
	go func() {
		for {
			select {
				case m := <-transport.msgQueue:
					switch m.Type {
						case "addToRing":
							transport.dhtNode.createTask("addToRing", m)
						case "updatePred":
							transport.dhtNode.setPredecessor(m)
						case "updateSucc":
							transport.dhtNode.setSuccessor(m)
						case "printRing":
							transport.dhtNode.createTask("printRing", m)
						case "printFinger":
							transport.dhtNode.printRingFingers(m)
						case "pred":
							transport.dhtNode.getPredecessor(m)
						case "response":
							transport.dhtNode.responseQueue <- m
						case "notify":
							transport.dhtNode.notify(m)
						case "lookup":
							go transport.dhtNode.lookup(m)
						case "fingerLookup":
							go transport.dhtNode.fingerLookup(m)
						case "finger":
							transport.dhtNode.createTask("updateFingers", m)
						case "statFinger":
							transport.dhtNode.setStaticFinger(m)
					}
				}	
			}		
		} ()
}

func (transport *Transport) send(msg *Msg) {
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)

	conn, err := net.DialUDP("udp", nil, udpAddr)

	bytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	_, err = conn.Write(bytes)
}