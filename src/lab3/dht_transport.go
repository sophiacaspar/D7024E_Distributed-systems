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
		if transport.dhtNode.online {
			msg := Msg{}
			err = dec.Decode(&msg)
			go func() { transport.msgQueue <- &msg } ()
		}

	}
}

func (transport *Transport) init_msgQueue() {
	go func() {
		for {
			select {
				case m := <-transport.msgQueue:
					//fmt.Println(transport.bindAddress, m.Type)
					switch m.Type {
						case "addToRing":
							transport.dhtNode.createTask("addToRing", m)
						case "updatePred":
							go transport.dhtNode.setPredecessor(m)
						case "updateSucc":
							go transport.dhtNode.setSuccessor(m)
						case "printRing":
							transport.dhtNode.createTask("printRing", m)
						case "printFinger":
							go transport.dhtNode.createTask("printRingFingers",m)
						case "pred":
							go transport.dhtNode.getPredecessor(m)
						case "response":
							transport.dhtNode.responseQueue <- m
						case "notify":
							go transport.dhtNode.notify(m)
							//transport.dhtNode.createTask("notify",m)
						case "lookup":
							//go transport.dhtNode.transport.send(createAckMsg(m.Dst, m.Origin))
							//fmt.Println(transport.bindAddress, "lookup", m.Key)
							go transport.dhtNode.lookup(m)
						case "lookupFound":
							transport.dhtNode.fingerMemory <- &Finger{m.LightNode[0], m.LightNode[1]}
						case "fingerLookup":
							go transport.dhtNode.fingerLookup(m)
						case "finger":
							transport.dhtNode.createTask("updateFingers", m)
						case "initFinger":
							go transport.dhtNode.initFingerTable(m)
						case "heartbeat":
							//transport.dhtNode.heartbeatQueue <- (createAckMsg(m.Dst, m.Origin))
							//go func () { transport.dhtNode.transport.send(createAckMsg(m.Dst, m.Origin))} ()
							transport.dhtNode.transport.send(createAckMsg(m.Dst, m.Origin))
						case "ack":
							transport.dhtNode.responseQueue <- m
					}
				}	
			}		
		} ()
}

func (transport *Transport) send(msg *Msg) {
	if transport.dhtNode.online {
		udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)

		conn, err := net.DialUDP("udp", nil, udpAddr)

		bytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
		}
		defer conn.Close()

		_, err = conn.Write(bytes)
	}
}

