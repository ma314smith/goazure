package goazure

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func sendHTTPRequest(req *http.Request) (responseBody []byte, err error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("goazure: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("goazure: failure status code: %v", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
