package goazure

import (
	"encoding/base64"
	"fmt"
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
func (a *ACS) GenerateUUID() string {
	return uuid.NewV4().String()
}

// GetToken returns a token that can be used to authenticate to Azure resources
func (a *ACS) GetToken(namespace string, scope *url.URL) (token string, err error) {
	req, err := a.buildRequest(namespace, scope)
	if err != nil {
		return "", err
	}

	body, err := sendHTTPRequest(req)
	if err != nil {
		return "", err
	}

	return a.processToken(body)
}

func (a *ACS) buildRequest(namespace string, scope *url.URL) (*http.Request, error) {
	acsURL := "https://" + namespace + "-sb.accesscontrol.windows.net/WRAPv0.9/"

	data := url.Values{}
	data.Add("wrap_name", a.IssuerName)
	data.Add("wrap_password", a.IssuerKey)
	data.Add("wrap_scope", scope.String())

	req, err := http.NewRequest("POST", acsURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("goazure: %s", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (a *ACS) processToken(body []byte) (token string, err error) {
	responseBodyValues := strings.Split(string(body), "&")
	tokenKeyValue := responseBodyValues[0]
	//expirationKeyValue := responseBodyValues[1]
	if !strings.Contains(tokenKeyValue, "wrap_access_token=") {
		return "", fmt.Errorf("goazure: ACS Authentication Failed")
	}
	tokenValue := strings.Replace(tokenKeyValue, "wrap_access_token=", "", 1)
	tokenValue, err = url.QueryUnescape(tokenValue)
	if err != nil {
		return "", fmt.Errorf("goazure: %s", err)
	}

	token = base64.StdEncoding.EncodeToString([]byte(tokenValue))
	return token, nil
}
