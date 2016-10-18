package dht

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
At start, check if "dataFolder" exist.
If not creates a new one with permission 0777
*/
func (dhtNode *DHTNode) makeFolder() {
	path := "dataFolder/" + dhtNode.nodeId

	if !Exists(path) {
		MkdirAll(path, 0777)
	}
}

/* Adds dataFile.txt to the nodes folder */
func (dhtNode *DHTNode) addFile(msg *MSG) {

	fileData := msg.data

	path := "dataFolder/" + dhtNode.nodeId + "/"

	/* If folder don't exist create one*/
	if !Exists(path) {
		MkdirAll(path, 0777)
	}

	path = "dataFolder/" + dhtNode.nodeId + "/" + dhtNode.nodeId + ".txt"
	createfile(path, fileData)	

	msg := createReplicationMsg("replication", dhtNode.transport.bindAddress, dhtNode.succesor[0], data)
	go func () { node1.transport.send(msg)}() 

}

/* Creates a new file named "dataFile" in folder "dataFolder" */
func createfile(path string, fileData string) {
	data := []byte(fileData)
	err := ioutil.WriteFile(path, data, 0777)

	check(err)

}

/*  */
func check(err error) {
	if err != nil {
		panic(err)
	}
}

/* Replicates the file to the nodes succesor */
func (dhtNode *DHTNode) replicate(msg *MSG) {
	if dhtNode.succesor[0] != "" {
		
		path := "dataFolder/" + msg.Origin + "/" msg.Origin + ".txt"

		// all files in dir
		files, _ := ioutil.ReadDir(path)







	}
}

/* If a node notice a predesessor is offline, it will take over 
responsebility for that file and move the file from it's backup 
folder to "/" */
func (dhtNode *DHTNode) takeResponsibility() {
	path := "dataFolder/" + 
}

/* When a predesessor node is back online, check the file 
in backupfolder, if any. Then remove the file in "/" that matches */
func (dhtNode *DHTNode) dropResponsibility() {
	path := "dataFolder/" + 
}

/* When a node is placed in the ring, check if sucessor has files 
belonging to the node and download it */
func getFileBackup() {
	path := "dataFolder/" + dhtNode.nodeId + "/"

	/* Checks if path exists */
	if Exists(path) {
		
		// all files in dir
		files, _ := ioutil.ReadDir(path)



	}

}
