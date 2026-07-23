package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dghubble/oauth1"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

func SetCredentials(credentials map[string]string) {
	for key, value := range credentials {
		os.Setenv(key, value)
	}
}
func GetCredentials(credentials ...string) (map[string]string, error) {
	mapCredentials := make(map[string]string)
	for _, credential := range credentials {
		value := os.Getenv(credential)
		if value == "" {
			return nil, fmt.Errorf("credential: %s not found", credential)
		}
		mapCredentials[credential] = value
	}
	return mapCredentials, nil
}
func GenerateNonce() string {
	// generates a number used once
	const charset = "absdefhhijklmnoprstuvwxyzABCDEDFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func CalculateSignature(method, targetUrl string, params map[string]string, consumerSecret, tokenSecret string) string {

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

func MakeOauthRequest(method, targetUrl string, params map[string]string) (url.Values, error) {
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

func AuthenticatedRequest(method, RequestUrl string) ([]byte, error) {

	// Get credentials
	credentials, err := GetCredentials("SMUGMUG_API_KEY", "SMUGMUG_API_SECRET", "SMUGMUG_ACCESS_TOKEN", "SMUGMUG_ACCESS_TOKEN_SECRET")
	if err != nil {
		fmt.Println("Error ocurred during retrieval of credentials: ", err)
		return nil, err
	}

	config := oauth1.NewConfig(credentials["SMUGMUG_API_KEY"], credentials["SMUGMUG_API_SECRET"])
	token := oauth1.NewToken(credentials["SMUGMUG_ACCESS_TOKEN"], credentials["SMUGMUG_ACCESS_TOKEN_SECRET"])
	httpClient := config.Client(oauth1.NoContext, token)
	req, err := http.NewRequest(method, RequestUrl, nil)
	if err != nil {
		fmt.Println("error ocurred during request generation: ", err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error ocurred during http request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body, nil
}
