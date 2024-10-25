package intest

import (
	"errors"
	"os"
	"strings"

	"github.com/go-rod/rod"
)

func Endpoints() []string {
	base := os.Getenv("GOAPP_URL")
	return []string{
		base,
		base + "/login",
	}
}

func LoginPage() *rod.Page {
	p := rod.New().MustConnect().MustPage(os.Getenv("GOAPP_URL")+"/login")
	return p.MustWaitStable()
}

func Register(user string, pass string,
	repeatPass string) *rod.Page {
	p := LoginPage()
	p.MustElement("#show-register").MustClick()
	p.MustElement("#register-username").MustInput(user)
	p.MustElement("#register-password").MustInput(pass)
	p.MustElement("#register-password-repeat").MustInput(repeatPass)
	p.MustElement("#register-button").MustClick()
	return p.MustWaitStable()
}

func CheckPageForErrorDivVisible(p *rod.Page) error {
	e := p.MustElement("#error")
	if strings.Contains(e.Object.ClassName, "invisible") {
		return errors.New("#error div still invisible")
	}
	return nil
}
