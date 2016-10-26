package dht

type Msg struct {
	Type 		string
	Origin		string
	Key			string
	Src       	string
	Dst      	string
	LightNode	[2]string 	// 0: address, 1: nodeID
	Bytes		[]byte
}

func createMessage(t, origin, key, src, dst string, ln [2]string, bytes []byte) *Msg {
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = origin
	Msg.Key = key
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = ln
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
	Msg.LightNode = [2]string{}
	Msg.Bytes = nil
	return Msg
}

func createJoinMsg(dst string, newNodeInfo [2]string) *Msg {
	Msg := &Msg{}
	Msg.Type = "addToRing"
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = newNodeInfo
	Msg.Bytes = nil
	return Msg
}

func createUpdatePSMsg(t, dst string, newNodeInfo [2]string) *Msg{
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = newNodeInfo
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
	Msg.LightNode = [2]string{}
	Msg.Bytes = nil
	return Msg
}

func createPrintFingerMsg(origin, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = "printFinger"
	Msg.Origin = origin
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.Bytes = nil
	return Msg
}

func createResponseMsg(src, dst string, ln [2]string) *Msg{
	Msg := &Msg{}
	Msg.Type = "response"
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = ln
	Msg.Bytes = nil
	return Msg
}

func createGetNodeMsg(t, src, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.Bytes = nil
	return Msg
}

func createNotifyMsg(src, dst string, ln [2]string) *Msg{
	Msg := &Msg{}
	Msg.Type = "notify"
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = ln
	Msg.Bytes = nil
	return Msg
}

func createLookupMsg(t, origin, key, src, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = origin
	Msg.Key = key
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.Bytes = nil
	return Msg
}

func createFingerMsg(src, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = "finger"
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.Bytes = nil
	return Msg
}

func createStatFingerMsg(src, dst string, finger [2]string) *Msg{
	Msg := &Msg{}
	Msg.Type = "statFinger"
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = finger
	Msg.Bytes = nil
	return Msg
}