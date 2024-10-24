package test

import (
	"testing"

	f "github.com/pgulb/flush-log/flush"
)

func TestValidateRegistryCreds(t *testing.T) {
	cases := []f.LastTriedCreds{
		{
			User:     "test",
			Password: "test",
		},
		{
			User:     "testtest",
			Password: "testtest",
		},
		{
			User:     "testtesttesttest",
			Password: "testtesttesttest",
		},
	}
	if err := f.ValidateRegistryCreds(
		"test", "test", "test", cases[0]); err == nil {
		t.Fatal("should error on already used creds")
	}
	if err := f.ValidateRegistryCreds(
		"test", "test", "test", cases[1]); err != nil {
		t.Fatal(err)
	}
	for _, c := range [][]interface{}{
		{"", "test", "test", cases[2]},
		{"test", "", "test", cases[2]},
		{"test", "test", "", cases[2]},
		{"", "", "", cases[2]},
		{"test", "", "", cases[2]},
		{"", "test", "", cases[2]},
		{"", "test", "test", cases[2]},
	} {
		if err := f.ValidateRegistryCreds(
			c[0].(string), c[1].(string), c[2].(string), c[3].(f.LastTriedCreds),
			); err == nil {
			t.Fatal("should error on already used creds")
		}
	}
}

func TestValidateLoginCreds(t *testing.T) {
	cases := []f.LastTriedCreds{
		{
			User:     "test",
			Password: "test",
		},
		{
			User:     "testtest",
			Password: "test",
		},
	}

	if err := f.ValidateLoginCreds(
		"test", "test", cases[0]); err == nil {
		t.Fatal("should error on already used creds")
	}
	if err := f.ValidateLoginCreds(
		"testtest", "test", cases[1]); err == nil {
		t.Fatal("should error on already used creds")
	}
	if err := f.ValidateLoginCreds(
		"test", "", f.LastTriedCreds{}); err == nil {
		t.Fatal("should error on empty pass")
	}
	if err := f.ValidateLoginCreds(
		"", "test", f.LastTriedCreds{}); err == nil {
		t.Fatal("should error on empty user")
	}
	if err := f.ValidateLoginCreds(
		"", "", f.LastTriedCreds{}); err == nil {
		t.Fatal("should error on empty creds")
	}
}
