package main

import (
	"./hashring"
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		fmt.Printf("Usage: ./test_go_hashring nodes_file keys_file\n")
		return
	}
	nodesFile, _ := os.Open(flag.Arg(0))
	keysFile, _ := os.Open(flag.Arg(1))

	nodesReader := bufio.NewReader(nodesFile)
	keysReader := bufio.NewReader(keysFile)

	nodes := []string{}
	keys := []string{}

	for {
		text, err := nodesReader.ReadString('\n')
		if err != nil {
			break
		}
		nodes = append(nodes, text[:len(text)-1])
	}

	for {
		text, err := keysReader.ReadString('\n')
		if err != nil {
			break
		}
		keys = append(keys, text[:len(text)-1])
	}

	ring := hashring.New(nodes)
	for i := range keys {
		fmt.Printf("%s -> %s\n", keys[i], ring.GetNode(keys[i]))
	}
}
