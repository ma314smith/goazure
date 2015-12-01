package goazure

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCallEndpoint(t *testing.T) {
	Convey("Given a valid ServieBusRelay object", t, func() {
		acs := ACS{IssuerName: os.Getenv("GOAZURE_ACSISSUERNAME"), IssuerKey: os.Getenv("GOAZURE_ACSISSUERKEY")}
		sbr := ServiceBusRelay{Namespace: "blackbaud", Scope: "servicerouter-qas", AccessControl: &acs}
		Convey("When an endpoint is called", func() {
			endpointPath := "/SFDC/AccountTeamService"
			soapAction := "http://webservices.blackbaud.com/sfdc/accountteamservice/GetAccountTeamBySiteID"
			soapBody := "<GetAccountTeamBySiteID xmlns=\"http://webservices.blackbaud.com/sfdc/accountteamservice/\"><siteID>5740</siteID></GetAccountTeamBySiteID>"
			result, err := sbr.CallEndpoint(endpointPath, soapAction, soapBody)
			Convey("Then a valid response is returned", func() {
				So(err, ShouldBeNil)
				So(string(result), ShouldContainSubstring, "GetAccountTeamBySiteIDResponse")
			})
		})
	})
	Convey("Given an invalid ServieBusRelay object", t, func() {
		acs := ACS{IssuerName: "invalid", IssuerKey: os.Getenv("GOAZURE_ACSISSUERKEY")}
		sbr := ServiceBusRelay{Namespace: "blackbaud", Scope: "servicerouter-qas", AccessControl: &acs}
		Convey("When an endpoint is called", func() {
			endpointPath := "SFDC/AccountTeamService"
			soapAction := "http://webservices.blackbaud.com/sfdc/accountteamservice/GetAccountTeamBySiteID"
			soapBody := "<GetAccountTeamBySiteID xmlns=\"http://webservices.blackbaud.com/sfdc/accountteamservice/\"><siteID>5740</siteID></GetAccountTeamBySiteID>"
			result, err := sbr.CallEndpoint(endpointPath, soapAction, soapBody)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(string(result), ShouldEqual, "")
			})
		})
	})
}
