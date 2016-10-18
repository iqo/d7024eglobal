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

		node.fingers.nodefingerlist[i] = &Finger{msg.Id, msg.Adress}
		//id := node.nodeId
		//adress := node.contact.ip + ":" + node.contact.port

		//node.fingers.nodefingerlist[i] = &FingerTable{id,adress,"","","","","",""}
		//node.fingers.nodefingerlist[i] = &Finger{id, adress}
	}
}

func (node *DHTNode) fingerTimer() {
	for {
		if node.alive {
			time.Sleep(time.Second * 8)
			node.createNewTask(nil, "updateFingers")
		} else {
			return
		}
	}
}

func (node *DHTNode) updateNetworkFingers() {
	node.PrintOutNetworkFingers()
	//fmt.Println(node.contact.port, "updating fingers")
	nodeAdress := node.contact.ip + ":" + node.contact.port
	booleanResponseTest := false
	for i := 0; i < bits; i++ {
		if node.fingers.nodefingerlist[i] != nil {
			x, _ := hex.DecodeString(node.nodeId)
			y, _ := calcFinger(x, (i + 1), bits)
			if y == "" {
				y = "00"
			}

			//fmt.Println("update lookup")
			fingerMsg := lookUpMessage(nodeAdress, y, nodeAdress, node.successor.adress)
			go node.transport.send(fingerMsg)
			responseTimmer := time.NewTimer(time.Second * 3)
			for booleanResponseTest != true {
				select {

				case responseCase := <-node.fingerQ:
					fmt.Println("found", responseCase, "---")
					node.fingers.nodefingerlist[i] = responseCase
					//fmt.Println("wtf", node.fingers.nodefingerlist[i])
					//createdFinger := &Finger{responseCase.id, responseCase.adress} //id eller key?
					//node.fingers.nodefingerlist[i] = createdFinger
					booleanResponseTest = true

				case <-responseTimmer.C:

					fmt.Println("timeout in updateNetworkFingers: ")
					booleanResponseTest = true
				}
			}
			booleanResponseTest = false
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
		fmt.Println("finger for node ", node.nodeId, "is <")
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

/*func (node *DHTNode) initLookUpNetworkFingers(key string, dhtnode *DHTNode) {
	fingerLookUpMsg := fingerLookUpMessage(node.transport.bindAddress, key, node.transport.bindAddress, dhtnode.transport.bindAddress)
	fmt.Println("Finger lookup")
	go func() {
		dhtnode.transport.send(fingerLookUpMsg)
	}()
}

func (dhtnode *DHTNode) LookUpNetworkFinger(msg *Msg) {
	srcAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	tempFingerTable := dhtnode.fingers.nodefingerlist
	lenOfFingerTable := len(tempFingerTable)

	for i := lenOfFingerTable; i > 0; i-- {
		nodeBetween := (between([]byte(dhtnode.nodeId), []byte(tempFingerTable[i-1].id), []byte(msg.Key)))
		if nodeBetween != true {
			LookUpFingerMsg := fingerLookUpMessage(msg.Origin, msg.Key, srcAdress, tempFingerTable[i-1].adress)
			go func() {
				dhtnode.transport.send(LookUpFingerMsg)
			}()
			return
		}
	}
	if dhtnode.nodeId == msg.Key {
		fmt.Println("node ", dhtnode.nodeId, "is responsible for key", msg.Key)
	} else {
		fmt.Println("node ", dhtnode.successor.nodeId, "is responsible for key", msg.Key)
	}
	//fmt.Println(dhtnode.successor.nodeId)
	return
}*/
