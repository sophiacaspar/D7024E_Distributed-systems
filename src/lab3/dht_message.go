package dht

type Msg struct {
	Type 		string
	Origin		string
	Key			string
	Src       	string
	Dst      	string
	LightNode	[2]string 	// 0: address, 1: nodeID
	Data		string
	FileName 	string
}

func createAckMsg(t, src, dst string) *Msg{
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = ""
	Msg.Key = ""
	Msg.Src = src
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
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
	Msg.Data = ""
	return Msg
}

func createUploadMsg(origin, dst, filename string, data string) *Msg{
	Msg := &Msg{}
	Msg.Type = "uploadData"
	Msg.Origin = origin
	Msg.Key = ""
	Msg.Src = ""
	Msg.Dst = dst
	Msg.LightNode = [2]string{}
	Msg.FileName = filename
	Msg.Data = data
	return Msg
}

 func createReplicateMsg(origin, dst, filename string, data string) *Msg{
 	Msg := &Msg{}
 	Msg.Type = "replicate"
 	Msg.Origin = origin
 	Msg.Key = ""
 	Msg.Src = ""
 	Msg.Dst = dst
 	Msg.LightNode = [2]string{}
	Msg.FileName = filename
	Msg.Data = data
 	return Msg
 }

 func createCheckSuccDataMsg(origin, dst string, ln [2]string) *Msg{
 	Msg := &Msg{}
 	Msg.Type = "checkSuccData"
 	Msg.Origin = origin
 	Msg.Key = ""
 	Msg.Src = ""
 	Msg.Dst = dst
 	Msg.LightNode = ln
	Msg.FileName = ""
	Msg.Data = ""
 	return Msg
 }

 func createDeleteFileMsg(t, origin, dst, filename string) *Msg{
 	Msg := &Msg{}
 	Msg.Type = t
 	Msg.Origin = origin
 	Msg.Key = ""
 	Msg.Src = ""
 	Msg.Dst = dst
 	Msg.LightNode = [2]string{}
 	Msg.FileName = filename
 	Msg.Data = ""
 	return Msg
 }

 func createDeleteBackupeMsg(origin, dst, filename string) *Msg{
 	Msg := &Msg{}
 	Msg.Type = "delBackup"
 	Msg.Origin = origin
 	Msg.Key = ""
 	Msg.Src = ""
 	Msg.Dst = dst
 	Msg.LightNode = [2]string{}
 	Msg.FileName = filename
 	Msg.Data = ""
 	return Msg
 }