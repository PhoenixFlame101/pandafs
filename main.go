package main

import (
	"os"
	"pandafs/master"
	"pandafs/mule"
)

func main() {
	nodeType := os.Args[1]
	switch nodeType {
	case "master":
		master.GetMasterNode().Start()
	case "mule":
		mule.GetMule().Start()
	default:
		panic("invalid node type")
	}
}
