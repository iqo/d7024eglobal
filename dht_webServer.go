package dht

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	//"html/template"
	//"io/ioutil"
	//"log"
	"net/http"
	//"regexp"
	//"time"
	"log"
)

func (dhtnode *DHTNode) startWebServer() {
	fmt.Println("starting webserver " + dhtnode.nodeId + "ipadres" + dhtnode.transport.BindAddress)
	//timeResp := time()
	router := httprouter.New()
	ipAdressOfNode := dhtnode.transport.BindAddress

	//router.GET("/", dhtnode.Index)
	router.GET("/contains", dhtnode.NodeContainsHandler)
	log.Fatal(http.ListenAndServe(ipAdressOfNode, router))
	//http.ListenAndServe(ipAdressOfNode, router)
}
