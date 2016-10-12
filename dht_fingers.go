package dht

import (
	"encoding/hex"
	"fmt"
	"time"
)

const bits int = 4

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

func (node *DHTNode) setFingers(msg *Msg) {
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

/*
func init_finger_table(n *DHTNode) [bits]*DHTNode{
	var templist [bits]*DHTNode
	for i := 0; i < bits; i++ {
		x,_ := hex.DecodeString(n.nodeId) //func DecodeString(s string) ([]byte, error)
		y, _ := calcFinger(x, (i+1), bits) // returnerar (string, []byte)
		if y == "" {
			y = "00"
		} else{}
		//fmt.Println("FINGER", i+1, "POINTS ON NODE", y)
		succ := n.lookup(y)  //ändra ej, vanlig lookup
		// calcFinger(n []byte, k int, m int) (string, []byte) {
		templist[i] = succ

	}
return templist
}
*/

func (node *DHTNode) updateFingers() {
	nodeAdress := node.contact.ip + ":" + node.contact.port
	for i := 0; i < bits; i++ {
		x, _ := hex.DecodeString(node.nodeId)
		y, _ := calcFinger(x, (i + 1), bits)
		booleanResponseTest := false
		if y == "" {
			y = "00"
		} else {
			responseTimmer := time.NewTimer(time.Second * 3)
			fingerMsg := fingerLookUpMessage(nodeAdress, y, nodeAdress, node.successor.adress)
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

/*if y == node.fingers.nodefingerlist[i].nodeId {
		} else {
			a := node.lookup(y)
			node.fingers.nodefingerlist[i] = a
		}
	}
	return node.fingers.nodefingerlist
}*/
