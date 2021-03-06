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
	successor   *TinyNode
	predecessor *TinyNode
	contact     Contact
	fingers     *FingerTable
	transport   *Transport
	ResponseQ   chan *Msg
	TaskQ       chan *Task
	HeartBeatQ  chan *Msg
	FingerQ     chan *Finger
	NodeLookQ   chan *Msg
	FileQ       chan *File
	//Path        string
	alive bool
}
type TinyNode struct {
	NodeId string
	Adress string
}

type Task struct {
	Message *Msg
	Type    string
}

func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port

	if nodeId == nil {
		adress := ip + ":" + port
		genNodeId := improvedGenerateNodeId(adress)
		//genNodeId := generateNodeId()
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}

	dhtNode.successor = &TinyNode{dhtNode.nodeId, ip + ":" + port}
	dhtNode.predecessor = &TinyNode{dhtNode.nodeId, ip + ":" + port}
	dhtNode.fingers = &FingerTable{}
	//ska new användas eller raden under?
	//dhtNode.fingers.Nodefingerlist = [bits]*DHTNode{}
	//eller denna kanske
	//dhtNode.fingers = &FingerTable{}
	dhtNode.alive = true
	dhtNode.ResponseQ = make(chan *Msg)
	dhtNode.TaskQ = make(chan *Task)
	dhtNode.HeartBeatQ = make(chan *Msg)
	dhtNode.FingerQ = make(chan *Finger)
	dhtNode.NodeLookQ = make(chan *Msg)
	dhtNode.FileQ = make(chan *File)
	dhtNode.createTransport()
	dhtNode.createFolder()
	return dhtNode
}

func (dhtNode *DHTNode) createTransport() {
	dhtNode.transport = &Transport{dhtNode, dhtNode.contact.ip + ":" + dhtNode.contact.port, nil}
	dhtNode.transport.msgQ = make(chan *Msg)
	dhtNode.transport.initmsgQ()
}

func (dhtNode *DHTNode) join(master *TinyNode) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	message := message("join", src, master.Adress, src, dhtNode.nodeId, "")
	dhtNode.transport.send(message)
	for {
		select {
		case r := <-dhtNode.ResponseQ:
			dhtNode.successor.Adress = r.Src
			dhtNode.successor.NodeId = r.Key

			dhtNode.setNetworkFingers(&Msg{"", "", "", "", "", &LiteNodeStruct{dhtNode.successor.Adress, dhtNode.successor.NodeId}, "", ""})
			fingerStart := fingerStartMessage(src, dhtNode.successor.Adress, dhtNode.transport.BindAddress, dhtNode.nodeId)
			fmt.Println("fingerstart join: ", fingerStart)
			go func() { dhtNode.transport.send(fingerStart) }()
			return
			//fmt.Println(dhtNode.nodeId, dhtNode.successor)
		}
	}
}

//Ligger du mellan noderna, nej. skicka join msg till nästa nod och kolla
//om han ligger mellan den noden och hans successor.
func (node *DHTNode) findSucc(msg *Msg) {
	var a = between([]byte(node.nodeId), []byte(node.successor.NodeId), []byte(msg.Key))
	if a {

		go node.transport.send(message("response", msg.Dst, msg.Origin, node.successor.Adress, node.successor.NodeId, ""))
		node.successor.Adress = msg.Origin
		node.successor.NodeId = msg.Key
	} else {
		node.transport.send(message("join", msg.Origin, node.successor.Adress, msg.Dst, msg.Key, ""))

	}
}

func (node *DHTNode) printNetworkRing(msg *Msg) {
	if msg.Origin != msg.Dst {

		fmt.Println(node.nodeId, node.successor)
		node.transport.send(printMessage(msg.Origin, node.successor.Adress))
	}
}

func (dhtNode *DHTNode) start_server() {
	go dhtNode.initTaskQ()
	go dhtNode.stableTimmer()
	go dhtNode.fingerTimer()
	go dhtNode.heartTimer()
	go dhtNode.transport.listen()
	go dhtNode.startWebServer()
}

func (dhtNode *DHTNode) notifyNetwork(msg *Msg) {
	if (dhtNode.predecessor.Adress == "") || between([]byte(dhtNode.predecessor.NodeId), []byte(dhtNode.nodeId), []byte(msg.LiteNode.Id)) {
		dhtNode.predecessor.Adress = msg.LiteNode.Adress
		dhtNode.predecessor.NodeId = msg.LiteNode.Id
	}
}

