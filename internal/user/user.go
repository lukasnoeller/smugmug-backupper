package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lukasnoeller/smugmug-backupper/internal/auth"
	"path"
)

type User struct {
	Name     string `json:"Name"`
	NickName string `json:"NickName"`
	Uris     struct {
		Node   Node `json:"Node"`
		Folder struct {
			Uri string `json:"Uri"`
		} `json:"Folder"`
	} `json:"Uris"`
}
type Node struct {
	Uri string `json:"Uri"`
}
type AuthUserResponse struct {
	Response struct {
		User User `json:"User"`
	} `json:"Response"`
}

func GetUser() User {

	body, err := auth.AuthenticatedRequest("GET", "https://api.smugmug.com/api/v2!authuser")
	var authUserResponse AuthUserResponse

	err = auth.Parse(body, &authUserResponse)

	if err != nil {
		fmt.Println("error ocurred during unmashalling of user ", err)
	}
	authUser := authUserResponse.Response.User
	return authUser
}
func GetRootNode() string {
	return path.Base(GetUser().Uris.Node.Uri)
}
func GetChildNodes(nodeUri string) {
	nodeRequestUrl := fmt.Sprintf("https://api.smugmug.com/api/v2/node/%s!children", nodeUri)
	body, err := auth.AuthenticatedRequest("GET", nodeRequestUrl)
	if err != nil {
		fmt.Println("error ocurred during request: ", err)
	}
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		fmt.Println("Invalid JSON:", err)
		return
	}
	fmt.Println(prettyJSON.String())

}
func GetEverything() {
	body, err := auth.AuthenticatedRequest("GET", "https://api.smugmug.com/api/v2!authuser")
	if err != nil {
		fmt.Println("error ocurred during request: ", err)
	}
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		fmt.Println("Invalid JSON:", err)
		return
	}

	fmt.Println(prettyJSON.String())
}
