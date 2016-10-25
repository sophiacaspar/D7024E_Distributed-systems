package dht

import (
	"io/ioutil"
	"os"
	"fmt"
	"time"
	/*"path/filepath"*/
	b64 "encoding/base64"
)

/*
At start, check if "dataFolder" exist.
If not creates a new one with permission 0777
*/
func (dhtNode *DHTNode) makeFolder() {
	path := "dataFolder/" + dhtNode.nodeId

	if !Exists(path) {
		os.MkdirAll(path, 0777)
	}
}

/* createUploadMsg Adds the file to the nodes folder as the name of the node 
and the filename */
func (dhtNode *DHTNode) addFile(msg *Msg) {
	
	path := "dataFolder/" + dhtNode.nodeId + "/" 

	FileName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	FileData, _ := b64.StdEncoding.DecodeString(msg.Data)

	hash := generateNodeId(string(FileName))
	
	if dhtNode.responsible(hash) != true{						
		sFileName := b64.StdEncoding.EncodeToString([]byte(FileName))
		sFileData := b64.StdEncoding.EncodeToString([]byte(FileData))
		msg := createUploadMsg(dhtNode.transport.bindAddress, dhtNode.predecessor[0], sFileName, sFileData)
		go func () {dhtNode.transport.send(msg)}() 

	} else {
	/* If folder don't exist create one*/
	if !Exists(path) {
		os.MkdirAll(path, 0777)
	}

	path = "dataFolder/" + dhtNode.nodeId + "/" + string(FileName)
	createfile(path, string(FileData))	

	/* After creating file send it to successor to be replicated */
	fmt.Println(msg.Dst + " replicates file to: " + dhtNode.successor[0] + "\n")
	
	sFileName := b64.StdEncoding.EncodeToString(FileName)
	sFileData := b64.StdEncoding.EncodeToString(FileData)

	NewMsg := createReplicateMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName, sFileData)
	go func () { dhtNode.transport.send(NewMsg)}()
	}
}
	
/* createCheckSuccDataMsg */
func (dhtNode *DHTNode) getSuccData(msg *Msg) {

	path := "dataFolder/" + dhtNode.nodeId + "/" 

	if Exists(path) {
		files, err := ioutil.ReadDir(path)

		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if f.Name() == ".DS_Store" {
				os.Remove(path + f.Name())
			} else {
				empty, err := IsEmptyDir(path)
					if err != nil {
						panic(err)
					}
				if !empty {
					if !f.IsDir(){
						file, _ := ioutil.ReadFile(path + f.Name())

						hash := generateNodeId(f.Name())
						
						if !(between([]byte(msg.LightNode[1]),[]byte(dhtNode.nodeId), []byte(hash))){
						
							sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
							sFileData := b64.StdEncoding.EncodeToString([]byte(file))
							msg := createUploadMsg(dhtNode.transport.bindAddress, msg.LightNode[0], sFileName, sFileData)
							go func () {dhtNode.transport.send(msg)}() 

							os.Remove(path + f.Name())
						} else {
							fmt.Println("Node " + msg.Dst + " with nodeId " + dhtNode.nodeId + " responsebility for " + f.Name() + "\n")
						}
					} 
				}
			}
			emptyer, error := IsEmptyDir(path)
			if error != nil {
				panic(err)
			}
			if emptyer {
				sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
				NewMsg := createDeleteFileMsg("deleteFileSucc", dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName)
				go func () { dhtNode.transport.send(NewMsg)}()
			} else {
				fmt.Println(path + " is not empty" + "\n")
			}	
		} 	
	}
}

func (dhtNode *DHTNode) checkfolder(nodeId string) {
path := "dataFolder/" + dhtNode.nodeId + "/"  + nodeId + "/" 
path2 := "dataFolder/" + dhtNode.nodeId + "/"
	
if Exists(path) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
		for _, f := range dir {
			if f.Name() == ".DS_Store" {
				os.Remove(path + f.Name())
			} else {
				files, _ := ioutil.ReadDir(path2)
				for _, f2 := range files {
					for _, f3 := range dir {
						if !f.IsDir() && f3.Name() == f2.Name(){
						fmt.Println("Remove file: "+ f3.Name() +" in "+ path2 + "\n" )
						os.Remove(path2 + f2.Name())
					}

					}
				}
			}
		}
	}
}

func (dhtNode *DHTNode) checkIfReplicate() {
	path := "dataFolder/" + dhtNode.nodeId + "/" 
	
	if Exists(path) {
		files, err := ioutil.ReadDir(path)

		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if f.Name() == ".DS_Store" {
				os.Remove(path + f.Name())
			} else {
				if !f.IsDir(){
					file, _ := ioutil.ReadFile(path + f.Name())

					sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
					sFileData := b64.StdEncoding.EncodeToString([]byte(file))

					NewMsg := createReplicateMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName, sFileData)
					go func () { dhtNode.transport.send(NewMsg)}()
				}
			}
		}
	}
}
	
/* Creates a new file named dhtNode.nodeId in folder "path" */
func createfile(path string, fileData string) {
	err := ioutil.WriteFile(path, []byte(fileData), 0777)
	check(err)
}

/*  */
func check(err error) {
	if err != nil {
		panic(err)
	}
}

