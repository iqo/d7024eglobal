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
	directorys, err := ioutil.ReadDir("storage/")
	adress := dht.contact.ip + ":" + dht.contact.port
	genAdress := improvedGenerateNodeId(adress)

	if err != nil {
		log.Fatal(err)
	}

	for _, tempDirectory := range directorys {
		if genAdress == tempDirectory.Name() {
			fmt.Fprint(w, "has folder ", tempDirectory.Name(), "\n")
			tempSecondDir, _ := ioutil.ReadDir("storage/" + tempDirectory.Name() + "/")
			for _, fileInDir := range tempSecondDir {
				fmt.Fprint(w, "this folder contains: ", fileInDir.Name())
			}
		}
	}
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
