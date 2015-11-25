package goazure

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/satori/go.uuid"
)

// ACS handles authentication with the Azure Access Control Service
type ACS struct {
	IssuerName string
	IssuerKey  string
}

// GenerateUUID returns a UUID that can be used in the Azure request
func (a ACS) GenerateUUID() string {
	return uuid.NewV4().String()
}

// GetToken returns a token that can be used to authenticate to Azure resources
func (a ACS) GetToken(namespace string, scope *url.URL) (string, error) {

	// set up the request
	acsURL := "https://" + namespace + "-sb.accesscontrol.windows.net/WRAPv0.9/"

	data := url.Values{}
	data.Add("wrap_name", a.IssuerName)
	data.Add("wrap_password", a.IssuerKey)
	data.Add("wrap_scope", scope.String())

	req, err := http.NewRequest("POST", acsURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("goazure: %s", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("goazure: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("goazure: %s", err)
	}

	// parse the token from the response body
	tokenKeyValue := strings.Split(string(body), "&")[0]
	if !strings.Contains(tokenKeyValue, "wrap_access_token=") {
		return "", fmt.Errorf("goazure: ACS Authentication Failed")
	}
	token := strings.Replace(tokenKeyValue, "wrap_access_token=", "", 1)
	token, err = url.QueryUnescape(token)
	if err != nil {
		return "", fmt.Errorf("goazure: %s", err)
	}

	b64EncodedToken := base64.StdEncoding.EncodeToString([]byte(token))

	return b64EncodedToken, nil
}
