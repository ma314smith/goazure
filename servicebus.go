package goazure

// ServiceBusRelay is used to call Azure Service Bus Relay endpoints.  ACS
// authentication is used.
type ServiceBusRelay struct {
	Namespace     string
	Scope         string
	ACSIssuerName string
	ACSIssuerKey  string
}

// CallEndpoint makes a SOAP call against the endpoint provided
func (s *ServiceBusRelay) CallEndpoint(endpointPath string, soapAction string, soapBody string) {
	if endpointPath[0] != []byte("/")[0] {
		endpointPath = "/" + endpointPath
	}
	//endpointURL := "https://" + s.Namespace + ".servicebus.windows.net/" + s.Scope + endpointPath

}

func (s *ServiceBusRelay) buildSOAPPayload(soapBody string) {

}
