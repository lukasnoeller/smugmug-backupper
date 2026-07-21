package auth

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
