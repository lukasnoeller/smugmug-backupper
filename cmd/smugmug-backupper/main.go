package main

import (
	"fmt"
	"github.com/lukasnoeller/smugmug-backupper/internal/user"
)

func main() {
	fmt.Println("User: ", user.GetUser())
	node := user.GetRootNode()
	user.GetChildNodes(node)
	//user.GetEverything()

}
