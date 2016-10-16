package dht

import (
	//"encoding/hex"
	"fmt"
	"time"
)

//const bits int = 3

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *tinyNode
	predecessor *tinyNode
	contact     Contact
	fingers     *FingerTable
	transport   *Transport
	responseQ   chan *Msg
	TaskQ       chan *Task
	heartBeatQ  chan *Msg
	fingerQ     chan *Finger
	nodeQ       chan *Msg

	alive bool
}
type tinyNode struct {
	nodeId string
	adress string
}

type Task struct {
	message *Msg
	Type    string
}

func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port

	if nodeId == nil {
		genNodeId := generateNodeId()
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}

	dhtNode.successor = &tinyNode{dhtNode.nodeId, ip + ":" + port}
	dhtNode.predecessor = &tinyNode{dhtNode.nodeId, ip + ":" + port}
	dhtNode.fingers = &FingerTable{}
	//ska new användas eller raden under?
	//dhtNode.fingers.nodefingerlist = [bits]*DHTNode{}
	//eller denna kanske
	//dhtNode.fingers = &FingerTable{}
	dhtNode.alive = true
	dhtNode.responseQ = make(chan *Msg)
	dhtNode.TaskQ = make(chan *Task)
	dhtNode.heartBeatQ = make(chan *Msg)
	dhtNode.fingerQ = make(chan *Finger)
	dhtNode.nodeQ = make(chan *Msg)
	dhtNode.createTransport()
	return dhtNode
}

func (dhtNode *DHTNode) createTransport() {
	dhtNode.transport = &Transport{dhtNode, dhtNode.contact.ip + ":" + dhtNode.contact.port, nil}
	dhtNode.transport.msgQ = make(chan *Msg)
	dhtNode.transport.initmsgQ()
}

func (dhtNode *DHTNode) join(master *tinyNode) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	message := message("join", src, master.adress, src, dhtNode.nodeId, nil)
	dhtNode.transport.send(message)
	for {
		select {
		case r := <-dhtNode.responseQ:
			dhtNode.successor.adress = r.Src
			dhtNode.successor.nodeId = r.Key
			dhtNode.setNetworkFingers(&Msg{"", "", "", "", nil, dhtNode.successor.adress, dhtNode.successor.nodeId, ""})
			fingerStart := fingerStartMessage(src, dhtNode.successor.adress, dhtNode.transport.bindAddress, dhtNode.nodeId)
			go func() { dhtNode.transport.send(fingerStart) }()
			return
			//fmt.Println(dhtNode.nodeId, dhtNode.successor)
		}
	}
}

//Ligger du mellan noderna, nej. skicka join msg till nästa nod och kolla
//om han ligger mellan den noden och hans successor.
func (node *DHTNode) findSucc(msg *Msg) {
	var a = between([]byte(node.nodeId), []byte(node.successor.nodeId), []byte(msg.Key))
	if a {

		node.transport.send(message("response", msg.Dst, msg.Origin, node.successor.adress, node.successor.nodeId, nil))
		node.successor.adress = msg.Origin
		node.successor.nodeId = msg.Key
	} else {
		node.transport.send(message("join", msg.Origin, node.successor.adress, msg.Dst, msg.Key, nil))

	}
}

func (node *DHTNode) printNetworkRing(msg *Msg) {
	if msg.Origin != msg.Dst {

		fmt.Println(node.nodeId, node.successor)
		node.transport.send(printMessage(msg.Origin, node.successor.adress))
	}
}

func (dhtNode *DHTNode) start_server() {
	go dhtNode.initTaskQ()
	go dhtNode.stableTimmer()
	go dhtNode.fingerTimer()
	go dhtNode.heartTimer()
	go dhtNode.transport.listen()
}

func (dhtNode *DHTNode) notifyNetwork(msg *Msg) {
	if (dhtNode.predecessor.adress == "") || between([]byte(dhtNode.predecessor.nodeId), []byte(dhtNode.nodeId), []byte(msg.Key)) {
		dhtNode.predecessor.adress = msg.Src
		dhtNode.predecessor.nodeId = msg.Key
	}
}

func (node *DHTNode) initTaskQ() {
	go func() {
		for {
			select {
			case t := <-node.TaskQ:
				switch t.Type {
				case "printRing": //test case
					node.printNetworkRing(t.message)
					//node.improvePrintRing(node.msg)
					//transport.send(&Msg{"printRing", "", v.Src, []byte("tjuuu")})
				case "join":
					node.findSucc(t.message)
				case "stabilize":
					//			fmt.Println("stabilize case: ", node.nodeId)
					node.stabilize()
				case "updateFingers":
					node.updateNetworkFingers()
				case "heartBeat":
					fmt.Println("initTask hearbeat")
					node.heartBeat()
				case "alive":
					fmt.Println("fuck")
				}
			}
		}
	}()
}

