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
		if y == " " {
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
					createdFinger := &Finger{responseCase.Id, responseCase.Adress} //id eller key?
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

func (node *DHTNode) initLookUpNetworkFingers(key string, dhtnode *DHTNode) {
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
}
