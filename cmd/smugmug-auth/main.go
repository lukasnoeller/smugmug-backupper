package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/lukasnoeller/smugmug-backupper/internal/auth"
)

func main() {
	// Get api key and secret
	credentials, err := auth.GetCredentials("SMUGMUG_API_KEY", "SMUGMUG_API_SECRET")
	if err != nil {
		fmt.Println("Error ocurred during retrieval of credentials: ", err)
		return
	}
	// Step 1: Get temporary request token
	fmt.Println("Connecting to smug  mug to receive temporary request token...")
	requestTokenurl := "https://api.smugmug.com/services/oauth/1.0a/getRequestToken"
	params := map[string]string{
		"oauth_callback":         "oob",
		"oauth_consumer_key":     credentials["SMUGMUG_API_KEY"],
		"oauth_nonce":            auth.GenerateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_version":          "1.0",
	}
	params["oauth_signature"] = auth.CalculateSignature("POST", requestTokenurl, params, credentials["SMUGMUG_API_SECRET"], "")

	tempAccessTokenResp, err := auth.MakeOauthRequest("POST", requestTokenurl, params)
	if err != nil {
		fmt.Println("error occurred during oauth request ", err)
		return
	}
	tempOauthToken := tempAccessTokenResp.Get("oauth_token")
	tempOauthSecret := tempAccessTokenResp.Get("oauth_token_secret")
	if tempOauthToken == "" || tempOauthSecret == "" {
		fmt.Println("temporary oauthToken and or oauthSecret could not be retrieved!")
		fmt.Println("response values: ", tempAccessTokenResp)
		return
	}
	fmt.Println("oauth_token: ", tempOauthToken)
	fmt.Println("oauth_secret: ", tempOauthSecret)

	// Step 2: Authorize in the Browser
	authUrl := fmt.Sprintf("https://api.smugmug.com/services/oauth/1.0a/authorize?oauth_token=%s", tempOauthToken)
	fmt.Println("Open your browser and authorize smug-mug-backupper to have access to your galleries at following link:")
	fmt.Println(authUrl)

	var pin string
	fmt.Printf("\nEnter 6-digit PIN here: ")
	fmt.Scanln(&pin)
	pin = strings.TrimSpace(pin)

	// Step 3
	accessTokenUrl := "https://api.smugmug.com/services/oauth/1.0a/getAccessToken"
	accessTokenParams := map[string]string{
		"oauth_consumer_key":     credentials["SMUGMUG_API_KEY"],
		"oauth_nonce":            auth.GenerateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_token":            tempOauthToken,
		"oauth_verifier":         pin,
		"oauth_version":          "1.0",
	}
	accessTokenParams["oauth_signature"] = auth.CalculateSignature("POST", accessTokenUrl, accessTokenParams, credentials["SMUGMUG_API_SECRET"], tempOauthSecret)
	accessTokenResp, err := auth.MakeOauthRequest("POST", requestTokenurl, params)
	if err != nil {
		fmt.Println("error occurred during oauth request ", err)
		return
	}
	outputCredentials := make(map[string]string)
	outputCredentials["SMUGMUG_ACCESS_TOKEN"] = accessTokenResp.Get("oauth_token")
	outputCredentials["SMUGMUG_ACCESS_TOKEN_SECRET"] = accessTokenResp.Get("oauth_token_secret")
	auth.SetCredentials(j)
	if outputCredentials["SMUGMUG_ACCESS_TOKEN"] == "" || outputCredentials["SMUGMUG_ACCESS_TOKEN_SECRET"] == "" {
		fmt.Println("oauthToken and or oauthSecret could not be retrieved!")
		fmt.Println("response values: ", accessTokenResp)
		return
	}
	fmt.Println("smugmug access token and secret values generated and set to environment")

}
