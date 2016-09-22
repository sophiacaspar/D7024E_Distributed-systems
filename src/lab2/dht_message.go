package dht

type Msg struct {
	Type 		string
	Origin		string
	Key			string
	Src       	string
	Dst      	string
	Bytes		[]byte
}

func createMessage(t, origin, key, src, dst string, bytes []byte) *Msg {
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = origin
	Msg.Key = key
	Msg.Src = src
	Msg.Dst = dst
	Msg.Bytes = bytes
	return Msg
}

func createAckMsg(src, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = "ack"
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = src
	Msg.Dst = dst
	Msg.Bytes = nil
	return Msg
}

func createJoinMsg(key, src, dst string) *Msg {
	Msg := &Msg{}
	Msg.Type = "addToRing"
	Msg.Origin = ""
	Msg.Key = key
	Msg.Src = src
	Msg.Dst = dst
	Msg.Bytes = nil
	return Msg
	
}

func createUpdatePSMsg(t, nodeID, src, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = ""
	Msg.Key = nodeID
	Msg.Src = src
	Msg.Dst = dst
	Msg.Bytes = nil
	return Msg
}

func createPrintMsg(origin, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = "printRing"
	Msg.Origin = origin
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.Bytes = nil
	return Msg
}
