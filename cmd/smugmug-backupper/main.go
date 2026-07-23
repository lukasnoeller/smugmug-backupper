package main

import (
	"fmt"
	"path"

	"github.com/lukasnoeller/smugmug-backupper/internal/user"
)

func main() {
	authUser := user.GetUser()
	fmt.Println("authUser: ", authUser)
	node := path.Base(authUser.Uris.Node.Uri)
	fmt.Println("node", node)

}
