package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
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
	requestTokenurl := "https://api.smugmug.com/services/oauth/1.0a/getRequestToken"
	params := map[string]string{
		"oauth_callback":         "oob",
		"oauth_consumer_key":     consumerKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_version":          "1.0",
	}
	params["oauth_signature"] = generateSignature("POST", requestTokenurl, params, consumerSecret, "")

	tempAccessTokenResp, err := makeOauthRequest("POST", requestTokenurl, params)
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
		"oauth_consumer_key":     consumerKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_token":            tempOauthToken,
		"oauth_verifier":         pin,
		"oauth_version":          "1.0",
	}
	accessTokenParams["oauth_signature"] = generateSignature("POST", accessTokenUrl, accessTokenParams, consumerSecret, tempOauthSecret)
	accessTokenResp, err := makeOauthRequest("POST", requestTokenurl, params)
	if err != nil {
		fmt.Println("error occurred during oauth request ", err)
		return
	}
	oauthToken := accessTokenResp.Get("oauth_token")
	oauthSecret := accessTokenResp.Get("oauth_token_secret")
	if oauthToken == "" || oauthSecret == "" {
		fmt.Println("oauthToken and or oauthSecret could not be retrieved!")
		fmt.Println("response values: ", accessTokenResp)
		return
	}
	fmt.Println("oauth_token: ", oauthToken)
	fmt.Println("oauth_secret: ", oauthSecret)

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

func generateSignature(method, targetUrl string, params map[string]string, consumerSecret, tokenSecret string) string {

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var paramParts []string
	for _, k := range keys {
		part := url.QueryEscape(k) + "=" + url.QueryEscape(params[k])
		paramParts = append(paramParts, part)
	}
	paramString := strings.Join(paramParts, "&")
	signatureBaseString := url.QueryEscape(method) + "&" +
		url.QueryEscape(targetUrl) + "&" +
		url.QueryEscape(paramString)

	signingKey := url.QueryEscape(consumerSecret) + "&" + url.QueryEscape(tokenSecret)

	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(signatureBaseString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))

}

func makeOauthRequest(method, targetUrl string, params map[string]string) (url.Values, error) {
	var authParts []string
	for k, v := range params {
		authParts = append(authParts, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	authHeader := "Oauth " + strings.Join(authParts, ", ")
	req, err := http.NewRequest(method, targetUrl, nil)
	if err != nil {
		println("error occurred during intialization request object")
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Length", "0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("error occurred during oauth request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("error occurred during read of response body")
		return nil, err
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		println("error occurred during parsing of query")
		println("body: ", string(body))
		return nil, err
	}
	return values, nil
}
