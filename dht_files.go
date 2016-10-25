package dht

import (
	//"encoding/hex"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
	//"strings"
)

func createFile(path, value string) {
	data := []byte(value)
	fmt.Println("data: ", value)
	fmt.Println("path: ", path)
	err := ioutil.WriteFile(path, data, 0777)
	errorChecker(err)
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
	if fileAlreadyExits(filePath) != true {
		os.Mkdir(filePath, 0777)
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
	defaultPath := "storage/"
	storagePath := defaultPath + dhtnode.nodeId + "/"

	fName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	fData, _ := b64.StdEncoding.DecodeString(msg.Data)

	generatedHash := improvedGenerateNodeId(string(fName))

	if dhtnode.resposibleNetworkNode(generatedHash) != true {
		StringFileName := b64.StdEncoding.EncodeToString([]byte(fName))
		StringFileData := b64.StdEncoding.EncodeToString([]byte(fData))
		uploadMsg := UpLoadMessage(dhtnode.transport.BindAddress, dhtnode.predecessor.Adress, StringFileName, StringFileData)
		go func() { dhtnode.transport.send(uploadMsg) }()

	} else {
		if !fileAlreadyExits(storagePath) {
			os.MkdirAll(storagePath, 0777)
		}
		storagePath = defaultPath + dhtnode.nodeId + "/" + string(fName)
		createFile(storagePath, string(fData))

		tempStringFileName := b64.StdEncoding.EncodeToString(fName)
		tempStringFileData := b64.StdEncoding.EncodeToString(fData)

		replicateMsg := ReplicateMessage(dhtnode.transport.BindAddress, dhtnode.successor.Adress, tempStringFileName, tempStringFileData)
		go func() { dhtnode.transport.send(replicateMsg) }()
	}
}
