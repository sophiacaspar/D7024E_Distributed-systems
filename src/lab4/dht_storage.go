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

func (dhtNode *DHTNode) getFiles(msg *Msg) {

	fmt.Println("current node get files:", dhtNode.transport.bindAddress, "origin node is:", msg.Origin)
	if dhtNode.successor[0] != msg.Origin {
		dhtNode.goThroughDir(msg)
		fmt.Println(dhtNode.transport.bindAddress, "sends forward to", dhtNode.successor[0])
		next := createGetFilesMsg(msg.Origin, dhtNode.transport.bindAddress, dhtNode.successor[0])
		go dhtNode.transport.send(next)
	} else {
		d := false
		d = dhtNode.goThroughDir(msg)
		if d == true {
			fmt.Println("sent ack to", msg.Origin)
			m := createFileResponseMsg(dhtNode.transport.bindAddress, msg.Origin, "DONE", "DONE")
			fmt.Println("send this message")
			go dhtNode.transport.send(m)
			//m := createAckMsg("ack", dhtNode.transport.bindAddress, msg.Origin)
			//go dhtNode.transport.send(m)
		}
	}
}

func (dhtNode *DHTNode) goThroughDir(msg *Msg) bool {
	path := "dataFolder/" + dhtNode.nodeId + "/" 
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if !f.IsDir() {
			fileData, _ := ioutil.ReadFile(path + f.Name())
			sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
			sFileData := b64.StdEncoding.EncodeToString(fileData)
			m := createFileResponseMsg(dhtNode.transport.bindAddress, msg.Origin, sFileName, sFileData)
			go dhtNode.transport.send(m)
		} 
	}
	return true
}


/* Sends back the node that should be responsible for the file */
func (dhtNode *DHTNode) responsibleForFile(filename, data string) {

	hash := generateNodeId(filename)
	dhtNode.lookup(hash)
	
	waitResponse := time.NewTimer(time.Millisecond * 1000)

		for {
				select {
				case s := <-dhtNode.fingerMemory:
					sFileName := b64.StdEncoding.EncodeToString([]byte(filename))
					sFileData := b64.StdEncoding.EncodeToString([]byte(data))
					msg := createUploadMsg(dhtNode.transport.bindAddress, s.ip, sFileName, sFileData)
					go func () {dhtNode.transport.send(msg)}() 
				case <-waitResponse.C:
					return
				}
			}
}



/* Adds the file to the nodes folder as the name of the node 
and the filename*/
func (dhtNode *DHTNode) addFile(msg *Msg) {

	nodeId := generateNodeId(msg.Dst)
	succnodeId := generateNodeId(dhtNode.successor[0])
	path := "dataFolder/" + nodeId + "/" 

	/* If folder don't exist create one*/
	if !Exists(path) {
		os.MkdirAll(path, 0777)
	}

	FileName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	FileData, _ := b64.StdEncoding.DecodeString(msg.Data)
	path = "dataFolder/" + nodeId + "/" + string(FileName)

	createfile(path, string(FileData))	

	/* After creating file send it to successor to be replicated */
	fmt.Print("--------------------------------------------" + "\n")
	fmt.Print(msg.Dst + " replicates file to: " + dhtNode.successor[0] + " with NodeId: " + succnodeId + "\n")
	fmt.Print("--------------------------------------------" + "\n")
	
	sFileName := b64.StdEncoding.EncodeToString(FileName)
	sFileData := b64.StdEncoding.EncodeToString(FileData)

	NewMsg := createReplicateMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName, sFileData)
	go func () { dhtNode.transport.send(NewMsg)}()

}

