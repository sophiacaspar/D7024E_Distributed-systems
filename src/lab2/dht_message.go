package dht

type Msg struct {
	Type 		string
	Origin		string
	Key			string
	Src       	string
	Dst      	string
	Bytes		[]byte
}

func createMessage(t string, origin string, key string, src string, dst string, bytes []byte) *Msg {
	Msg := &Msg{}
	Msg.Type = t
	Msg.Origin = origin
	Msg.Key = key
	Msg.Src = src
	Msg.Dst = dst
	Msg.Bytes = bytes
	return Msg
}
func createACk(type, dst, src string)
