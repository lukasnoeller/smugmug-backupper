package user

import (
	"encoding/json"
	"fmt"
	"github.com/lukasnoeller/smugmug-backupper/internal/auth"
)

type User struct {
	Name     string `json:"Name"`
	NickName string `json:"NickName"`
	Uris     struct {
		Node struct {
			Uri string `json:"Uri"`
		} `json:"Node"`
	} `json:"Uris"`
}
type AuthUserResponse struct {
	Response struct {
		User User `json:"User"`
	} `json:"Response"`
}

func GetUser() User {

	body, err := auth.AuthenticatedRequest("GET", "https://api.smugmug.com/api/v2!authuser")
	var authUserResponse AuthUserResponse
	err = json.Unmarshal([]byte(body), &authUserResponse)
	if err != nil {
		fmt.Println("error ocurred during unmarshalling of response body: ", err)
	}
	authUser := authUserResponse.Response.User
	return authUser
}
