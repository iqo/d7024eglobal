package dht

import (
	//"encoding/hex"
	"fmt"
	"time"
)

func (dhtnode *DHTNode) resposibleNetworkNode(key string) bool {
	if dhtnode.predecessor.NodeId == key {
		//fmt.Println("this is not know ")
		return false
	}
	if dhtnode.nodeId == key {
		//fmt.Println("this is know ")
		return true
	}

	//beeweetNodes := (between([]byte(dhtnode.predecessor.nodeId), []byte(dhtnode.nodeId), []byte(key)))
	//return beeweetNodes
	return (between([]byte(dhtnode.predecessor.NodeId), []byte(dhtnode.nodeId), []byte(key)))
}

/*func (dhtnode *DHTNode) findNextAlive(key int) string {
	dhtAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	//fmt.Println("dht adress:", dhtAdress, "node fingerlist adress", dhtnode.fingers.Nodefingerlist[key].adress)
	notDead := AliveMessage(dhtAdress, dhtnode.fingers.Nodefingerlist[key].adress)
	go dhtnode.transport.send(notDead)
	timerResp := time.NewTimer(time.Millisecond * 100)
	for {
		select {
		case r := <-dhtnode.ResponseQ:
			if r.LiteNode.adress != "" {
				//fmt.Println("lookUp ", r.Adress)
				return r.LiteNode.adress
			} else {
				return dhtnode.findNextAlive(key + 1)
			}
		case <-timerResp.C:
			fmt.Println(dhtnode.contact.port, "no response from", dhtnode.fingers.Nodefingerlist[key].adress)
			if key < (bits - 1) {
				return dhtnode.findNextAlive(key + 1)
			}
		}
	}
}*/

/*func (dhtnode *DHTNode) improvedNetworkLookUp(msg *Msg) {
	dhtAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	timerResp := time.NewTimer(time.Millisecond * 100)
	if dhtnode.resposibleNetworkNode(msg.Key) {
		fmt.Println("")
		fmt.Println("dhtnodeid", dhtnode.nodeId)
		fmt.Println("")
		foundMsg := nodeFoundMessage(dhtAdress, msg.Origin, dhtAdress, dhtnode.nodeId)
		go dhtnode.transport.send(foundMsg)

		for {
			select {
			case <-dhtnode.ResponseQ:
				return
			case <-timerResp.C:
				fmt.Println("-------------------------------------------------")
				fmt.Println("fuck timer ")
				fmt.Println("-------------------------------------------------")

				dhtnode.transport.send(foundMsg)
			}
		}
		return
	} else {
		fmt.Println("-------------------------------------------------")
		fmt.Println("else fuck")
		fmt.Println("-------------------------------------------------")
		//fmt.Println("fin next alive")
		next := dhtnode.findNextAlive(0)
		fmt.Println("-------------------------------------------------")
		fmt.Println("next")
		fmt.Println("-------------------------------------------------")
		lookUpMsg := lookUpMessage(msg.Origin, msg.Key, dhtAdress, next)
		go dhtnode.transport.send(lookUpMsg)
		//fmt.Println(dhtnode.nodeId)
	}
	return
}*/

func (dhtnode *DHTNode) initNetworkLookUp(key string) {
	nodeadress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	if dhtnode.resposibleNetworkNode(key) {
		dhtnode.FingerQ <- &Finger{dhtnode.nodeId, nodeadress}
	} else {
		lookUpMsg := lookUpMessage(nodeadress, key, nodeadress, dhtnode.successor.Adress)
		go dhtnode.transport.send(lookUpMsg)
	}
}

func (dhtnode *DHTNode) improvedNetworkLookUp(msg *Msg) {
	NodeAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	timeResp := time.NewTimer(time.Second * 1)

	if dhtnode.resposibleNetworkNode(msg.Key) {
		nodeFoundMsg := nodeFoundMessage(NodeAdress, msg.Origin, NodeAdress, dhtnode.nodeId)
		go dhtnode.transport.send(nodeFoundMsg)
		timeResp.Stop()
	} else {
		lookUpMsg := lookUpMessage(msg.Origin, msg.Key, NodeAdress, dhtnode.successor.Adress)
		go dhtnode.transport.send(lookUpMsg)
		timeResp.Reset(time.Second * 1)

		for {
			select {
			case <-dhtnode.NodeLookQ:
				return

			case <-timeResp.C:
				//fmt.Println("fuck this lookup")
				return
			}
		}
	}
}

func (dhtnode *DHTNode) findNextAlive(fing *Finger) string {
	tempFinger := fing
	timeResp := time.NewTimer(time.Millisecond * 500)
	dhtAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	lenOfFingerList := len(dhtnode.fingers.Nodefingerlist)
	tempFingerList := dhtnode.fingers.Nodefingerlist
	for i := 0; i < lenOfFingerList; i++ {
		if tempFingerList[i].Id > fing.Id || fing.Id != dhtnode.successor.NodeId {
			aliveMsg := AliveMessage(dhtAdress, tempFingerList[i].Adress)
			go dhtnode.transport.send(aliveMsg)
			tempFinger = tempFingerList[i]
			break
		}
	}
	for {
		select {
		case <-dhtnode.ResponseQ:
			if dhtnode.successor.Adress != fing.Id {
				dhtnode.successor.Adress = tempFinger.Adress
				dhtnode.successor.NodeId = tempFinger.Id
			}
			return tempFinger.Adress
		case <-timeResp.C:
			fmt.Println("no resp")
			return dhtnode.findNextAlive(tempFinger)
		}
	}

}
