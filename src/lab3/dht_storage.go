package dht

import (
	"io/ioutil"
	"os"
	/*"path/filepath"
	b64 "encoding/base64"*/
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

/* Adds the file to the nodes folder as the nodes nodeId.txt */
func (dhtNode *DHTNode) addFile(msg *Msg) {

	path := "dataFolder/" + dhtNode.nodeId + "/"

	/* If folder don't exist create one*/
	if !Exists(path) {
		os.MkdirAll(path, 0777)
	}

	path = "dataFolder/" + dhtNode.nodeId + "/" + dhtNode.nodeId + ".txt"
	createfile(path, msg.Data)	

	/* After creating file send it to successor to be replicated */
	NewMsg := createReplicateMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], msg.Data)
	go func () { dhtNode.transport.send(NewMsg)}()

}

/* Creates a new file named dhtNode.nodeId in folder "dataFolder" */
func createfile(path string, fileData []byte) {
	err := ioutil.WriteFile(path, fileData, 0777)

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
		
	path := "dataFolder/" + msg.Origin + "/" 

	/* If folder don't exist create one */
	if !Exists(path) {
		os.MkdirAll(path, 0777)

	}

	path = "dataFolder/" + msg.Origin + "/" + msg.Origin + ".txt"
	
	/* If file already exist, remove it and create it again 
	else, just create it */
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
		createfile(path, msg.Data)
  
	} else {
		createfile(path, msg.Data)
	}

}

/* If a node notice a predesessor is offline, it will take over 
responsebility for that file and move the file from it's backup 
folder to "/" */
func (dhtNode *DHTNode) takeResponsibility(predesessor_nodeId string) {
	path := "dataFolder/" + predesessor_nodeId + "/"

	files, err := ioutil.ReadDir(path)

		if err != nil {
		panic(err)
	}

	for _, f := range files {
			file, err := ioutil.ReadFile(path + f.Name())

			path = "dataFolder/" + dhtNode.nodeId + "/" + f.Name()
			createfile(path, file)

			NewMsg := createReplicateMsg(dhtNode.transport.bindAddress, dhtNode.successor[0], file)
			go func () { dhtNode.transport.send(NewMsg)}()
		}

}

/* When a predesessor node is back online, check for a file in 
the backupfolder, if any. Then remove the file in "dataFolder/" that matches */
func (dhtNode *DHTNode) dropResponsibility(nodeId string) {

	path := "dataFolder/" + nodeId + "/"

	/* If folder exist, remove every file in folder that's not named
	as the predesessors nodeId */
	if _, err := os.Stat(path); err == nil {
		files, error := ioutil.ReadDir(path)

		if error != nil {
			panic(error)
		}

		for _, f := range files {
			if ((path + f.Name()) != (path + nodeId + ".txt")) {
				os.Remove(path + f.Name())
			}
		}
		
	}

	path2 := "dataFolder/" + dhtNode.nodeId + "/"

	/* If folder exist, remove everything in folder that not named
	as the dhtnodes nodeId */
	if _, err := os.Stat(path); err == nil {
		files, error := ioutil.ReadDir(path)

		if error != nil {
			panic(error)
		}

		for _, f := range files {
			if ((path + f.Name()) != (path + dhtNode.nodeId + ".txt")) {
				os.Remove(path + f.Name())
			}
		}
	}
}

/* When a node is placed in the ring, check if sucessor has files 
belonging to the node and download it */
func (dhtNode *DHTNode) getFileBackup(msg *Msg) {
	path := "dataFolder/" + msg.Origin + "/"

	/* Checks if path exists */
	if Exists(path) {
		path = "dataFolder/" + msg.Origin + "/" + msg.Origin + ".txt"
	
		/* If file exist, send file back to node */
		if _, err := os.Stat(path); err == nil {
			file, error := ioutil.ReadFile(path)

			NewMsg := createUploadMsg(dhtNode.transport.bindAddress, msg.Origin, file)
			go func () { dhtNode.transport.send(NewMsg)}()

			dhtNode.dropResponsibility(msg.Origin)
		}
	}
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