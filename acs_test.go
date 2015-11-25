package goazure

import (
	"net/url"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetToken(t *testing.T) {
	Convey("Given a valid ACS config", t, func() {
		acs := ACS{IssuerName: os.Getenv("GOAZURE_ACSISSUERNAME"), IssuerKey: os.Getenv("GOAZURE_ACSISSUERKEY")}
		Convey("When a token is requested", func() {
			scope, _ := url.Parse("http://blackbaud.servicebus.windows.net/servicerouter-qas")
			token, err := acs.GetToken("blackbaud", scope)
			Convey("Then a token is returned", func() {
				So(err, ShouldBeNil)
				So(token, ShouldNotEqual, "")
			})
		})
	})
}
