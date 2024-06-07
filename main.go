package main

import (
	"fmt"
	"hash/fnv"
)

// MOD value for mod
// var MOD int = 6

// func sumKey(str string) int {
// 	byteArr := []byte(str)

// 	sum := 0
// 	for _, value := range byteArr {
// 		sum += int(value)
// 	}

// 	return sum
// }

func contains(slice []string, key string) bool {
	for _, value := range slice {
		if value == key {
			return true
		}
	}

	return false
}

// StorageNode struct to hold info regarding the nodes
type StorageNode struct {
	name string
	host string
}

var servers = map[string][]string{}

func (node *StorageNode) putFile(file string) {
	servers[node.name] = append(servers[node.name], file)
}

func (node *StorageNode) fetchFile(file string) bool {
	return contains(servers[node.name], file)
}

var storageNodes = []StorageNode{
	{
		name: "A",
		host: "10.131.213.12",
	},
	{
		name: "B",
		host: "10.131.217.11",
	},
	{
		name: "C",
		host: "10.131.142.46",
	},
	{
		name: "D",
		host: "10.131.189.18",
	},
	{
		name: "E",
		host: "10.131.210.10",
	},
	{
		name: "F",
		host: "10.131.231.32",
	},
}

func hashFn(key string, totalServers int) int {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	hashValue := hash.Sum32()

	return int(hashValue) % totalServers
}

func uploadFn(path string) {
	index := hashFn(path, len(storageNodes))
	node := storageNodes[index]

	node.putFile(path)
	fmt.Printf("File has been uploaded to server %s\n\n", node.name)
}

func fetchFn(path string) {
	index := hashFn(path, len(storageNodes))
	node := storageNodes[index]

	if ok := node.fetchFile(path); !ok {
		fmt.Println("Cannot find the file!")
		return
	}

	fmt.Printf("Found the file %s on node %s\n\n", path, node.name)
}

func main() {
	for true {
		fmt.Println("Enter choice: ")
		fmt.Printf("Put file: \t1\nFetch file:\t2\nEnd program:\t3\n")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Print("Enter file path: ")
			var path string
			fmt.Scanln(&path)

			uploadFn(path)
		case 2:
			fmt.Print("Enter file path: ")
			var path string
			fmt.Scanln(&path)

			fetchFn(path)
		case 3:
			fmt.Print("Ending the program!\n")
			return
		}
	}
}
