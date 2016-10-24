package dht

import (
	"fmt"
	"html/template"
    "encoding/json"
	//"io/ioutil"
	"net/http"
	"github.com/httprouter-master" //https://github.com/julienschmidt/httprouter
	"bufio"
	"os"
	"log"
)

type Page struct {
    Address     string
}

type File struct {
	Filename	string `json="filename"`
    Data        string `json="data"`

}

type FileList struct {
    Files       []*File `json="files"`
}

/*****************************************
*** Adds content on website            ***
*****************************************/
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
    f := &Page{Address: "localhost:1115"}

    style := loadWebsite("skynet/style.html")
    script := loadWebsite("skynet/script.html")
	htmlStr := loadWebsite("skynet/webpage.html")
	t, _ := template.New("webpage").Parse(style + script + htmlStr)

    t.Execute(w, f)
}

/*****************************************
*** Opens a file and reads the content ***
*** one line at the time, the whole    ***
*** file is returned as one string     ***
*****************************************/
func loadWebsite(filename string) string{
	htmlStr := ""
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        htmlStr = htmlStr + scanner.Text() + "\n"
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return htmlStr
}

/*****************************************
*** Starts the http-server with the    ***
*** different commands (get, post etc) ***
*****************************************/
func (dhtNode *DHTNode) startWebserver() {
    router := httprouter.New()
    router.GET("/", IndexHandler)
    router.GET("/storage", GetHandler)
    router.POST("/storage", PostHandler)
    router.PUT("/storage/:KEY", PutHandler)
    router.DELETE("/storage/:KEY", DeleteHandler)

    log.Fatal(http.ListenAndServe(dhtNode.contact.ip+":"+dhtNode.contact.port, router))
}


/*****************************************
*** Gets all files in chord network    ***
*** and encodes them for the website   ***
*****************************************/
func GetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    f1 := &File{Filename: "test1", Data: "this is text for file1"}
    f2 := &File{Filename: "test2", Data: "this is text for file2"}
    f3 := &File{Filename: "test3", Data: "this is text for file3"}
    fList := &FileList{Files: []*File{f1, f2, f3}}

    jsonBody, err := json.Marshal(fList)
    w.WriteHeader(200)
    w.Write(jsonBody)

    if err != nil {
        w.WriteHeader(500)
        log.Fatal(err)
    }
}

/*****************************************
*** Handles files that are uploaded on ***
*** website and adds them on servern   ***
*** after decoding it                  ***
*****************************************/
func PostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    dec := json.NewDecoder(r.Body)
    file := File{}
    err := dec.Decode(&file)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("asdfghjklölkjhgfdsdfghjkl", file)
}

/*****************************************
*** Handles updates of text in a file  ***
*** and updates same info on server    ***
*****************************************/
func PutHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    dec := json.NewDecoder(r.Body)
    file := File{Filename: ps.ByName("KEY")}
    err := dec.Decode(&file)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("--------------------------------------->", file)
}

/*****************************************
*** Deletes a file from server         ***
*****************************************/
func DeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    //golangfunc.deletedataonfile(ps.ByName("KEY"))
}