/* createReplicateMsg Replicates the file to the nodes succesor */
func (dhtNode *DHTNode) replicate(msg *Msg) {

	nodeId := generateNodeId(msg.Origin)
		
	path := "dataFolder/" + dhtNode.nodeId + "/" + nodeId + "/" 

	/* If folder exist */
	if !Exists(path) {
		os.MkdirAll(path, 0777)
	}

	FileName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	FileData, _ := b64.StdEncoding.DecodeString(msg.Data)

	path2 := "dataFolder/" + dhtNode.nodeId + "/" + nodeId + "/" + string(FileName)
	
	/* If file already exist, remove it and create it again 
	else, just create it */
	if _, err := os.Stat(path2); err == nil {
		os.Remove(path2)
		createfile(path2, string(FileData))
  
	} else {
		createfile(path2, string(FileData))
	}
	dhtNode.checkfolder(nodeId)
}

/* If a node notice a predesessor is offline, it will take over 
responsebility for that file and move the file from it's backup 
folder to "/" */
func (dhtNode *DHTNode) takeResponsibility() {

	path := "dataFolder/" + dhtNode.nodeId + "/" + dhtNode.predecessor[1] + "/"
	path3 := "dataFolder/" + dhtNode.nodeId + "/" 

	/* If folder exist*/
	if Exists(path) {
		fmt.Println("Takes responsebility of file in: " + path + "\n")
		files, err := ioutil.ReadDir(path)

		if err != nil {
		panic(err)
		}
	
		for _, f := range files {
			file, _ := ioutil.ReadFile(path + f.Name())

			path2 := "dataFolder/" + dhtNode.nodeId + "/" + f.Name()
			createfile(path2, string(file))

			fmt.Println("Moves file " + f.Name() + " in " + path + " to " + path2 + "\n" )

			sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
			sFileData := b64.StdEncoding.EncodeToString(file)

			/* After creating file send it to successor to be replicated */
			fmt.Println(dhtNode.nodeId + " replicates file to: " + dhtNode.successor[1] + "\n")

			NewMsg := createReplicateMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName, sFileData)
			go func () { dhtNode.transport.send(NewMsg)}()

			os.Remove(path + f.Name())
		}

		empty, err := IsEmptyDir(path)
		if err != nil {
			panic(err)
		}
		if empty {
			dhtNode.removefolder(path3, dhtNode.predecessor[1])
		} else {
			fmt.Printf(path + " is not empty" + "\n")
		}
	} else {
		fmt.Println(path + " don't exist" + "\n")
	} 
}

func (dhtNode *DHTNode) deleteFile(msg *Msg) {
	path := "dataFolder/" + dhtNode.nodeId + "/" 
	FileName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	/* If folder exist, remove every file in folder that's not named
	as the predesessors nodeId */
	if _, err := os.Stat(path); err == nil {
		files, error := ioutil.ReadDir(path)
		if error != nil {
			panic(error)
		}
		for _, f := range files {
			if f.Name() == string(FileName) {
				fmt.Println("Deleting file: " + path + f.Name() + "\n")
				os.Remove(path + f.Name())
				sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
				NewMsg := createDeleteFileMsg("deleteFileSucc", dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName)
				go func () { dhtNode.transport.send(NewMsg)}()
			}
		}
	}
}

func (dhtNode *DHTNode) removefolder(path string, dirName string) {

	fmt.Println(dhtNode.nodeId + " removes folder: " + dirName + " in path: " + path + "\n")

	os.RemoveAll(path + dirName)
}

/* createSuccFileDeleteMsg */
func (dhtNode *DHTNode) deleteFileSucc(msg *Msg) {

	FileName, _ := b64.StdEncoding.DecodeString(msg.FileName)

	path := "dataFolder/" + dhtNode.nodeId + "/" + dhtNode.predecessor[1] + "/" 

	if Exists(path) {
		os.Remove(path + string(FileName))
	}
	empty, err := IsEmptyDir(path)
	if err != nil {
		panic(err)
	}
	if empty {
		path2 := "dataFolder/" + dhtNode.nodeId + "/"
		dhtNode.removefolder(path2, dhtNode.predecessor[1])
	} else {
		fmt.Printf(path + " is not empty" + "\n")
	}	
}

/* If backupfolderNode dies and gets back alive. It's successor deletes the 
folder named as it pre predecessor */
func (dhtNode *DHTNode) deleteBackupSucc(msg *Msg) {

	path := "dataFolder/" + dhtNode.nodeId + "/" 

	if Exists(path) {
		files, err := ioutil.ReadDir(path)

		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if f.Name() == ".DS_Store" {
				os.Remove(path + f.Name())
			} else if f.IsDir() && f.Name() == msg.FileName {
				dhtNode.removefolder(path, f.Name())
			}
		}
	}	
}

func IsEmptyDir(name string) (bool, error) {
	entries, err := ioutil.ReadDir(name)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}
	
// Do the file exist?
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

/* Sends back the node that should be responsible for the file */
func (dhtNode *DHTNode) responsibleForFile(filename, data string) {
	FileName, _ := b64.StdEncoding.DecodeString(filename)
	
	hash := generateNodeId(string(FileName))

	fmt.Println("Hash: " + string(FileName) + " to: " + hash + "\n")

	dhtNode.lookup(hash)
	
	waitResponse := time.NewTimer(time.Millisecond * 2000)

		for {
				select {
				case s := <-dhtNode.fingerMemory:
					fmt.Println(" <<<<<<<<<<<<<<<<<<<<<< First time uploading file to: " + s.ip + " <<<<<<<<<<<<<<<<<<<<<< ")

					msg := createUploadMsg(dhtNode.transport.bindAddress, s.ip, filename, data)
					go func () {dhtNode.transport.send(msg)}() 
					return
					
				case <-waitResponse.C:
					return
				}
			}
}