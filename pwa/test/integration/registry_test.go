package intest

import (
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	p, b := Register("user_registry_test", "pass_registry_test",
	"pass_registry_test")
	defer b.MustClose()
	defer p.MustClose()
	e := p.MustElement("#hidden-hello")
	if !strings.Contains(e.MustText(), "user_registry_test") {
		t.Fatal("user not registered")
	}
}

func TestRegisterBadUsernameChars(t *testing.T) {
	p, b := Register("user_ęśąćź", "pass_registry_test", "pass_registry_test")
	defer b.MustClose()
	defer p.MustClose()
	if err := CheckRegisterHintVisible(p); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterEmptyUsername(t *testing.T) {
	p, b := Register("", "pass_registry_test", "pass_registry_test")
	defer b.MustClose()
	defer p.MustClose()
	if err := CheckRegisterHintVisible(p); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterShortPass(t *testing.T)()  {
	p, b := Register("user_registry_test2", "asd", "asd")
	defer b.MustClose()
	defer p.MustClose()
	if err := CheckRegisterHintVisible(p); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterEmptyPass(t *testing.T)()  {
	p, b := Register("user_registry_test3", "", "")
	defer b.MustClose()
	defer p.MustClose()
	if err := CheckErrorDivText(p, "fill all required fields"); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterPasswordsDiffer(t *testing.T)()  {
	p, b := Register("user_registry_test3",
		"asdasdasdasdasddasd", "qweqweqweqweqwe")
		defer b.MustClose()
		defer p.MustClose()
	if err := CheckErrorDivText(p, "passwords don't match"); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterUsernameTaken(t *testing.T)()  {
	p, b := Register("user_registry_test3", "pass_registry_test3",
	"pass_registry_test3")
	p2, b2 := Register("user_registry_test3", "pass_registry_test3",
	"pass_registry_test3")
	defer b2.MustClose()
	defer p2.MustClose()
	defer b.MustClose()
	defer p.MustClose()
	if err := CheckErrorDivText(p2, "username already exists"); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterSameCredsTwoTimes(t *testing.T)()  {
	p, b := RegisterDoubleClickButton("asd", "asd",
	"asd")
	defer b.MustClose()
	defer p.MustClose()
	if err := CheckErrorDivText(p, "you already tried those credentials"); err != nil {
		t.Fatal(err)
	}
}
