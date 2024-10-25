package intest

import (
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	p := Register("user_registry_test", "pass_registry_test",
	"pass_registry_test")
	defer p.MustClose()
	e := p.MustElement("#fetched-flushes")
	if !strings.Contains(e.MustText(), "user_registry_test") {
		t.Fatal("user not registered")
	}
}

func TestRegisterBadUsernameChars(t *testing.T) {
	p := Register("user_ęśąćź", "pass_registry_test", "pass_registry_test")
	defer p.MustClose()
	if err := CheckPageForErrorDivVisible(p); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterEmptyUsername(t *testing.T) {
	p := Register("", "pass_registry_test", "pass_registry_test")
	defer p.MustClose()
	if err := CheckPageForErrorDivVisible(p); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterShortPass(t *testing.T)()  {
	p := Register("user_registry_test2", "asd", "asd")
	defer p.MustClose()
	if err := CheckPageForErrorDivVisible(p); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterEmptyPass(t *testing.T)()  {
	p := Register("user_registry_test3", "", "")
	defer p.MustClose()
	if err := CheckPageForErrorDivVisible(p); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterPasswordsDiffer(t *testing.T)()  {
	p := Register("user_registry_test3",
		"asdasdasdasdasddasd", "qweqweqweqweqwe")
	defer p.MustClose()
	if err := CheckPageForErrorDivVisible(p); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterUsernameTaken(t *testing.T)()  {
	p := Register("user_registry_test", "pass_registry_test",
	"pass_registry_test")
	defer p.MustClose()
	if err := CheckPageForErrorDivVisible(p); err == nil {
		t.Fatal("should throw error when used the same username twice")
	}
}
