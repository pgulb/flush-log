package test

import (
	"testing"

	f "github.com/pgulb/flush-log/flush"
)

func TestGetApiUrl(t *testing.T) {
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
