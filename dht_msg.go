package dht

import (
//"fmt"
)

type Msg struct {
	Origin string
	Key    string //värdet
	Src    string //från noden som kalla
	Dst    string //destinationsadress
	//Data    string //transport funktionen, msg.Data
	Data     string
	LiteNode *LiteNodeStruct
	Type     string // type of message thats is being sent
	FileName string
}

type LiteNodeStruct struct {
	Adress string
	Id     string
}

func message(t, origin, dst, src, key string, Data string) *Msg {
	msg := &Msg{}
	msg.Type = t
	msg.LiteNode = &LiteNodeStruct{"", ""}
	//msg.Adress = ""
	//msg.Id = ""
	msg.Origin = origin
	msg.Src = src
	msg.Dst = dst
	msg.Data = Data
	msg.Key = key
	return msg
}

func joinMessage(dst string) *Msg {
	msg := &Msg{}
	msg.Type = "addToRing"
	msg.LiteNode = &LiteNodeStruct{"", ""}
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.Origin = "" //origin?
	msg.Src = ""
	msg.Dst = dst
	msg.Data = ""
	//msg.Key = key
	return msg
}

func printMessage(origin, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "printRing"
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.LiteNode = &LiteNodeStruct{"", ""}
	msg.Origin = origin
	msg.Src = ""
	msg.Dst = dst
	msg.Data = ""
	//msg.Key = key
	return msg
}

func notifyMessage(src, dst, adress, id string) *Msg {
	msg := &Msg{}
	msg.Type = "notify"
	//add adress to struct
	msg.LiteNode = &LiteNodeStruct{adress, id}
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.Origin = ""
	msg.Key = ""
	msg.Src = src
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func getPredMessage(origin, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "pred"
	msg.LiteNode = &LiteNodeStruct{"", ""}
	msg.Origin = origin
	msg.Src = ""
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func responseMessage(src, dst, adress, id string) *Msg {
	msg := &Msg{}
	msg.Type = "response"
	msg.LiteNode = &LiteNodeStruct{adress, id}
	/*msg.Adress = adress
	msg.Id = id*/
	msg.Origin = ""
	msg.Src = src
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func lookUpMessage(origin, key, src, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "lookup"
	msg.Key = key
	msg.LiteNode = &LiteNodeStruct{"", ""}
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.Origin = origin
	msg.Src = src
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func fingerLookUpMessage(origin, key, src, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "fingerLookup"
	msg.Key = key
	msg.LiteNode = &LiteNodeStruct{"", ""}
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.Origin = origin
	msg.Src = src
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func fingerPrintMessage(origin, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "fingerPrint"
	msg.Key = ""
	msg.LiteNode = &LiteNodeStruct{"", ""}
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.Origin = origin
	msg.Src = ""
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func heartBeatMessage(origin, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "heartBeat"
	//msg.Key = ""
	msg.LiteNode = &LiteNodeStruct{"", ""}
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.Origin = origin
	msg.Src = ""
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func heartBeatAnswer(origin, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "heartAnswer"
	msg.Key = ""
	msg.LiteNode = &LiteNodeStruct{"", ""}
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.Origin = origin
	msg.Src = ""
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func AliveMessage(origin, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "isAlive"
	msg.LiteNode = &LiteNodeStruct{"", ""}
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.Origin = origin
	msg.Src = ""
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func nodeFoundMessage(origin, dst, adress, key string) *Msg {
	msg := &Msg{}
	msg.Type = "nodeFound"
	msg.LiteNode = &LiteNodeStruct{adress, key}
	msg.Origin = origin
	msg.Src = ""
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func ackMsg(src, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "ack"
	msg.LiteNode = &LiteNodeStruct{"", ""}
	/*msg.Adress = ""
	msg.Id = ""*/
	msg.Origin = ""
	msg.Src = src
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func fingerStartMessage(src, dst, adress, id string) *Msg {
	msg := &Msg{}
	msg.Type = "fingerStart"
	msg.LiteNode = &LiteNodeStruct{adress, id}
	msg.Origin = ""
	msg.Src = src
	msg.Dst = dst
	msg.Data = ""
	return msg
}
func LookAckMessage(src, dst string) *Msg {
	msg := &Msg{}
	msg.Type = "LookAck"
	msg.LiteNode = &LiteNodeStruct{"", ""}
	msg.Origin = ""
	msg.Src = src
	msg.Dst = dst
	msg.Data = ""
	return msg
}

func UpLoadMessage(origin, dst, fname, data string) *Msg {
	msg := &Msg{}
	msg.Type = "Upload"
	msg.LiteNode = &LiteNodeStruct{"", ""}
	msg.Origin = ""
	msg.Src = ""
	msg.Dst = dst
	msg.Data = data
	msg.FileName = fname
	return msg
}

func ReplicateMessage(origin, dst, fname, data string) *Msg {
	msg := &Msg{}
	msg.Type = "Replicate"
	msg.LiteNode = &LiteNodeStruct{"", ""}
	msg.Origin = origin
	msg.Src = ""
	msg.Dst = dst
	msg.Data = data
	msg.FileName = fname
	return msg
}