func (dhtNode *DHTNode) newSuccessor() {

	nodeId := generateNodeId(dhtNode.successor[0])
	path := "dataFolder/" + dhtNode.nodeId + "/"

	fmt.Print("--------------------------------------------" + "\n")
	fmt.Print(dhtNode.transport.bindAddress + " new successor " + dhtNode.successor[0] + "\n")
	fmt.Print("--------------------------------------------" + "\n")

	
	/* If folder exist */
	if Exists(path) {
		files, err := ioutil.ReadDir(path)

		if err != nil {
			panic(err)
		}
	
	for _, f := range files {
		if !f.IsDir() && dhtNode.successor[0] != "" {
			file, _ := ioutil.ReadFile(path + f.Name())

			fmt.Print("--------------------------------------------" + "\n")
			fmt.Print(dhtNode.transport.bindAddress + " sends " + f.Name() + " to new successor " + dhtNode.successor[0] + " with NodeId: " + nodeId + "\n")
			fmt.Print("--------------------------------------------" + "\n")

			sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
			sFileData := b64.StdEncoding.EncodeToString(file)

			NewMsg := createReplicateMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName, sFileData)
			go func () { dhtNode.transport.send(NewMsg)}()
		}
	}

	}
}

func (dhtNode *DHTNode) newPredecessor(msg *Msg) {

	//oldPreNodeId := generateNodeId(msg.Src)
	NodeId := generateNodeId(msg.Origin)

	path := "dataFolder/" + dhtNode.nodeId + "/" + NodeId + "/"

	fmt.Print("--------------------------------------------" + "\n")
	fmt.Print(dhtNode.nodeId + " new predesessor " + NodeId +"\n")
	fmt.Print("--------------------------------------------" + "\n")

	/* If folder exist, send all data back to new predecessor*/
	if Exists(path) {
		files, err := ioutil.ReadDir(path)
		fmt.Println("in newPredecessor (files):",files)

		if err != nil {
			panic(err)
		}

		for _, f := range files {
			if !f.IsDir() {
			file, _ := ioutil.ReadFile(path + f.Name())
			fmt.Println("in newPredecessor (file):", file)

			fmt.Print("--------------------------------------------" + "\n")
			fmt.Print("Sends old file: " + f.Name() + " to new predecessor " + msg.Origin + " from path " + path + "\n")
			fmt.Print("--------------------------------------------" + "\n")

			sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
			sFileData := b64.StdEncoding.EncodeToString(file)

			NewMsg := createUploadMsg(dhtNode.transport.bindAddress, msg.Origin, sFileName, sFileData)
			go func () { dhtNode.transport.send(NewMsg)}()

			os.Remove(path + f.Name())

			} else {

				dir, _ := ioutil.ReadDir(path + f.Name())

				for _, f2 := range dir {

					file2, _ := ioutil.ReadFile(path + f2.Name())

					fmt.Print("--------------------------------------------" + "\n")
					fmt.Print("Sends back file : " + f2.Name() + " to new predecessor " + msg.Origin + " from path " + path + "\n")
					fmt.Print("--------------------------------------------" + "\n")
				
					sFileName := b64.StdEncoding.EncodeToString([]byte(f2.Name()))
					sFileData := b64.StdEncoding.EncodeToString(file2)

					NewMsg := createUploadMsg(dhtNode.transport.bindAddress, msg.Origin, sFileName, sFileData)
					go func () { dhtNode.transport.send(NewMsg)}()

					if f2.Name() != NodeId {
						os.RemoveAll(path + f2.Name())
					}

				}
			}
		}
	}
}

/* Creates a new file named dhtNode.nodeId in folder "path" */
func createfile(path string, fileData string) {
	err := ioutil.WriteFile(path, []byte(fileData), 0777)

	fmt.Print("--------------------------------------------" + "\n")
	fmt.Print("Creates file in path: " + path + " with data: " + string(fileData) + "\n")
	fmt.Print("--------------------------------------------" + "\n")
	
	check(err)

}

/*  */
func check(err error) {
	if err != nil {
		panic(err)
	}
}

