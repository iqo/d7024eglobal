package dht

//go test -test.run TestDHT1

import (
	"fmt"
	"testing"
	"time"
)

func TestDHT2(t *testing.T) {
	/*id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"
	id0 := "00"

	node0 := makeDHTNode(&id0, "localhost", "1110")
	node1 := makeDHTNode(&id1, "localhost", "1111")
	node2 := makeDHTNode(&id2, "localhost", "1112")
	node3 := makeDHTNode(&id3, "localhost", "1113")
	node4 := makeDHTNode(&id4, "localhost", "1114") //listen node do not start
	node5 := makeDHTNode(&id5, "localhost", "1115")
	node6 := makeDHTNode(&id6, "localhost", "1116")
	node7 := makeDHTNode(&id7, "localhost", "1117")*/

	node0 := makeDHTNode(nil, "localhost", "1110")
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")

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
	node1.setNetworkFingers(&Msg{"", "", "", "", "", &LiteNodeStruct{node1.successor.Adress, node1.successor.NodeId}, "", ""})
	node2.start_server()
	node3.start_server()
	node7.start_server()
	node0.start_server()
	node5.start_server()
	node6.start_server()
	time.Sleep(time.Second * 3)

	src := node1.contact.ip + ":" + node1.contact.port
	//dst := node2.contact.ip + ":" + node2.contact.port
	Master := &TinyNode{node1.nodeId, src}

	//node1.PrintRingProc()

	node2.join(Master)
	node3.join(Master)
	node7.join(Master)
	node0.join(Master)
	node6.join(Master)
	node5.join(Master)
	//node1.isTheNodeAlive()
	//node1.killTheNode()
	time.Sleep(time.Second * 15)
	fmt.Println("")
	//node5.PrintOutNetworkFingers()
	//node1.initPrintNetworkFingers(node2)
	fmt.Println("")
	node2.initLookUpNetworkFinger("03", node3)

	initFileUpload(node3)
	time.Sleep(time.Second * 15)
	node5.killTheNode()
	//node5.killTheNode()
	time.Sleep(time.Second * 15)
	fmt.Println("")
	node1.initPrintNetworkFingers(node2)
	//node5.PrintOutNetworkFingers()
	fmt.Println("")
	time.Sleep(time.Second * 15)
	node5.bringNodeBack(Master)
	time.Sleep(time.Second * 15)
	node1.initPrintNetworkFingers(node2)
	time.Sleep(time.Second * 15)
	node1.initPrintNetworkFingers(node2)

	//node1.PrintOutNetworkFingers()
	//node1.isTheNodeAlive()

	//node1.initLookUpNetworkFingers("08", node3)

	//node1.initNetworkLookUp("01", node1)
	//time.Sleep(time.Second * 10)
	//node1.initPrintNetworkFingers(node2)

	node4.transport.listen()

	//Glöm inte lägga till en timer på "20000sek" så inte allt dör.

	time.Sleep(2000 * time.Second)

}
