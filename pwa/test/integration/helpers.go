package intest

import (
	"errors"
	"fmt"
	"log"
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

func LoginPage() (*rod.Page, *rod.Browser) {
	b := rod.New().MustConnect()
	p := b.MustPage(os.Getenv("GOAPP_URL")+"/login")
	p = p.MustWaitStable()
	return p, b
}

func Register(user string, pass string,
	repeatPass string) (*rod.Page, *rod.Browser) {
	p, b := LoginPage()
	p.MustElement("#show-register").MustClick()
	p.MustElement("#register-username").MustInput(user)
	p.MustElement("#register-password").MustInput(pass)
	p.MustElement("#register-password-repeat").MustInput(repeatPass)
	p.MustElement("#register-button").MustClick()
	p = p.MustWaitStable()
	return p, b
}

func RegisterDoubleClickButton(user string, pass string,
	repeatPass string) (*rod.Page, *rod.Browser) {
	p, b := Register(user, pass, repeatPass)
	p.MustElement("#register-button").MustClick()
	p = p.MustWaitStable()
	return p, b
}

func CheckErrorDivText(p *rod.Page, text string) error {
	p.MustWaitDOMStable()
	e := p.MustElement("#error")
	errText := e.MustProperty("innerHTML").String()
	log.Printf("#error - '%s'\n", errText)
	if !strings.Contains(errText, text) {
		return fmt.Errorf(
			"#error div - not found '%s' in '%s'", text, errText)
	}
	return nil
}

func CheckRegisterHintVisible(p *rod.Page) error {
	e := p.MustElement("#register-error")
	if strings.Contains(e.Object.ClassName, "invisible") {
		return errors.New("#credentials hint still invisible")
	}
	return nil
}