func (node *DHTNode) initTaskQ() {
	go func() {
		for {
			select {
			case t := <-node.TaskQ:
				switch t.Type {
				case "printRing": //test case
					node.printNetworkRing(t.Message)
					//node.improvePrintRing(node.msg)
					//transport.send(&Msg{"printRing", "", v.Src, []byte("tjuuu")})
				case "join":
					go node.findSucc(t.Message)
				case "stabilize":
					//			fmt.Println("stabilize case: ", node.nodeId)
					//go node.stabilize()
					node.stabilize()
				case "updateFingers":
					go node.updateNetworkFingers()
					//go node.updateNetworkFingers()
				case "heartBeat":
					//fmt.Println("initTask hearbeat")
					//go node.heartBeat()
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
	predOfSucc := getPredMessage(nodeAdress, node.successor.Adress) // id eller adress?
	go node.transport.send(predOfSucc)
	time := time.NewTimer(time.Millisecond * 2000)
	for {
		select {
		case r := <-node.ResponseQ:
			//fmt.Println("case 1 stab: ")

			between := (between([]byte(node.nodeId), []byte(node.successor.NodeId), []byte(r.LiteNode.Id))) && r.LiteNode.Id != "" && node.nodeId != r.LiteNode.Id //r.key = "" för att connecta sista nodens successor
			if between {
				node.successor.Adress = r.LiteNode.Adress //origin eller source
				node.successor.NodeId = r.LiteNode.Id
				//	fmt.Println("beetween")
				//return
			}
			//ska notifymessage ha fler variabler?
			notifyMsg := notifyMessage(nodeAdress, node.successor.Adress, nodeAdress, node.nodeId)

			go node.transport.send(notifyMsg)
			//	fmt.Println("node id:", node.nodeId, "node successor id:", node.successor, "node predecessor id:", node.predecessor)
			return
		case timer := <-time.C: //timer
			fmt.Println("Stabilize timeout error updating suscessor:", timer)
			node.updateSucc(node.successor.NodeId)
			return
		}
	}
}

func (dhtnode *DHTNode) stableTimmer() {
	for {
		if dhtnode.alive {
			time.Sleep(time.Millisecond * 2000)
			go dhtnode.createNewTask(nil, "stabilize")
		} else {
			return
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
	node.successor.Adress = msg.Src
	node.successor.NodeId = msg.Key
}

func (node *DHTNode) setPred(msg *Msg) {
	node.predecessor.Adress = msg.Src
	node.predecessor.NodeId = msg.Key
}

func (node *DHTNode) getPred(msg *Msg) {
	responseMsg := responseMessage(msg.Dst, msg.Origin, node.predecessor.Adress, node.predecessor.NodeId)

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

func (dhtnode *DHTNode) killTheNode() {
	fmt.Println("killing node ", dhtnode.nodeId)
	dhtnode.alive = false
	/*dhtnode.successor.NodeId = ""
	dhtnode.successor.Adress = ""
	dhtnode.predecessor.Adress = ""
	dhtnode.predecessor.NodeId = ""*/
}

func (dhtnode *DHTNode) updateSucc(key string) {
	tempFinger := &Finger{dhtnode.successor.NodeId, dhtnode.successor.Adress}
	dhtAdress := dhtnode.contact.ip + ":" + dhtnode.contact.port
	//getPredOfFinger := getNodeMessage(dhtAdress, dhtnode.fingers.Nodefingerlist[key].adress)
	lenOfFingerList := len(dhtnode.fingers.Nodefingerlist)
	for i := 0; i < lenOfFingerList; i++ {
		if dhtnode.fingers.Nodefingerlist[i].Id > key {
			FingerPred := getPredMessage(dhtAdress, dhtnode.fingers.Nodefingerlist[i].Adress)
			go dhtnode.transport.send(FingerPred)
			tempFinger = dhtnode.fingers.Nodefingerlist[i]
			break
		}
	}
	timerResp := time.NewTimer(time.Millisecond * 500)
	for {
		select {
		case <-dhtnode.ResponseQ:
			dhtnode.successor.Adress = tempFinger.Adress
			dhtnode.successor.NodeId = tempFinger.Id
			notifyMsg := notifyMessage(dhtAdress, tempFinger.Adress, dhtAdress, dhtnode.nodeId)
			go dhtnode.transport.send(notifyMsg)
			return

		case <-timerResp.C:
			dhtnode.updateSucc(tempFinger.Id)
			return
		}
	}
}

func (dhtnode *DHTNode) bringNodeBack(master *TinyNode) {
	src2 := dhtnode.contact.ip + ":" + dhtnode.contact.port
	if dhtnode.alive == false {
		dhtnode.alive = true

		go dhtnode.initTaskQ()
		go dhtnode.stableTimmer()
		go dhtnode.fingerTimer()
		go dhtnode.heartTimer()
		go dhtnode.transport.listen()
		//go dhtNode.startWebServer()
		//dhtnode.start_server()
		dhtnode.successor.NodeId = dhtnode.nodeId
		dhtnode.successor.Adress = src2
		master.NodeId = dhtnode.nodeId
		master.Adress = src2
		dhtnode.join(master)
		fmt.Println("node ", dhtnode.nodeId, "rejoining the ring ")
	}
}
