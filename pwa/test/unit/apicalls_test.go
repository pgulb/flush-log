package test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	f "github.com/pgulb/flush-log/flush"
)

func TestBasicAuth(t *testing.T) {
	cases := []map[string]string{
		{
			"username": "test",
			"password": "test",
			"expected":  "dGVzdDp0ZXN0",
		},
		{
			"username": "hergo45uhg324uhg5_123123",
			"password": "asdfasdfasdfasdf234t234fv234f234f2vdfgsdfgĘŚĄĆĄŚĆęęśąąćźżźć",
			"expected": "aGVyZ280NXVoZzMyNHVoZzVfMTIzMTIzOmFzZGZhc2RmYXNkZmFzZGYyMzR0MjM0ZnYyMzRmMjM0ZjJ2ZGZnc2RmZ8SYxZrEhMSGxITFmsSGxJnEmcWbxIXEhcSHxbrFvMW6xIc=",
		},
		{
			"username": "juzer",
			"password": `êö2Q¯ÀÂÀ#9û\Z¥VIG#Ñüär¡ßZ®f£Ûo´p·ÓDµÙÏ2oûó'Ô/Ôã"t"`,
			"expected": "anV6ZXI6w6rDtjJRwq/DgMOCw4AjOcO7XFrCpVZJRyPDkcO8w6RywqHDn1rCrmbCo8Obb8K0cMK3w5NEwrXDmcOPMm/Du8OzJ8OUL8OUw6MidCI=",
		},
	}

	for _, c := range cases {
		basicAuth := f.BasicAuth(c["username"], c["password"])
		if basicAuth != c["expected"] {
			t.Fatal("expected", c["expected"], "got", basicAuth)
		}
	}
}

func TestCloseBody(t *testing.T) {
	r := &http.Response{
		Body: io.NopCloser(&bytes.Buffer{}),
	}
	f.CloseBody(r)
}

func TestAuthorizedRequest(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"method":   "GET",
			"url":      "http://localhost",
			"basicAuth": f.BasicAuth("test", "test"),
			"expected":  http.Request{
				Method: "GET",
				URL:    &url.URL{
					Scheme: "http",
					Host: "localhost",
				},
				Header: http.Header{"Authorization": []string{"Basic dGVzdDp0ZXN0"}},
			},
		},
		{
			"method":   "POST",
			"url":      "https://example.com",
			"basicAuth": f.BasicAuth("test2", "test2"),
			"expected":  http.Request{
				Method: "POST",
				URL:    &url.URL{
					Scheme: "https",
					Host: "example.com",
				},
				Header: http.Header{"Authorization": []string{"Basic dGVzdDI6dGVzdDI="}},
			},
		},
		{
			"method":   "PUT",
			"url":      "https://localhost:1234",
			"basicAuth": f.BasicAuth("test3", "test3"),
			"expected":  http.Request{
				Method: "PUT",
				URL:    &url.URL{
					Scheme: "https",
					Host: "localhost:1234",
				},
				Header: http.Header{"Authorization": []string{"Basic dGVzdDM6dGVzdDM="}},
			},
		},
	}
	for _, c := range cases {
		req, err := f.AuthorizedRequest(
			c["method"].(string), c["url"].(string), c["basicAuth"].(string))
		if err != nil {
			t.Fatal(err)
		}
		if req.URL.String() != c["expected"].(http.Request).URL.String() {
			t.Fatal("expected", c["expected"].(http.Request).URL, "got", req.URL)
		}
		if req.Method != c["expected"].(http.Request).Method {
			t.Fatal("expected", c["expected"].(http.Request).Method, "got", req.Method)
		}
		if req.Header.Get("Authorization") != c["expected"].(http.Request).Header.Get("Authorization") {
			t.Fatal("expected", c["expected"].(http.Request).Header.Get("Authorization"), "got", req.Header.Get("Authorization"))
		}
	}
}