func (node *DHTNode) stabilize() {
	nodeAdress := node.contact.ip + ":" + node.contact.port
	predOfSucc := getNodeMessage(nodeAdress, node.successor.adress) // id eller adress?
	go node.transport.send(predOfSucc)
	time := time.NewTimer(time.Millisecond * 5000)
	for {
		select {
		case r := <-node.responseQ:
			//fmt.Println("case 1 stab: ")

			between := (between([]byte(node.nodeId), []byte(node.successor.nodeId), []byte(r.Key))) && r.Key != "" //r.key = "" för att connecta sista nodens successor
			if between {
				node.successor.adress = r.Src //origin eller source
				//node.successor.adress = msg.Origin
				//node.successor.nodeId = msg.Key
				node.successor.nodeId = r.Key
				//	fmt.Println("beetween")
				//return
			}
			//ska notifymessage ha fler variabler?
			N := notifyMessage(nodeAdress, node.successor.adress, nodeAdress, node.nodeId)

			go node.transport.send(N)
			//	fmt.Println("node id:", node.nodeId, "node successor id:", node.successor, "node predecessor id:", node.predecessor)
			return
		case timer := <-time.C: //timer
			fmt.Println("TIMER ERROR:", timer)
			node.updateSucc(1)
			return
		}
	}
}

func (dhtnode *DHTNode) stableTimmer() {
	for {
		if dhtnode.alive {
			time.Sleep(time.Millisecond * 5000)
			dhtnode.createNewTask(nil, "stabilize")
		}
	}
}

func (node *DHTNode) createNewTask(msg *Msg, typeOfTask string) {
	if node.alive {
		task := &Task{msg, typeOfTask}
		node.TaskQ <- task
	}
}

func (node *DHTNode) setSucc(msg *Msg) {
	node.successor.adress = msg.Src
	node.successor.nodeId = msg.Key
}

func (node *DHTNode) setPred(msg *Msg) {
	node.predecessor.adress = msg.Src
	node.predecessor.nodeId = msg.Key
}

func (node *DHTNode) getPred(msg *Msg) {
	//fmt.Println("hej getpred")
	//fmt.Println("src:",msg.Dst,"dst:", msg.Src,"node pred adress:", node.predecessor.adress,"node pred. node id:", node.predecessor.nodeId)
	//fmt.Println("dst:", msg.Src)
	//fmt.Println("node pred adress:", node.predecessor.adress)
	//fmt.Println("node pred node id:", node.predecessor.nodeId)

	responseMsg := responseMessage(msg.Dst, msg.Src, node.predecessor.adress, node.predecessor.nodeId)

	go func() {
		node.transport.send(responseMsg)
	}()
}

func (node *DHTNode) PrintRingProc() {
	src := node.contact.ip + ":" + node.contact.port
	go func() {
		for {
			time.Sleep(time.Second * 2)
			fmt.Println()
			node.TaskQ <- &Task{printMessage(src, ""), "printRing"}
		}
	}()
}

/*func (dhtnode *DHTNode) networkLookup(msg *Msg) {
	nodeAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port

	if between([]byte(dhtnode.nodeId), []byte(dhtnode.successor.nodeId), []byte(msg.Key)) {
		if dhtnode.nodeId == msg.Key {
			//fmt.Println(dhtnode.nodeId)
			respMsg := responseMessage(nodeAdress, msg.Origin, nodeAdress, dhtnode.nodeId)
			go func() { dhtnode.transport.send(respMsg) }()
			//return
		} else {
			//fmt.Println(dhtnode.successor.nodeId)
			respMsg := responseMessage(nodeAdress, msg.Origin, dhtnode.successor.adress, dhtnode.successor.nodeId)
			go func() { dhtnode.transport.send(respMsg) }()
			//return
		}
	} else {
		//fmt.Println("lookup else ")
		lookUpMsg := lookUpMessage(msg.Origin, msg.Key, nodeAdress, dhtnode.successor.adress)
		go func() { dhtnode.transport.send(lookUpMsg) }()
	}
	//fmt.Println(dhtnode.successor.nodeId)
}

//skicka till taskQ!!!
func (node *DHTNode) initNetworkLookUp(key string, dhtnode *DHTNode) {
	lookUpMsg := lookUpMessage(node.transport.bindAddress, key, node.transport.bindAddress, dhtnode.transport.bindAddress)
	fmt.Println("hej")
	go func() {
		dhtnode.transport.send(lookUpMsg)
	}()
}*/

func (dhtnode *DHTNode) killTheNode() {
	fmt.Println("killing node ", dhtnode.nodeId)
	dhtnode.alive = false
	dhtnode.successor.nodeId = ""
	dhtnode.successor.adress = ""
	dhtnode.predecessor.adress = ""
	dhtnode.predecessor.nodeId = ""
}

/*func (dhtnode *DHTNode) isTheNodeAlive() bool {
	if dhtnode.alive {
		return true
	} else {
		return false
	}
}*/

func (dhtnode *DHTNode) updateSucc(key int) {
	dhtAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	getPredOfFinger := getNodeMessage(dhtAdress, dhtnode.fingers.nodefingerlist[key].adress)
	go dhtnode.transport.send(getPredOfFinger)

	timerResp := time.NewTimer(time.Second * 1)
	for {
		select {
		case r := <-dhtnode.responseQ:
			if r.Id == "" {
				dhtnode.successor.adress = dhtnode.fingers.nodefingerlist[key].adress
				dhtnode.successor.nodeId = dhtnode.fingers.nodefingerlist[key].id
				fmt.Println("update succ done")
			}
			notify := notifyMessage(dhtAdress, dhtnode.fingers.nodefingerlist[key].adress, dhtAdress, dhtnode.nodeId)
			go dhtnode.transport.send(notify)
			return

		case <-timerResp.C:
			if key < (bits - 1) {
				dhtnode.updateSucc(key + 1)
			}
			return
		}
	}
}
