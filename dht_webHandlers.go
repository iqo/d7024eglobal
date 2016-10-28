package dht

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
)

type File struct {
	FileName string
	Data     string
}

type FileList struct {
	AllFiles []*File
}

func (dht *DHTNode) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "this is ", dht.nodeId, " and has ipadres ", dht.transport.BindAddress)
}

func (dht *DHTNode) NodeContainsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//	directorys, err := ioutil.ReadDir("storage/")
	adress := dht.contact.ip + ":" + dht.contact.port
	genAdress := improvedGenerateNodeId(adress)
	directorys, err := ioutil.ReadDir("storage/" + genAdress + "/")

	if err != nil {
		log.Fatal(err)
	}

	for _, tempDirectory := range directorys {
		if tempDirectory.IsDir() {
			fmt.Fprint(w, dht.transport.BindAddress, " contains backup for ", tempDirectory.Name())
		} else {
			fmt.Fprint(w, dht.transport.BindAddress, " contains the file ", tempDirectory.Name())
		}
		//fmt.Fprint(w, dht.transport.BindAddress," contains "tempDirectory.Name())

	}
}

func (dhtNode *DHTNode) UpdateKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("_______________________________________________________________________")
	fmt.Println(p.ByName("key"))
	fmt.Println("_______________________________________________________________________")
	fmt.Fprint(w, p.ByName("key"))
}

/*func (dht *DHTNode) viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}*/

/*func loadFile(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}*/

func (dht *DHTNode) NodeContainsFileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//	directorys, err := ioutil.ReadDir("storage/")
	//adress := dht.contact.ip + ":" + dht.contact.port
	//genAdress := improvedGenerateNodeId(adress)
	directorys, err := ioutil.ReadDir("storage/")

	if err != nil {
		log.Fatal(err)
	}

	for _, tempDirectory := range directorys {

		folderNodeAdress, err := ioutil.ReadDir("storage/" + tempDirectory.Name())
		if err != nil {
			log.Fatal(err)
		}
		for _, sencondTemp := range folderNodeAdress {
			if !sencondTemp.IsDir() {
				fmt.Fprint(w, "system contains these files\n")
				fmt.Fprint(w, sencondTemp.Name(), "\n")
			}
		}
	}
}
