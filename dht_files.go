package dht

import (
	//"encoding/hex"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
	"time"
	//"io"
	//"strings"
)

func createFile(path, value string) {
	data := []byte(value)

	//fmt.Println("data: ", value)
	//fmt.Println("path: ", path)
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
		//		fmt.Println("file created no errors")
		//		fmt.Println("")
	}
}

func (dhtnode *DHTNode) createFolder() {
	path := "storage/" + dhtnode.nodeId

	if !fileAlreadyExits(path) {
		fmt.Println(dhtnode.nodeId, "does not have a folder,  creating folder", path)
		os.MkdirAll(path, 0777)
	}
}

func (dhtnode *DHTNode) upload(msg *Msg) {
	defaultPath := "storage/"
	storagePath := defaultPath + dhtnode.nodeId + "/"

	fName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	fData, _ := b64.StdEncoding.DecodeString(msg.Data)
	successornode := improvedGenerateNodeId(dhtnode.successor.Adress)
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
		fmt.Println("------	replicating message to", successornode, " ------")
		go func() { dhtnode.transport.send(replicateMsg) }()
	}
}

func (dhtnode *DHTNode) replicator(msg *Msg) {
	generatedId := improvedGenerateNodeId(msg.Origin)
	defaultPath := "storage/"

	storagePath := defaultPath + dhtnode.nodeId + "/" + generatedId + "/"
	if !fileAlreadyExits(storagePath) {
		os.MkdirAll(storagePath, 0777)
	}

	StringFileName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	StringFileData, _ := b64.StdEncoding.DecodeString(msg.Data)
	//StringFileName := b64.StdEncoding.EncodeToString([]byte(msg.FileName))
	//StringFileData := b64.StdEncoding.EncodeToString([]byte(msg.Data))

	SeconddaryStoragePath := defaultPath + dhtnode.nodeId + "/" + generatedId + "/" + string(StringFileName)

	if _, err := os.Stat(SeconddaryStoragePath); err == nil {
		os.Remove(SeconddaryStoragePath)
		//	fmt.Println("replicater err == nil")
		createFile(SeconddaryStoragePath, string(StringFileData))
	} else {

		createFile(SeconddaryStoragePath, string(StringFileData))
	}
	dhtnode.folderCheck(generatedId)

}

func (dhtnode *DHTNode) responsible(filename, data string) {
	FName, _ := b64.StdEncoding.DecodeString(filename)
	generatedHash := improvedGenerateNodeId(string(FName))
	dhtnode.initNetworkLookUp(generatedHash)
	respTimer := time.NewTimer(time.Second * 2)
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

func (node *DHTNode) disconnectedNodeResponsibility() {
	path := "storage/" + node.nodeId + "/" + node.predecessor.NodeId + "/"
	dhtnode := improvedGenerateNodeId(node.transport.BindAddress)
	thirdPath := "storage/" + node.nodeId + "/"
	fmt.Println("YOLOOOO")
	if !fileAlreadyExits(path) {
		fmt.Println(path + "does not exist")
		fmt.Println("EHEHEHEHEHHE ")

	} else {
		fmt.Println("responsibility of file from", dhtnode, "to: "+path+"\n")
		files, err := ioutil.ReadDir(path)

		//behövs detta?
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			f, _ := ioutil.ReadFile(path + file.Name())

			//this needs to be here otherwise it complainse file.Name undefined.
			secondaryPath := "storage/" + node.nodeId + "/" + file.Name()
			//-------------------------------------------
			successorNode := improvedGenerateNodeId(node.successor.Adress)

			createFile(secondaryPath, string(f))
			fmt.Println("File", file.Name(), "moved to", secondaryPath)

			fName := b64.StdEncoding.EncodeToString([]byte(file.Name()))
			fData := b64.StdEncoding.EncodeToString(f)

			fmt.Println(node.nodeId, "replicates files to:", successorNode)

			repMsg := ReplicateMessage(node.transport.BindAddress, node.successor.Adress, fName, fData)
			go func() { node.transport.send(repMsg) }()
			os.Remove(path + file.Name())
		}
		nothing, err := IsEmpty(path)
		if err != nil {
			panic(err)
		}
		if nothing {
			node.removeDirectory(thirdPath, node.predecessor.NodeId)
		} else {
			fmt.Println("PATH NOT EMPTY")
		}
	}
}

