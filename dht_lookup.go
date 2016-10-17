package dht

import (
	//"encoding/hex"
	"fmt"
	"time"
)

func (dhtnode *DHTNode) resposibleNetworkNode(key string) bool {
	if dhtnode.predecessor.nodeId == key {
		//fmt.Println("this is not know ")
		return false
	}
	if dhtnode.nodeId == key {
		//fmt.Println("this is know ")
		return true
	}

	//beeweetNodes := (between([]byte(dhtnode.predecessor.nodeId), []byte(dhtnode.nodeId), []byte(key)))
	//return beeweetNodes
	return (between([]byte(dhtnode.predecessor.nodeId), []byte(dhtnode.nodeId), []byte(key)))
}

func (dhtnode *DHTNode) findNextAlive(key int) string {
	dhtAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	//fmt.Println("dht adress:", dhtAdress, "node fingerlist adress", dhtnode.fingers.nodefingerlist[key].adress)
	notDead := AliveMessage(dhtAdress, dhtnode.fingers.nodefingerlist[key].adress)
	go dhtnode.transport.send(notDead)
	timerResp := time.NewTimer(time.Millisecond * 100)
	for {
		select {
		case r := <-dhtnode.responseQ:
			if r.Adress != "" {
				//fmt.Println("lookUp ", r.Adress)
				return r.Adress
			} else {
				return dhtnode.findNextAlive(key + 1)
			}
		case <-timerResp.C:
			fmt.Println(dhtnode.contact.port, "no response from", dhtnode.fingers.nodefingerlist[key].adress)
			if key < (bits - 1) {
				return dhtnode.findNextAlive(key + 1)
			}
		}
	}
}

func (dhtnode *DHTNode) improvedNetworkLookUp(msg *Msg) {
	dhtAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	timerResp := time.NewTimer(time.Millisecond * 100)
	if dhtnode.resposibleNetworkNode(msg.Key) {
		foundMsg := nodeFoundMessage(dhtAdress, msg.Origin, dhtAdress, dhtnode.nodeId)
		dhtnode.transport.send(foundMsg)

		for {
			select {
			case <-dhtnode.responseQ:
				return
			case <-timerResp.C:
				dhtnode.transport.send(foundMsg)
			}
		}
		return
	} else {
		//fmt.Println("fin next alive")
		next := dhtnode.findNextAlive(0)
		lookUpMsg := lookUpMessage(msg.Origin, msg.Key, dhtAdress, next)
		dhtnode.transport.send(lookUpMsg)
		//fmt.Println(dhtnode.nodeId)
	}
	return
}

func (node *DHTNode) initNetworkLookUp(key string, dhtnode *DHTNode) {
	lookUpMsg := lookUpMessage(node.transport.bindAddress, key, node.transport.bindAddress, dhtnode.transport.bindAddress)
	fmt.Println("hej")
	go func() {
		dhtnode.transport.send(lookUpMsg)
	}()
}