/* Replicates the file to the nodes succesor */
func (dhtNode *DHTNode) replicate(msg *Msg) {

	nodeId := generateNodeId(msg.Origin)
		
	path := "dataFolder/" + dhtNode.nodeId + "/" + nodeId + "/" 

	/* If folder exist */
	if !Exists(path) {
		os.MkdirAll(path, 0777)
	}

	FileName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	FileData, _ := b64.StdEncoding.DecodeString(msg.Data)
	fmt.Println(FileData, msg.Data)
	path2 := "dataFolder/" + dhtNode.nodeId + "/" + nodeId + "/" + string(FileName)
	
	/* If file already exist, remove it and create it again 
	else, just create it */
	if _, err := os.Stat(path2); err == nil {
		os.Remove(path2)
		createfile(path2, string(FileData))
  
	} else {
		createfile(path2, string(FileData))
	}

}

/* If a node notice a predesessor is offline, it will take over 
responsebility for that file and move the file from it's backup 
folder to "/" */
func (dhtNode *DHTNode) takeResponsibility(predesessor_nodeId string) {

	nodeId := generateNodeId(predesessor_nodeId)
	path := "dataFolder/" + dhtNode.nodeId + "/" + nodeId + "/"

	/* If folder exist*/
	if Exists(path) {
	files, err := ioutil.ReadDir(path)

		if err != nil {
		panic(err)
		}
	
	for _, f := range files {
		file, _ := ioutil.ReadFile(path + f.Name())

		path2 := "dataFolder/" + dhtNode.nodeId + "/" + f.Name()
		createfile(path2, string(file))

		fmt.Print("--------------------------------------------" + "\n")
		fmt.Print("Moves file " + f.Name() + " in " + path + " to " + path2 + "\n" )
		fmt.Print("--------------------------------------------" + "\n")

		sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
		sFileData := b64.StdEncoding.EncodeToString(file)

		NewMsg := createReplicateMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName, sFileData)
		go func () { dhtNode.transport.send(NewMsg)}()
	}

	} else {

	fmt.Print("--------------------------------------------" + "\n")
	fmt.Print("No files in folder " + path + "\n" )
	fmt.Print("--------------------------------------------" + "\n")

	}
}

/* When a predesessor node is back online, check for a file in 
the backupfolder, if any. Then remove the file in "dataFolder/" that matches */
func (dhtNode *DHTNode) dropResponsibility(preNodeId string, fileName string) {

	nodeId := generateNodeId(preNodeId)	

	path := "dataFolder/" + dhtNode.nodeId + "/" + nodeId + "/"
	path2 := "dataFolder/" + dhtNode.nodeId + "/" 

	/* If folder exist, remove every file in folder that's not named
	as the predesessors nodeId */

	fmt.Print("--------------------------------------------" + "\n")
	fmt.Print(dhtNode.transport.bindAddress + " drops responsibility for file: " + fileName + " in " + path2 + "\n")
	fmt.Print("--------------------------------------------" + "\n")

	if _, err := os.Stat(path); err == nil {
		files, error := ioutil.ReadDir(path)

		if error != nil {
			panic(error)
		}

		for _, f := range files {
			if f.Name() == fileName {
				os.Remove(path + f.Name())
			}
		}
	}

	if _, err := os.Stat(path2); err == nil {
		files, error := ioutil.ReadDir(path2)

		if error != nil {
			panic(error)
		}

		for _, f := range files {
			if f.Name() == fileName {
				sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
				os.Remove(path2 + f.Name())
				
				NewMsg := createDeleteFileMsg("deleteFileSucc", dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName)
				go func () { dhtNode.transport.send(NewMsg)}()
			}
		}
	}
}

