package main

import (
	"fmt"
	"io"
	"net/http"

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
	httpClient := config.Client(oauth1.NoContext, token)
	req, err := http.NewRequest("GET", "https://api.smugmug.com/api/v2/folder/user/utanoller", nil)
	if err != nil {
		fmt.Println("error ocurred during request generation: ", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error ocurred during http request: ", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", body)
}
