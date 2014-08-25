/*
	Tests for server.
	Run with go test -v
*/

package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestGetVersion(t *testing.T) {
	t.Parallel()

	Convey("Should be able to get current version", t, func() {
		res, err := http.Get("http://127.0.0.1:8080/api/version")
		So(err, ShouldBeNil)
		Convey("Should be able to get current version")
		So(res.StatusCode, ShouldEqual, 200)
	})
}
