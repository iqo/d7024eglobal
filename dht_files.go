package dht

import (
	//"encoding/hex"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func createFile(path, value, string) {
	data := []byte(value)
	fmt.Println("data: ", value)
	fmt.Println("path: ", path)
	err := ioutil.WriteFile(path, data, 0777)
	check(err)
}

func fileAlreadyExits(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (dhtnode *DHTNode) uploadFile(filePath, key, value string) {
	if fileAlreadyExits(path) != true {
		os.Mkdir(path, 0777)
	}
	createFile(filePath+key, value)
}

func errorChecker(e error) {
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println("file created no errors")
	}
}

func (dhtnode *DHTNode) createFolder() {
	path := "folder/" + dhtnode.nodeId
	if !fileAlreadyExits(path) {
		os.MkdirAll(path, 0777)
	}
}

func (dhtnode *DHTNode) initUpload(msg *Msg) {
	generatedNodeid := improvedGenerateNodeId(msg.Dst)
	nodeIdForSuccessor := improvedGenerateNodeId(dhtnode.successor.Adress)
	desiredPath := "storage/" + generatedNodeid + "/"

	if fileAlreadyExits(desiredPath) != true {
		os.Mkdir(desiredPath, 077)
	}
	FName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	FData, _ := b64.StdEncoding.DecodeString(msg.Bytes)
	desiredPath = "storage/" + generatedNodeid + "/" + string(FName)
	createFile(desiredPath, string(FName))
	fNameEncodedToString :=
	fDataEncodedToString :=

}
