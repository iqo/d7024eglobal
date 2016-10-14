package dht

import (
	"fmt"
	"time"
)

func (dhtnode *DHTNode) heartBeat() {
	//nodeAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	//heartMsg := heartBeatMessage(nodeAdress, dhtnode.predecessor.adress)
	fmt.Println(dhtnode.predecessor.nodeId, "has adress ", dhtnode.predecessor.adress)
	/*go func() { dhtnode.transport.send(heartMsg) }()
	waitTimer := time.NewTimer(time.Second * 4)
	for {
		select {
		case <-dhtnode.heartBeatQ:
			fmt.Println("stil alive baby", dhtnode.predecessor.adress)
			return

		case <-waitTimer.C:
			fmt.Println("heartstop", dhtnode.contact.port)
			dhtnode.predecessor.adress = ""
			dhtnode.predecessor.nodeId = ""
			dhtnode.stabilize()
			return
		}
	}*/
}

func (dhtnode *DHTNode) heartTimer() {
	for {
		//fmt.Println("heart timer")
		time.Sleep(time.Second * 3)
		dhtnode.createNewTask(nil, "heartBeat")
	}
}
