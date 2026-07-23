package main

import (
	"fmt"
	"github.com/lukasnoeller/smugmug-backupper/internal/element"
	"github.com/lukasnoeller/smugmug-backupper/internal/user"
)

func main() {
	fmt.Println("User: ", user.GetUser())
	nodes := user.GetRootNode()
	cnr := element.GetChildNodes(nodes)
	element.GetElements(cnr)
	//user.GetEverything()

}
