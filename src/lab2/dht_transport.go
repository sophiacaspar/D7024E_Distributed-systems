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
						case "hello": 
							fmt.Println(string(m.Key))
							transport.send(&Msg{"reply","", m.Src, []byte("fuck off")})
						case "reply":
							fmt.Println(string(m.Key))
						case "addToRing": 
							fmt.Println(string(m.Key))
							transport.send(&Msg{"ack","", m.Src, []byte("adding to ring")})
							transport.dhtNode.
					}
				}	
			}		
		} ()
}

func (transport *Transport) send(msg *Msg) {
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)

	conn, err := net.DialUDP("udp", nil, udpAddr)

	bytes, err := json.Marshal(msg)
	defer conn.Close()

	_, err = conn.Write(bytes)

	if err != nil {
		fmt.Println(err)
	}
}


