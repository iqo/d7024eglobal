package dht

import (
	"encoding/hex"
	"fmt"
	"time"
)

const bits int = 3

/*type FingerTable struct{
	nodefingerlist [bits]*DHTNode
}
*/

type FingerTable struct {
	nodefingerlist [bits]*Finger
}

type Finger struct {
	id     string
	adress string
}

func (node *DHTNode) setNetworkFingers(msg *Msg) {
	for i := 0; i < bits; i++ {
		id := node.nodeId
		adress := node.contact.ip + ":" + node.contact.port

		//node.fingers.nodefingerlist[i] = &FingerTable{id,adress,"","","","","",""}
		node.fingers.nodefingerlist[i] = &Finger{id, adress}
	}
}

func (node *DHTNode) fingerTimer() {
	for {
		time.Sleep(time.Second * 3)
		node.createNewTask(nil, "updateFingers")
	}
}

func (node *DHTNode) updateNetworkFingers() {
	nodeAdress := node.contact.ip + ":" + node.contact.port
	for i := 0; i < bits; i++ {
		x, _ := hex.DecodeString(node.nodeId)
		y, _ := calcFinger(x, (i + 1), bits)
		booleanResponseTest := false
		if y == "" {
			y = "00"
		} else {
			responseTimmer := time.NewTimer(time.Second * 3)
			fingerMsg := lookUpMessage(nodeAdress, y, nodeAdress, node.successor.adress)
			go func() {
				node.transport.send(fingerMsg)
			}()
			for booleanResponseTest != true {
				select {

				case responseCase := <-node.responseQ:
					createdFinger := &Finger{responseCase.Adress, responseCase.Id} //id eller key?
					node.fingers.nodefingerlist[i] = createdFinger
					booleanResponseTest = true

				case e := <-responseTimmer.C:

					fmt.Println(e, "timeout: ")
					booleanResponseTest = true
				}
			}
		}
	}
}

func (node *DHTNode) PrintOutNetworkFingers() {
	len_list := len(node.fingers.nodefingerlist)
	for i := 0; i < len_list; i++ {
		fmt.Println(node.fingers.nodefingerlist[i])
	}
}

func (node *DHTNode) printNetworkFingers(msg *Msg) {
	if msg.Origin != msg.Dst {
		fmt.Println("finger for node: ", node.nodeId, "is <")
		node.PrintOutNetworkFingers()
		fmt.Println(">")
		fingerPrintMsg := fingerPrintMessage(msg.Origin, node.successor.adress)
		go func() { node.transport.send(fingerPrintMsg) }()
	} else {
		fmt.Println("finger for node: ", node.nodeId, "is <")
		node.PrintOutNetworkFingers()
		fmt.Println(">")
	}
}

func (dhtnode *DHTNode) initPrintNetworkFingers(node *DHTNode) {
	printMsg := fingerPrintMessage(dhtnode.transport.bindAddress, node.transport.bindAddress)
	go func() {
		dhtnode.transport.send(printMsg)
	}()
}
