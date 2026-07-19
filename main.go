package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	url := "https://api.smugmug.com/api/v2"
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
