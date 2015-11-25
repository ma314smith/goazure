package goazure

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// ServiceBusRelay is used to call Azure Service Bus Relay endpoints.  ACS
// authentication is used.
type ServiceBusRelay struct {
	Namespace     string
	Scope         string
	AccessControl *ACS
}

// CallEndpoint makes a SOAP call against the endpoint provided
func (s *ServiceBusRelay) CallEndpoint(endpointPath string, soapAction string, soapBody string) (responseBody []byte, err error) {
	if endpointPath[0] != []byte("/")[0] {
		endpointPath = "/" + endpointPath
	}
	endpointURL := "https://" + s.Namespace + ".servicebus.windows.net/" + s.Scope + endpointPath

	envelope, err := s.buildSOAPEnvelope(soapBody)
	if err != nil {
		return nil, fmt.Errorf("goazure: error calling buildSOAPEnvelope: %s", err)
	}

	req, err := http.NewRequest("POST", endpointURL, strings.NewReader(envelope))
	if err != nil {
		return nil, fmt.Errorf("goazure: %s", err)
	}
	req.Header.Add("SOAPAction", soapAction)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")

	// send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("goazure: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("goazure: response status code: %v", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func (s *ServiceBusRelay) buildSOAPEnvelope(soapBody string) (string, error) {
	uuid := s.AccessControl.GenerateUUID()
	scopeURL, err := url.Parse("http://" + s.Namespace + ".servicebus.windows.net/" + s.Scope)
	if err != nil {
		return "", err
	}
	token, err := s.AccessControl.GetToken(s.Namespace, scopeURL)
	if err != nil {
		return "", err
	}

	requestBody := "<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"<s:Header>" +
		"<RelayAccessToken xmlns=\"http://schemas.microsoft.com/netservices/2009/05/servicebus/connect\">" +
		"<wsse:BinarySecurityToken wsu:Id=\"uuid:" + uuid + "\" ValueType=\"http://schemas.xmlsoap.org/ws/2009/11/swt-token-profile-1.0\" " +
		"EncodingType=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary\" " +
		"xmlns:wsse=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd\" " +
		"xmlns:wsu=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd\">" +
		token +
		"</wsse:BinarySecurityToken></RelayAccessToken>" +
		"</s:Header>" +
		"<s:Body>" +
		soapBody +
		"</s:Body></s:Envelope>"

	return requestBody, nil
}
