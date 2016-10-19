package dht

import (
	"fmt"
	"html/template"
	//"io/ioutil"
	"net/http"
	"github.com/httprouter-master" //https://github.com/julienschmidt/httprouter
	"bufio"
	"os"
	"log"
)

type Page struct {
	Filename	string
	Address		string
	//Body 		[]byte

	//Address		string
}

/*
func loadPage(filename string) (*Page, error) {
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Filename: filename, Body: body}, nil
}
*/
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
    f := &Page{Filename: "testfilename", Address: "localhost:1115"}

	htmlStr := test()
	t, _ := template.New("webpage").Parse(htmlStr)
    //t, _ = t.Parse(htmlStr)
    //t, _ = t.Parse("<head><title>skynet</title><script src='webserver.js'></script></head><h1>The Grand Repository in the Sky</h1>testing {{.Filename}}\n<table style='width:80%'>\n<tr><th>Select new file: <input type='file' id='myFile'> <button onclick=''>Upload</button></th> <th><form action='/storage' method='POST'><textarea name='message' rows='10' cols='30'></textarea><br><input type='submit' value='Update'></form></th></tr><tr> <th> {{.Filename}} <button onclick=''>Open</button> <button onclick=''>Delete</button></th></tr></table>")

    //t, _ := template.ParseFiles("webpage.html")
/*
    f, err := loadPage("webpage.html")
    if err != nil {
        fmt.Println(err)
    }
   */ 

    t.Execute(w, f)
    //fmt.Fprintf(w, string(f))
}

func test() string{
	htmlStr := ""
    file, err := os.Open("webpage.html")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        //fmt.Println(scanner.Text())
        htmlStr = htmlStr + scanner.Text() + "\n"
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return htmlStr
}

func HelloHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s %s!\n", ps.ByName("name"), ps.ByName("lastname"))
}

func (dhtNode *DHTNode) startWebserver() {
    router := httprouter.New()
    router.GET("/", IndexHandler)
    router.GET("/hello/:name/:lastname", HelloHandler)
    router.POST("/storage", postData)

    log.Fatal(http.ListenAndServe(dhtNode.contact.ip+":"+dhtNode.contact.port, router))
}


func postData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("asdfghjklölkgfdsasdfghjklöljgfdsasdfghjklögdsaasdfgh", r.Body)
	
}