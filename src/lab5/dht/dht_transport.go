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
		} else{return}

	}
}

func (transport *Transport) init_msgQueue() {
	go func() {
		for {
			select {
				case m := <-transport.msgQueue:
				switch m.Type {
					case "addToRing":
						go transport.dhtNode.createTask("addToRing", m)
					case "setPred":
						go transport.dhtNode.setPredecessor(m)
					case "setSucc":
						go transport.dhtNode.setSuccessor(m)
					case "printRing":
						transport.dhtNode.createTask("printRing", m)
					case "printFinger":
						go transport.dhtNode.createTask("printRingFingers",m)
					case "pred":
						go transport.dhtNode.getPredecessor(m)	
					case "response":
						go func(){transport.dhtNode.responseQueue <- m}()
					case "notify":
						go transport.dhtNode.notify(m)
					case "lookup":
						go transport.dhtNode.transport.send(createAckMsg("lookupAck", m.Dst, m.Src))
						go transport.dhtNode.lookupNext(m)
					case "lookupFound":
						go func(){transport.dhtNode.fingerMemory <- &Finger{m.LightNode[0], m.LightNode[1]}}()
					case "fingerLookup":
						go transport.dhtNode.fingerLookup(m)
					case "finger":
						go transport.dhtNode.createTask("updateFingers", m)
					case "initFinger":
						go transport.dhtNode.initFingerTable(m)
					case "checkFinger":
						transport.dhtNode.transport.send(createResponseMsg(m.Dst, m.Origin, [2]string{transport.bindAddress, transport.dhtNode.nodeId}))
					case "heartbeat":
						if transport.dhtNode.online {transport.dhtNode.transport.send(createHeartbeatAnswer(m.Dst, m.Origin))} 
					case "heartbeatAnswer":
						go func(){transport.dhtNode.heartbeatQueue <- m}()
					case "isAlive":
						if transport.dhtNode.online{transport.dhtNode.transport.send(createResponseMsg(m.Dst, m.Origin, [2]string{transport.bindAddress, transport.dhtNode.nodeId}))}
					case "ack":
						go func() {transport.dhtNode.responseQueue <- m}()
					case "lookupAck":
						go func(){transport.dhtNode.lookupQueue <- m}()
					case "uploadData":
						go transport.dhtNode.addFile(m)
					case "replicate":
						go transport.dhtNode.replicate(m)
					case "checkSuccData":
						go transport.dhtNode.getSuccData(m)
					case "deleteFileSucc":
						go transport.dhtNode.deleteFileSucc(m)
					case "deleteFile":
						go transport.dhtNode.deleteFile(m)
					case "delBackup":
						go transport.dhtNode.deleteBackupSucc(m)
					case "updateFile":
						go transport.dhtNode.updateFile(m)
					case "getFiles":
						go transport.dhtNode.getFiles(m)
					case "fileResponse":
						go func(){transport.dhtNode.fileResponse <- &File{m.FileName, m.Data}}()
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