/* When a node is placed in the ring, check if successor has files 
belonging to the node and download it */
func (dhtNode *DHTNode) getFileBackup(msg *Msg) {
	
	nodeId := generateNodeId(msg.Origin)	
	path := "dataFolder/" + dhtNode.nodeId + "/" + nodeId + "/"
	
	/* Checks if path exists */
	if Exists(path) {
		files, err := ioutil.ReadDir(path)

		if err != nil {
			panic(err)
		}

		for _, f := range files {
			if !f.IsDir() {
				file, _ := ioutil.ReadFile(path + f.Name())

				sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))
				sFileData := b64.StdEncoding.EncodeToString(file)

				fmt.Print("--------------------------------------------" + "\n")
				fmt.Print(dhtNode.transport.bindAddress + " sends back file: " + f.Name() + " to " + msg.Origin + "\n")
				fmt.Print("--------------------------------------------" + "\n")

				NewMsg := createUploadMsg(dhtNode.transport.bindAddress, msg.Origin, sFileName, sFileData)
				go func () { dhtNode.transport.send(NewMsg)}()

				dhtNode.dropResponsibility(msg.Origin, f.Name())
			} 
		}	
	} else {
		fmt.Print("--------------------------------------------" + "\n")
		fmt.Print("No files in folder " + path + "\n" )
		fmt.Print("--------------------------------------------" + "\n")
	}
}

/* Called from webpage, deletes file */
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

				fmt.Print("--------------------------------------------" + "\n")
				fmt.Print("Deleting file: " + path + f.Name() + "\n")
				fmt.Print("--------------------------------------------" + "\n")

				os.Remove(path + f.Name())

				sFileName := b64.StdEncoding.EncodeToString([]byte(f.Name()))

				NewMsg := createDeleteFileMsg("deleteFileSucc", dhtNode.transport.bindAddress, dhtNode.successor[0], sFileName)
				go func () { dhtNode.transport.send(NewMsg)}()
			}
		}
	}
}

func (dhtNode *DHTNode) removefolder(path string, dirName string) {
	
	fmt.Print("--------------------------------------------" + "\n")
	fmt.Print(dhtNode.nodeId + " removes folder: " + dirName + " in path: " + path + "\n")
	fmt.Print("--------------------------------------------" + "\n")

	os.RemoveAll(path + dirName)
}

func (dhtNode *DHTNode) deleteFileSucc(msg *Msg) {
	nodeId := generateNodeId(msg.Origin)
	path := "dataFolder/" + dhtNode.nodeId + "/" + nodeId + "/"

	FileName, _ := b64.StdEncoding.DecodeString(msg.FileName)

	if Exists(path) {
		files, err := ioutil.ReadDir(path)

		if err != nil {
			panic(err)
		}

		for _, f := range files {
			if f.Name() == string(FileName) {

				fmt.Print("--------------------------------------------" + "\n")
				fmt.Print(dhtNode.nodeId + " deleted file " + f.Name() + " in folder: " + nodeId + "\n")
				fmt.Print("--------------------------------------------" + "\n")

				os.Remove(path + f.Name())
			}
		}
		newPath := "dataFolder/" + dhtNode.nodeId + "/"
		dhtNode.removefolder(newPath, nodeId)
	}
}

func (dhtNode *DHTNode) startUpdateFile(filename, data string) {
    hash := generateNodeId(filename)
    dhtNode.lookup(hash)
    waitResponse := time.NewTimer(time.Millisecond*1000)
        for {
            select {
                case n := <- dhtNode.fingerMemory:
                    sFileName := b64.StdEncoding.EncodeToString([]byte(filename))
					sFileData := b64.StdEncoding.EncodeToString([]byte(data))
                    msg := createUpdateFileMsg("updateFile", dhtNode.transport.bindAddress, n.ip, sFileName, sFileData)  
                    go func () { dhtNode.transport.send(msg)}()
                    return
                case  <- waitResponse.C:
                    fmt.Println("^^^^^^^^^^^^^^^^^^^ UPDATE TIMEOUT ^^^^^^^^^^^^^^")
                    return
            }
        }
}

func (dhtNode *DHTNode) updateFile(msg *Msg) {
	FileName, _ := b64.StdEncoding.DecodeString(msg.FileName)
	Data, _ := b64.StdEncoding.DecodeString(msg.Data)
	path := "dataFolder/" + dhtNode.nodeId + "/" + string(FileName)
	ioutil.WriteFile(path, []byte(Data), 0777)
	m := createReplicateMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], msg.FileName, msg.Data)
	go func () { dhtNode.transport.send(m)}()
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