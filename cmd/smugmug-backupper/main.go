package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/lukasnoeller/smugmug-backupper/internal/auth"
)

func main() {

	// Get credentials
	credentials, err := auth.GetCredentials("SMUGMUG_API_KEY", "SMUGMUG_API_SECRET", "SMUGMUG_ACCESS_TOKEN", "SMUGMUG_ACCESS_TOKEN_SECRET")
	if err != nil {
		fmt.Println("Error ocurred during retrieval of credentials: ", err)
		return
	}

	config := oauth1.NewConfig(credentials["SMUGMUG_API_KEY"], credentials["SMUGMUG_API_SECRET"])
	token := oauth1.NewToken(credentials["SMUGMUG_ACCESS_TOKEN"], credentials["SMUGMUG_ACCESS_TOKEN_SECRET"])
	method := "GET"
	targetUrl := "https://api.smugmug.com/api/v2/folder"
	req, err := http.NewRequest(method, targetUrl, nil)
	if err != nil {
		println("error occurred during intialization request object")
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)

}