//eventuellt skriv om denna!!
//http://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty
func IsEmpty(name string) (bool, error) {
	entries, err := ioutil.ReadDir(name)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

/*    f, err := os.Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1) // Or f.Readdir(1)
    if err == io.EOF {
        return true, nil
    }
    return false, err // Either not empty or error, suits both cases
}*/

func (node *DHTNode) removeDirectory(path, directory string) {
	nodeId := improvedGenerateNodeId(node.nodeId)
	fmt.Println("node:", node, "---> ", nodeId, "removes path+directory->", path+directory)
	os.RemoveAll(path + directory)
}

//skriv färdigt denna shit
func (node *DHTNode) deleteSuccessorFile(msg *Msg) {
	nodeIdentity := improvedGenerateNodeId(msg.Origin)
	path := "storage/" + node.nodeId + "/" + nodeIdentity + "/"
	secondaryPath := "storage/" + node.nodeId + "/"
	fName, _ := b64.StdEncoding.DecodeString(msg.FileName)

	if !fileAlreadyExits(path) {
		fmt.Println("successorfile not empty")
	} else {
		fmt.Println("crashes before remove delete successor file")
		os.Remove(path + string(fName))
	}
	fmt.Println("here it crashes!! -- path is ->", path)
	nothing, err := IsEmpty(path)
	fmt.Println("NOTHING ->>", nothing)
	if err != nil {
		panic(err)
	}
	if !nothing {
		fmt.Println(path, " not empty")
	} else {
		fmt.Println("delete successor file crashes here")
		node.removeDirectory(secondaryPath, node.predecessor.NodeId)
	}
}

func (node *DHTNode) deleteSuccessorBackup(msg *Msg) {
	path := "storage/" + node.nodeId + "/"
	if !fileAlreadyExits(path) {
		fmt.Println("")
	} else {
		f, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range f {
			if file.IsDir() && file.Name() == msg.FileName {
				node.removeDirectory(path, file.Name())
			} else {
				//	fmt.Println("deleteSuccessorBackup problem")
			}
		}
	}
}

func (node *DHTNode) replicateNecessary() {
	path := "storage/" + node.nodeId + "/"

	//invertera denna funktion?
	if !fileAlreadyExits(path) {
		fmt.Println("file doesn't exist in replicate.")
	} else {
		f, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range f {
			if !file.IsDir() {
				f, _ := ioutil.ReadFile(path + file.Name())
				fName := b64.StdEncoding.EncodeToString([]byte(file.Name()))
				fData := b64.StdEncoding.EncodeToString([]byte(f))

				replicateMsg := ReplicateMessage(node.transport.BindAddress, node.successor.Adress, fName, fData)
				//fmt.Println("***********	replicating because necessary	***********")
				go func() { node.transport.send(replicateMsg) }()
			}

		}
	}
}

func (node *DHTNode) folderCheck(fileName string) {
	path := "storage/" + node.nodeId + "/" + fileName + "/"
	secondaryPath := "storage/" + node.nodeId + "/"

	if !fileAlreadyExits(path) {
		fmt.Println("file doesn't exist in foldercheck")
	} else {
		fileinfo, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, folderOne := range fileinfo {
			files, _ := ioutil.ReadDir(secondaryPath)
			for _, folderTwo := range files {
				for _, folderthree := range fileinfo {
					if !folderOne.IsDir() && folderthree.Name() == folderTwo.Name() {
						fmt.Println("removing file:", folderthree.Name(), "in", secondaryPath)
						os.Remove(secondaryPath + folderTwo.Name())
					}
				}
			}

		}
	}
}

/*func (dhtNode *DHTNode) folderCheck(nodeId string) {
	path := "storage/" + dhtNode.nodeId + "/" + nodeId + "/"
	path2 := "storage/" + dhtNode.nodeId + "/"

	if fileAlreadyExits(path) {
		dir, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, ffolderOne := range dir {
			if ffolderOne.Name() == ".DS_Store" {
				os.Remove(path + f.Name())
			} else {
				files, _ := ioutil.ReadDir(path2)
				for _, folderTwo := range files {
					for _, folderThree := range dir {
						if !f.IsDir() && folderThree.Name() == folderTwo.Name() {
							fmt.Println("Remove file: " + folderThree.Name() + " in " + path2 + "\n")
							os.Remove(path2 + folderTwo.Name())
						}
					}
				}
			}
		}
	}
}*/

func (node *DHTNode) dataFromSuccessor(msg *Msg) {
	path := "storage/" + node.nodeId + "/"
	if !fileAlreadyExits(path) {
		fmt.Println("no data from successor")
	} else {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			nothing, err := IsEmpty(path)
			if err != nil {
				panic(err)
			}

			if nothing {
				fmt.Println(path, "is empty")
			} else {
				if !file.IsDir() {
					f, _ := ioutil.ReadFile(path + file.Name())
					id := improvedGenerateNodeId(file.Name())

					if between([]byte(msg.LiteNode.Adress), []byte(node.nodeId), []byte(id)) {
						fmt.Println("")
					} else {
						fName := b64.StdEncoding.EncodeToString([]byte(file.Name()))
						fData := b64.StdEncoding.EncodeToString([]byte(f))

						upload := UpLoadMessage(node.transport.BindAddress, msg.LiteNode.Adress, fName, fData)
						go func() { node.transport.send(upload) }()

						deletefile := deleteFileMessage(node.transport.BindAddress, node.successor.Adress, fName)
						go func() { node.transport.send(deletefile) }()
						os.Remove(path + file.Name())

					}
				}
			}
		}
	}
}
