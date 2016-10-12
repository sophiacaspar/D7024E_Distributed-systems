package dht

type Msg struct {
	Type 		string
	Origin		string
	Key			string
	Src       	string
	Dst      	string
	LightNode	[2]string 	// 0: address, 1: nodeID
	Data		[]byte
}

func createMessage(t, origin, key, src, dst string, ln [2]string, bytes []byte) *Msg {
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = origin
	Msg.Key = key
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = ln
	Msg.Data = bytes
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
	Msg.Data = nil
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
	Msg.Data = nil
	return Msg
}

func createSetPreSuccMsg(t, dst string, newNodeInfo [2]string) *Msg{
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = newNodeInfo
	Msg.Data = nil
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
	Msg.Data = nil
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
	Msg.Data = nil
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
	Msg.Data = nil
	return Msg
}

func createGetNodeMsg(t, origin, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = origin
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.Data = nil
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
	Msg.Data = nil
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
	Msg.Data = nil
	return Msg
}

func createLookupFoundMsg(origin, dst string, node [2]string) *Msg{
	Msg := &Msg{}
	Msg.Type = "lookupFound"
	Msg.Origin = origin
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = node
	Msg.Data = nil
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
	Msg.Data = nil
	return Msg
}

func createInitFingerMsg(src, dst string, finger [2]string) *Msg{
	Msg := &Msg{}
	Msg.Type = "initFinger"
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = finger
	Msg.Data = nil
	return Msg
}

func createHeartbeatMsg(origin, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = "heartbeat"
	Msg.Origin = origin
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.Data = nil
	return Msg
}

func createHeartbeatAnswer(origin, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = "heartbeatAnswer"
	Msg.Origin = origin
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.Data = nil
	return Msg
}

func createAliveMsg(origin, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = "isAlive"
	Msg.Origin = origin
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.Data = nil
	return Msg
}



