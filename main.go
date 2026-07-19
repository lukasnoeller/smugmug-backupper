package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// Get api key and secret
	consumerKey := os.Getenv("SMUGMUG_API_KEY")
	consumerSecret := os.Getenv("SMUGMUG_API_SECRET")
	if consumerKey == "" || consumerSecret == "" {
		fmt.Println("Error: missing env variable values..")
		fmt.Println("SMUGMUG_API_KEY and/or SMUGMUG_API_SECRET missing")
		return
	}
	// Step 1: Get temporary request token
	fmt.Println("Connecting to smug  mug to receive temporary request token...")
	requestTokenurl := "https://api.smugmug.com/services/oath/1.0a/getRequestToken"
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error ocurred during HTTP request: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error while reading body: ", err)
		return
	}
	fmt.Println(string(body))

}

func generateNonce() string {
	// generates a number used once
	const charset = "absdefhhijklmnoprstuvwxyzABCDEDFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func makeOauthRequest(method, targetUrl string, params map[string]string) (url.Values, error) {
	var authParts []string
	for k, v := range params {
		authParts = append(authParts, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	authHeader := "Oauth " + strings.Join(authParts, ", ")
	req, err := http.NewRequest(method, targetUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Length", "0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	values , err := 

}
