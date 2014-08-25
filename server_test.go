/*
	Tests for server.
	Run with go test -v
*/
package main

import (
	"bytes"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

var (
	baseURL       = "http://127.0.0.1:3000/api/"
	test_email    = "foo"
	test_password = "bar"
)

func createJSON(user *User, t *testing.T) *bytes.Buffer {
	data, err := json.Marshal(user)
	if err != nil {
		t.FailNow()
	}
	return bytes.NewBuffer(data)
}

func TestServer(t *testing.T) {
	t.Parallel()
	Convey("Should be able to do request", t, func() {
		res, err := http.Get(baseURL + "version")
		So(err, ShouldBeNil)
		Convey("Should be able to get current version", func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})
	Convey("Should be able to login with user", t, func() {
		client := new(http.Client)
		user := &User{Email: test_email, Password: test_password}
		data := createJSON(user, t)
		res, err := client.Post(baseURL+"login", "application/json", data)
		if err != nil {
			t.FailNow()
		}
		So(res.StatusCode, ShouldEqual, 200)
	})
}
