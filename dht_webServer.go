package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Page struct {
	PageAress string
}

type FileProperties struct {
	Filename string
	FileData string
}

type FileList struct {
	AllFilesInList []*FileProperties
}
