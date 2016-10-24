package dht

/*import (
	//"encoding/hex"
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

func (dhtnode *DHTNode) initUpload(msg *Msg) {
	if dhtnode.resposibleNetworkNode(msg.key) {
		data := strings.Split(msg.Bytes, ";")
		dataInFile := ""
		for _, tempData := range data[1:] {
			dataInFile = dataInFile + tempData
		}
		dhtnode.uploadFile(filePath, key, value)

	}
}
*/
