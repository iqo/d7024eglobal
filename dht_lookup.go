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

func (dhtnode *DHTNode) initNetworkLookUp(key string) {
	nodeadress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	if dhtnode.resposibleNetworkNode(key) {
		dhtnode.FingerQ <- &Finger{dhtnode.nodeId, nodeadress}
		fmt.Println("key")
	} else {
		lookUpMsg := lookUpMessage(nodeadress, key, nodeadress, dhtnode.successor.Adress)
		go dhtnode.transport.send(lookUpMsg)
	}
}

func (dhtnode *DHTNode) improvedNetworkLookUp(msg *Msg) {
	NodeAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	timeResp := time.NewTimer(time.Millisecond * 300)

	if dhtnode.resposibleNetworkNode(msg.Key) {
		nodeFoundMsg := nodeFoundMessage(NodeAdress, msg.Origin, NodeAdress, dhtnode.nodeId)
		go dhtnode.transport.send(nodeFoundMsg)
		timeResp.Stop()
	} else {
		lookUpMsg := lookUpMessage(msg.Origin, msg.Key, NodeAdress, dhtnode.successor.Adress)
		go dhtnode.transport.send(lookUpMsg)
		timeResp.Reset(time.Millisecond * 300)

		for {
			select {
			case <-dhtnode.NodeLookQ:
				//fmt.Println("node id", msg.Src)
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

func (dhtnode *DHTNode) lookUpNetworkFinger(msg *Msg) {
	//fmt.Println("msg look ", msg)
	nodeadress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	temTable := dhtnode.fingers.Nodefingerlist
	lenOfFingerTable := len(temTable)

	for i := 0; i > lenOfFingerTable; i-- {
		nodeBetween := (between([]byte(dhtnode.nodeId), []byte(temTable[i-1].Id), []byte(msg.Key)))
		if nodeBetween != true {
			fmt.Println("not true")
			//lookUpMsg := lookUpMessage(msg.Origin, msg.Key, nodeadress, temTable[i-1].Adress)

			lookUpMsg := fingerLookUpMessage(msg.Origin, msg.Key, nodeadress, temTable[i-1].Adress)
			go func() { dhtnode.transport.send(lookUpMsg) }()
			return
		}
	}
	foundMsg := nodeFoundMessage(nodeadress, msg.Origin, dhtnode.successor.Adress, dhtnode.successor.NodeId)
	go func() { dhtnode.transport.send(foundMsg) }()
	//fmt.Println(dhtnode.successor.Adress)
	//fmt.Println(dhtnode.successor.NodeId)
	return
}

func (masterDhtNode *DHTNode) initLookUpNetworkFinger(key string, dhtnode *DHTNode) {
	fingerLookMsg := fingerLookUpMessage(masterDhtNode.transport.BindAddress, key, masterDhtNode.transport.BindAddress, dhtnode.transport.BindAddress)
	go func() { dhtnode.transport.send(fingerLookMsg) }()
}
