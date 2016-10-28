package dht

import (
	"fmt"
	"time"
)

func (dhtnode *DHTNode) heartBeat() {
	nodeAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	heartMsg := heartBeatMessage(nodeAdress, dhtnode.predecessor.Adress)
	//fmt.Println(dhtnode.predecessor.nodeId, "has adress ", dhtnode.predecessor.adress)
	waitTimer := time.NewTimer(time.Second * 1)
	go dhtnode.transport.send(heartMsg)
	for {
		select {
		case <-dhtnode.HeartBeatQ:
			//fmt.Println("stil alive baby", dhtnode.predecessor.adress)
			return

		case <-waitTimer.C:
			fmt.Println("heartstop", dhtnode.contact.port)
			dhtnode.disconnectedNodeResponsibility()
			dhtnode.predecessor.Adress = ""
			dhtnode.predecessor.NodeId = ""
			dhtnode.createNewTask(nil, "stabilize")
			return
		}
	}
}

func (dhtnode *DHTNode) heartTimer() {
	for {
		if dhtnode.alive {
			//fmt.Println("heart8 timer")
			time.Sleep(time.Second * 1)
			go dhtnode.createNewTask(nil, "heartBeat")
		} else {
			return
		}
	}
}