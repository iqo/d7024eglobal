package dht

import (
	"encoding/json"
	"fmt"
	"net"
)

type Transport struct {
	node        *DHTNode
	bindAddress string // rad 20, bindadress måste finnas.
	msgQ        chan *Msg
}

func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress)
	//	fmt.Println("transport bindaddress:", transport.bindAddress)
	conn, err := net.ListenUDP("udp", udpAddr)
	conn.SetReadBuffer(10000)
	conn.SetWriteBuffer(10000)
	if err != nil {
		fmt.Println("error LISTEN function is:", err)
	}
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg := Msg{}
		err = dec.Decode(&msg)
		go func() {
			transport.msgQ <- &msg
		}()

	}

}

func (transport *Transport) send(msg *Msg) {
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("error SEND function is:", err)
	}
	encoded, err := json.Marshal(msg)
	defer conn.Close()
	_, err = conn.Write(encoded)

}

func (transport *Transport) initmsgQ() {
	go func() {
		for {
			select {
			case msg := <-transport.msgQ:
				switch msg.Type {
				case "fingerPrint": //test case
					transport.node.printNetworkFingers(msg)
					//transport.node.TaskQ <- &Task{msg, "printRingFingers"}
				//	transport.send(&Msg{"printRing", "", v.Src, []byte("tjuuu")})
				case "reply": //test
					fmt.Println("hej:", string(msg.Bytes))
				case "printRing":
					transport.node.TaskQ <- &Task{msg, "printRing"} //transport.node.printRing()
					//transport.send(&Msg{"ring", "", v.Src, []byte(transport.node.printRing())})
				case "addToRing":
					transport.node.printNetworkRing(msg)
				case "response":
					transport.node.responseQ <- msg
				case "join":
					transport.node.TaskQ <- &Task{msg, "join"}
				case "notify":
					//		fmt.Println("notify network")
					transport.node.notifyNetwork(msg)
				case "pred":
					transport.node.getPred(msg)
				case "lookup":
					//fmt.Println("initmsgQ lookup: ")
					go transport.node.networkLookup(msg)
				case "fingerLookup":
					go transport.node.LookUpNetworkFinger(msg)
					//go transport.node.lookupFingers(msg)
				}
			}
		}
	}()
}
