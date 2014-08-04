package hashring

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestHashring(t *testing.T) {
	nodesFile, _ := os.Open("nodes")
	keysFile, _ := os.Open("keys")
	
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


	w, err := os.Create("go_results")
	if err != nil {
		t.Error("failed to create write file")
		return
	}

	ring := New(nodes)
	for i := range keys {
		fmt.Fprintf(w, "%s -> %s\n", keys[i], ring.GetNode(keys[i]))
	}

}

