package dht

import (
	"encoding/json"
	"fmt"
	"net"
)

type Transport struct {
	Node        *DHTNode
	BindAddress string // rad 20, bindadress måste finnas.
	msgQ        chan *Msg
}

func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.BindAddress)
	//	fmt.Println("transport bindaddress:", transport.BindAddress)
	conn, err := net.ListenUDP("udp", udpAddr)
	conn.SetReadBuffer(10000)
	conn.SetWriteBuffer(10000)
	if err != nil {
		fmt.Println("error LISTEN function is:", err)
	}
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		if transport.Node.alive {
			msg := Msg{}
			err = dec.Decode(&msg)
			go func() {
				transport.msgQ <- &msg
			}()
		} else {
			return
		}
	}

}

func (transport *Transport) send(msg *Msg) {
	if transport.Node.alive {
		udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			fmt.Println("error SEND function is:", err)
		}
		encoded, err := json.Marshal(msg)
		defer conn.Close()
		_, err = conn.Write(encoded)
	}
}

func (transport *Transport) initmsgQ() {
	go func() {
		for {
			select {
			case msg := <-transport.msgQ:
				switch msg.Type {
				case "fingerPrint": //test case
					transport.Node.printNetworkFingers(msg)
					//transport.Node.TaskQ <- &Task{msg, "printRingFingers"}
				//	transport.send(&Msg{"printRing", "", v.Src, []byte("tjuuu")})
				case "reply": //test
					fmt.Println("hej:", string(msg.Data))
				case "printRing":
					go func() { transport.Node.TaskQ <- &Task{msg, "printRing"} }() //transport.Node.printRing()
					//transport.send(&Msg{"ring", "", v.Src, []byte(transport.Node.printRing())})
				case "addToRing":
					transport.Node.printNetworkRing(msg)
				case "response":
					go func() { transport.Node.ResponseQ <- msg }()
				case "join":
					go func() { transport.Node.TaskQ <- &Task{msg, "join"} }()
				case "notify":
					//		fmt.Println("notify network")
					go transport.Node.notifyNetwork(msg)
				case "pred":
					go transport.Node.getPred(msg)
				case "lookup":
					//fmt.Println("lookup transport ", msg)
					go transport.Node.transport.send(LookAckMessage(msg.Dst, msg.Src))
					go transport.Node.improvedNetworkLookUp(msg)
					//fmt.Println("initmsgQ lookup: ")
					//go transport.Node.networkLookup(msg)
				case "fingerLookup":
					go transport.Node.lookUpNetworkFinger(msg)
					//go transport.Node.lookupFingers(msg)
				case "heartBeat":
					if transport.Node.alive {
						transport.Node.transport.send(heartBeatAnswer(msg.Dst, msg.Origin))
						//transport.Node.transport.send(heartBeatAnswer(msg.Origin, msg.Dst))
					}
				case "heartAnswer":
					go func() { transport.Node.HeartBeatQ <- msg }()
				case "isAlive":
					if transport.Node.alive {
						transport.Node.transport.send(responseMessage(msg.Dst, msg.Origin, transport.BindAddress, transport.Node.nodeId))
					}
				case "nodeFound":
					go func() {
						transport.Node.FingerQ <- &Finger{msg.LiteNode.Id, msg.LiteNode.Adress} /*<- &Finger{msg.LiteNode.id, msg.LiteNode.adress}*/
					}()
				case "ack":
					go func() { transport.Node.ResponseQ <- msg }()
				case "fingerStart":
					go transport.Node.setNetworkFingers(msg)
				case "LookAck":
					go func() { transport.Node.NodeLookQ <- msg }()

				/* REPLICATE STUFF DOWN HERE*/

				case "Upload":
					go transport.Node.upload(msg)
				case "dataFromSuccessor":
					go transport.Node.dataFromSuccessor(msg)
				case "Replicate":
					go transport.Node.replicator(msg)
				case "deleteBackup":
					go transport.Node.deleteSuccessorBackup(msg)
				case "DeleteSuccessorFile":
					go transport.Node.deleteSuccessorFile(msg)

				}
			}
		}
	}()
}