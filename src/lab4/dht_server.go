package dht

import (
	"fmt"
	"html/template"
    "encoding/json"
    "time"
	"net/http"
	"github.com/httprouter-master" //https://github.com/julienschmidt/httprouter
	"bufio"
	"os"
	"log"
    b64 "encoding/base64"
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
func (dhtNode *DHTNode) IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
    p := &Page{Address: dhtNode.transport.bindAddress}

    style := loadWebsite("skynet/style.html")
    script := loadWebsite("skynet/script.html")
	htmlStr := loadWebsite("skynet/webpage.html")
	t, _ := template.New("webpage").Parse(style + script + htmlStr)

    t.Execute(w, p)
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
    router.GET("/", dhtNode.IndexHandler)
    router.GET("/storage", dhtNode.GetHandler)
    router.POST("/storage", dhtNode.PostHandler)
    router.PUT("/storage/:KEY", dhtNode.PutHandler)
    router.DELETE("/storage/:KEY", dhtNode.DeleteHandler)

    log.Fatal(http.ListenAndServe(dhtNode.contact.ip+":"+dhtNode.contact.port, router))
}

/*****************************************
*** Gets all files in chord network    ***
*** and encodes them for the website   ***
*****************************************/
func (dhtNode *DHTNode) GetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var fileList []*File
    var newFileList []*File
    var done = false

    // sends message to get all existing files on system
    go dhtNode.transport.send(createGetFilesMsg(dhtNode.transport.bindAddress, dhtNode.transport.bindAddress, dhtNode.transport.bindAddress))
    waitResponse := time.NewTimer(time.Millisecond*700)
        for done != true {
            select {
                case f := <- dhtNode.fileResponse:
                    if f.Filename != "DONE" {
                        fileList = append(fileList, f)
                        waitResponse.Reset(time.Millisecond*700)
                    } else {
                        done = true
                    }
                case  <- waitResponse.C:
                    fmt.Println("^^^^^^^^^^^^^^^^^^^ GET TIMEOUT ^^^^^^^^^^^^^^")
                    done = true
            }
        }
    // decodes each file and puts them in new list to print
    for _, f := range fileList {   
        fmt.Println(f.Filename, f.Data)
        fileName, _ := b64.StdEncoding.DecodeString(f.Filename)
        fileData, _ := b64.StdEncoding.DecodeString(f.Data)
        eFile := &File{Filename: string(fileName), Data: string(fileData)}
        newFileList = append(newFileList, eFile)
    }

    fList := &FileList{Files: newFileList}

    // returns the JSON encoding and writes it to webpage
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
func (dhtNode *DHTNode) PostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    dec := json.NewDecoder(r.Body)
    file := File{}
    err := dec.Decode(&file)
    if err != nil {
        log.Fatal(err)
    }
    sFileName := b64.StdEncoding.EncodeToString([]byte(file.Filename))
    sFileData := b64.StdEncoding.EncodeToString([]byte(file.Data))

    dhtNode.responsibleForFile(sFileName, sFileData)
    //msg := createUploadMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName, sFileData)  
    //go func () { dhtNode.transport.send(msg)}()
}

/*****************************************
*** Handles updates of text in a file  ***
*** and updates same info on server    ***
*****************************************/
func (dhtNode *DHTNode) PutHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    dec := json.NewDecoder(r.Body)
    file := File{Filename: ps.ByName("KEY")}
    err := dec.Decode(&file)
    if err != nil {
        log.Fatal(err)
    }
    dhtNode.startUpdateFile(file.Filename, file.Data)
}

/*****************************************
*** Deletes a file from server         ***
*****************************************/
func (dhtNode *DHTNode) DeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    hash := generateNodeId(ps.ByName("KEY"))
    dhtNode.lookup(hash)
    waitResponse := time.NewTimer(time.Millisecond*1000)
    for {
        select {
            case n := <- dhtNode.fingerMemory:
                sFileName := b64.StdEncoding.EncodeToString([]byte(ps.ByName("KEY")))
                msg := createDeleteFileMsg("deleteFile", dhtNode.transport.bindAddress, n.ip, sFileName)  
                go func () { dhtNode.transport.send(msg)}()
                return
            case  <- waitResponse.C:
                fmt.Println("^^^^^^^^^^^^^^^^^^^DELETE TIMEOUT ^^^^^^^^^^^^^^")
                return
        }
    }
}
