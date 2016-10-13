package dht

//go test -test.run TestDHT1

import (
	//"fmt"
	"testing"
	"time"
)

func TestDHT2(t *testing.T) {
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id7 := "07"
	id0 := "00"

	node0 := makeDHTNode(&id0, "localhost", "1110")
	node1 := makeDHTNode(&id1, "localhost", "1111")
	node2 := makeDHTNode(&id2, "localhost", "1112")
	node3 := makeDHTNode(&id3, "localhost", "1113")
	node4 := makeDHTNode(&id4, "localhost", "1114")
	//	node5 := makeDHTNode(nil, "localhost", "1115")
	//	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(&id7, "localhost", "1117")

	//	key1 := "2b230fe12d1c9c60a8e489d028417ac89de57635"
	//	key2 := "87adb987ebbd55db2c5309fd4b23203450ab0083"
	//	key3 := "74475501523a71c34f945ae4e87d571c2c57f6f3"

	/*node1.addToRing(node2)
	node1.addToRing(node3)
	node1.addToRing(node4)
	node4.addToRing(node5)
	node3.addToRing(node6)
	node3.addToRing(node7)
	node3.addToRing(node8)
	node7.addToRing(node9)*/

	//	fmt.Println("TEST: " + node1.lookup(key1).nodeId + " is responsible for " + key1)
	//	fmt.Println("TEST: " + node1.lookup(key2).nodeId + " is responsible for " + key2)
	//	fmt.Println("TEST: " + node1.lookup(key3).nodeId + " is responsible for " + key3)

	node1.start_server()
	node2.start_server()
	node3.start_server()
	node7.start_server()
	node0.start_server()

	src := node1.contact.ip + ":" + node1.contact.port
	//dst := node2.contact.ip + ":" + node2.contact.port
	master := &tinyNode{node1.nodeId, src}
	//node1.PrintRingProc()
	node2.join(master)
	node3.join(master)
	node0.join(master)
	node7.join(master)
	time.Sleep(time.Second * 5)

	node1.initLookUpNetworkFingers("08", node3)

	//node1.initNetworkLookUp("01", node1)
	//node1.initPrintNetworkFingers(node2)

	node4.transport.listen()

	//Glöm inte lägga till en timer på "20000sek" så inte allt dör.

}
