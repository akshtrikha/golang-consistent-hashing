package main

import (
	"fmt"
	"hash/fnv"
	"slices"
	"sort"
)

// func contains(arr []string, file string) bool {
// 	for _, value := range arr {
// 		if value == file {
// 			return true
// 		}
// 	}

// 	return false
// }

func hashFn(key string, totalServers int) int {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	hashValue := hash.Sum32()

	return int(hashValue) % totalServers
}

// ServerHostNames map to store the hostname of the servers against server names
var ServerHostNames = make(map[string]string)

// ServerNode struct to hold info regarding the nodes
type ServerNode struct {
	name    string
	host    string
	// content []string
}

// ConsistentHash represents an array based implementation of consistent hashing algorithm
type ConsistentHash struct {
	keys         []int
	nodes        []ServerNode
	totalServers int
}

// Init function to initialize the Consistent Hash
func (c *ConsistentHash) Init(totalServers int) {
	// define the capacity of your consistent hash
	c.totalServers = totalServers
}

// AddNode function to add a node to the consistent Hash
// Returns the key from the hash space where it was added
func (c *ConsistentHash) AddNode(node ServerNode) (int, error) {

	// handling case when the hash space is full
	if len(c.keys) == c.totalServers {
		return -1, fmt.Errorf("Hash Space is full. Cannot insert node %+v", node)
	}

	// getting the node/server key using the hash function
	// hash the node hostname
	key := hashFn(node.host, c.totalServers)

	// find the index where the key should be inserted in the keys array
	// this will be the index where the Storage Node will be added im
	// the nodes array
	index := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] > key })

	// if we have already seen the key i.e., node already is present
	// for the same key, we raise Collision Error
	if index > 0 && c.keys[index-1] == key {
		return -1, fmt.Errorf("Collision Error for key %d", key)
	}

	// insert the node_id and the key at the same 'index' location
	// this insertion will keep nodes and keys sorted w.r.t. keys
	c.nodes = slices.Insert(c.nodes, index, node)
	c.keys = slices.Insert(c.keys, index, key)

	return key, nil
}

// RemoveNode removes the node and returns the key from the hash space
// on which the node was placed
func (c *ConsistentHash) RemoveNode(node ServerNode) (int, error) {

	// handling the case when the hash space is empty
	if len(c.keys) == 0 {
		return -1, fmt.Errorf("Hash Spcae is empty. Did not find node %+v", node)
	}

	// get the key from the hash function
	// hash the node hostname
	key := hashFn(node.host, c.totalServers)

	// find the index where the key would reside in the keys
	index := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] == key })

	// if key does not exist in the array,
	// raise error
	if index >= len(c.keys) || c.keys[index] != key {
		return -1, fmt.Errorf("Node %+v does not exist", node)
	}

	c.keys = slices.Delete(c.keys, index, index+1)
	c.nodes = slices.Delete(c.nodes, index, index+1)

	return key, nil
}

// Assign - Given an item, the function returns the node it is associated with
func (c *ConsistentHash) Assign(item string) (*ServerNode, error) {
	// get the key for the item
	key := hashFn(item, c.totalServers)

	// we find the first node to the right of this key
	// if bisect_right returns index which is out of bounds
	// then we circle back to the first in the array in a circular fashion
	index := (sort.Search(len(c.keys), func(i int) bool { return c.keys[i] > key })) % len(c.keys)

	if index >= len(c.keys) {
		return &ServerNode{}, fmt.Errorf("Cannot find a server for %s", item)
	}

	return &c.nodes[index], nil
}

func main() {

	c := ConsistentHash{}
	c.Init(50)

	// for true {
	// 	fmt.Printf("Current no. of server: %v\n", len(c.keys))
	// 	fmt.Println("Enter choice: ")
	// 	fmt.Printf("Put file: \t\t1\nFetch file:\t\t2\nServer Options:\t\t3\nEnd the program: \t4\n")

	// 	var choice int
	// 	fmt.Scanln(&choice)

	// 	switch choice {
	// 	case 1:
	// 		if len(c.keys) == 0 {
	// 			fmt.Println("Cannot put/fetch files until servers are added!")
	// 		} else {
	// 			fmt.Print("Enter file path: ")
	// 			var path string
	// 			fmt.Scanln(&path)

	// 			server, err := c.Assign(path)
	// 			if err != nil {
	// 				fmt.Printf("Ran into an error: %v\n", err)
	// 			} else {
	// 				(*server).content = append(server.content, path)
	// 				fmt.Printf("Stored file at the server: %+v\n", *server)
	// 			}
	// 		}
	// 	case 2:
	// 		if len(c.keys) == 0 {
	// 			fmt.Println("Cannot put/fetch files until servers are added!")
	// 		} else {
	// 			fmt.Print("Enter file path: ")
	// 			var path string
	// 			fmt.Scanln(&path)

	// 			server, err := c.Assign(path)
	// 			if err != nil {
	// 				fmt.Printf("Ran into an error: %v\n", err)
	// 			} else {
	// 				if contains((*server).content, path) {
	// 					fmt.Printf("File is stored at the server: %+v\n", *server)
	// 				} else {
	// 					fmt.Println("Cannot find the file")
	// 				}
	// 			}
	// 		}
	// 	case 3:
	// 		fmt.Println("Enter the choice:")
	// 		fmt.Printf("Add a server\t1\nRemove a server\t2\n")

	// 		var serverChoice int
	// 		fmt.Scanln(&serverChoice)

	// 		switch serverChoice {
	// 		case 1:
	// 			var server ServerNode

	// 			fmt.Println("Enter the name of the server:")
	// 			fmt.Scanln(&server.name)

	// 			fmt.Println("Enter the hostname of the server:")
	// 			fmt.Scanln(&server.host)

	// 			fmt.Printf("Server Details: %+v\n", server)

	// 			ServerHostNames[server.name] = server.host

	// 			key, err := c.AddNode(server)
	// 			if err != nil {
	// 				fmt.Printf("Ran into an error: %v\n", err)
	// 			} else {
	// 				fmt.Printf("Successfully added a server at %v\n", key)
	// 			}
	// 		case 2:
	// 			var serverName string
	// 			fmt.Println("Enter the name of the server to remove:")
	// 			fmt.Scanln(&serverName)

	// 			serverHostName := ServerHostNames[serverName]
	// 			server := ServerNode{serverName, serverHostName, []string{}}

	// 			key, err := c.RemoveNode(server)
	// 			if err != nil {
	// 				fmt.Printf("Ran into an error: %v\n", err)
	// 			} else {
	// 				fmt.Printf("Successfully removed the server %+v at %v\n", server, key)
	// 			}
	// 		}
	// 	case 4:
	// 		fmt.Print("Ending the program!\n")
	// 		return
	// 	case 5:
	// 		fmt.Println("Here are the servers:")
	// 		for _, server := range c.nodes {
	// 			fmt.Printf("%+v\n", server)
	// 		}
	// 	}
	// }
}
