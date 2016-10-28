package dht

import (
	"encoding/hex"
	"fmt"
	"time"
)

const bits int = 3

/*type FingerTable struct{cls
	Nodefingerlist [bits]*DHTNode
}
*/

type FingerTable struct {
	Nodefingerlist [bits]*Finger
}

type Finger struct {
	Id     string
	Adress string
}

func (node *DHTNode) setNetworkFingers(msg *Msg) {
	finger := &Finger{msg.LiteNode.Id, msg.LiteNode.Adress}
	for i := 0; i < bits; i++ {

		node.fingers.Nodefingerlist[i] = finger
	}
}

func (node *DHTNode) fingerTimer() {
	for {
		if node.alive {
			time.Sleep(time.Second * 5)
			node.createNewTask(nil, "updateFingers")
		} else {
			return
		}
	}
}

func (node *DHTNode) updateNetworkFingers() {
	//fmt.Println(node.contact.port, "updating fingers")
	nodeAdress := node.contact.ip + ":" + node.contact.port
	var booleanResponseTest = false
	for i := 0; i < bits; i++ {
		if node.fingers.Nodefingerlist[i] != nil {
			x, _ := hex.DecodeString(node.nodeId)
			y, _ := calcFinger(x, (i + 1), bits)
			if y == "" {
				y = "00"
			}

			//fmt.Println("update lookup")
			fingerMsg := lookUpMessage(nodeAdress, y, nodeAdress, node.successor.Adress)
			go node.transport.send(fingerMsg)
			responseTimmer := time.NewTimer(time.Second * 2)
			for booleanResponseTest != true {
				select {

				case responseCase := <-node.FingerQ:

					node.fingers.Nodefingerlist[i] = responseCase
					//fmt.Println("wtf", node.fingers.Nodefingerlist[i])
					//createdFinger := &Finger{responseCase.id, responseCase.adress} //id eller key?
					//node.fingers.Nodefingerlist[i] = createdFinger
					booleanResponseTest = true

				case <-responseTimmer.C:

		//			fmt.Println("timeout in updateNetworkFingers: ")
					booleanResponseTest = true
				}
			}
			booleanResponseTest = false
		}
	}
}

func (node *DHTNode) PrintOutNetworkFingers() {
	len_list := len(node.fingers.Nodefingerlist)
	for i := 0; i < len_list; i++ {
		fmt.Println(node.fingers.Nodefingerlist[i])
	}
}

func (node *DHTNode) printNetworkFingers(msg *Msg) {
	if msg.Origin != msg.Dst {
		fmt.Println("finger for node: ", node.nodeId, "is <")
		node.PrintOutNetworkFingers()
		fmt.Println(">")
		fingerPrintMsg := fingerPrintMessage(msg.Origin, node.successor.Adress)
		go func() { node.transport.send(fingerPrintMsg) }()
	} else {
		fmt.Println("finger for node ", node.nodeId, "is <")
		node.PrintOutNetworkFingers()
		fmt.Println(">")
	}
}

func (dhtnode *DHTNode) initPrintNetworkFingers(node *DHTNode) {
	printMsg := fingerPrintMessage(dhtnode.transport.BindAddress, node.transport.BindAddress)
	go func() {
		dhtnode.transport.send(printMsg)
	}()
}
