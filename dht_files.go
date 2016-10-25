package dht

import (
	//"encoding/hex"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
	"time"
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
	//_, err := os.Stat(name)
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

/*func (dhtnode *DHTNode) uploadFile(filePath, key, value string) {
	if fileAlreadyExits(filePath) != true {
		os.Mkdir(filePath, 0777)
	}
	createFile(filePath+key, value)
}*/

func errorChecker(e error) {
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println("file created no errors")
	}
}

func (dhtnode *DHTNode) createFolder() {
	path := "storage/" + dhtnode.nodeId
	fmt.Println(dhtnode.nodeId, "does not have a folder,  creating folder", path)
	if !fileAlreadyExits(path) {
		os.MkdirAll(path, 0777)
	}
}

func (dhtnode *DHTNode) upload(msg *Msg) {
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

func (dhtnode *DHTNode) replicator(msg *Msg) {
	generatedId := improvedGenerateNodeId(msg.Origin)
	defaultPath := "storage/"

	storagePath := defaultPath + dhtnode.nodeId + "/" + generatedId + "/"
	if !fileAlreadyExits(storagePath) {
		os.MkdirAll(storagePath, 077)
	}

	StringFileName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	StringFileData, _ := b64.StdEncoding.DecodeString(msg.Data)
	//StringFileName := b64.StdEncoding.EncodeToString([]byte(msg.FileName))
	//StringFileData := b64.StdEncoding.EncodeToString([]byte(msg.Data))

	SeconddaryStoragePath := defaultPath + "/" + dhtnode.nodeId + "/" + generatedId + string(StringFileName)

	_, err := os.Stat(SeconddaryStoragePath)
	if err == nil {
		os.Remove(SeconddaryStoragePath)
		createFile(SeconddaryStoragePath, string(StringFileData))
	} else {
		createFile(SeconddaryStoragePath, string(StringFileData))
	}

}

func (dhtnode *DHTNode) responsible(filename, data string) {
	respTimer := time.NewTimer(time.Second * 2)
	FName, _ := b64.StdEncoding.DecodeString(filename)
	generatedHash := improvedGenerateNodeId(string(FName))
	dhtnode.initNetworkLookUp(generatedHash)
	for {
		select {
		case fingerResp := <-dhtnode.FingerQ:
			fmt.Println("uploading file to folder", fingerResp.Id)
			upLoadMsg := UpLoadMessage(dhtnode.transport.BindAddress, fingerResp.Adress, filename, data)
			go func() { dhtnode.transport.send(upLoadMsg) }()
			return
		case <-respTimer.C:
			return
		}
	}
}

func initFileUpload(dhtnode *DHTNode) {
	//filePath := "C:\Users\Niklas\gocode\github\Mox_D7024E\github.com\d7024eglobal\readme" //fuck windows
	//filePath := "C/Users/Niklas/gocode/github/Mox_D7024E/github.com/d7024eglobal/readme"
	filePath := "readme/"
	filesInPath, err := ioutil.ReadDir(filePath)
	if err != nil {
		panic(err)
	}

	for _, temp := range filesInPath {
		readFile, _ := ioutil.ReadFile(filePath + temp.Name())

		stringFile := b64.StdEncoding.EncodeToString([]byte(temp.Name()))
		stringData := b64.StdEncoding.EncodeToString(readFile)

		dhtnode.responsible(stringFile, stringData)
	}
}
