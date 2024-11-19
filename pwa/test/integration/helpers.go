package intest

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func Endpoints() []string {
	base := os.Getenv("GOAPP_URL")
	return []string{
		base,
		base + "/login",
	}
}

func LauncherSystemBrowser() string {
	path, has := launcher.LookPath()
	if !has {
		log.Fatal("browser not installed")
	}
	log.Println("using browser: ", path)
	return path
}

func LoginPage() (*rod.Page, *rod.Browser) {
	log.Println("using LoginPage()")
	u := launcher.New().Bin(LauncherSystemBrowser()).MustLaunch()
	b := rod.New().ControlURL(u).MustConnect()
	p := b.MustPage(os.Getenv("GOAPP_URL") + "/login")
	go p.EachEvent(func(e *proto.RuntimeConsoleAPICalled) {
		if e.Type == proto.RuntimeConsoleAPICalledTypeLog {
			fmt.Println("BROWSER:", p.MustObjectsToJSON(e.Args))
		}
	})()
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("return from LoginPage()")
	return p, b
}

func Register(user string, pass string,
	repeatPass string) (*rod.Page, *rod.Browser) {
	p, b := LoginPage()
	log.Println("using Register()")
	p.MustElement("#show-register").MustClick()
	p.MustElement("#register-username").MustInput(user)
	p.MustElement("#register-password").MustInput(pass)
	p.MustElement("#register-password-repeat").MustInput(repeatPass)
	p.MustElement("#register-button").MustClick()
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	if err := CheckErrorDivText(p, "register failed"); err == nil {
		log.Fatal("#error contains 'register failed' text")
	}
	log.Println("return from Register()")
	return p, b
}

func RegisterDoubleClickButton(user string, pass string,
	repeatPass string) (*rod.Page, *rod.Browser) {
	p, b := Register(user, pass, repeatPass)
	log.Println("using RegisterDoubleClickButton()")
	p.MustElement("#register-button").MustClick()
	log.Println("return from RegisterDoubleClickButton()")
	return p, b
}

func CheckErrorDivText(p *rod.Page, text string) error {
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

func RegisterAndGoToNew(user string, pass string,
	repeatPass string) (*rod.Page, *rod.Browser) {
	log.Println("using RegisterAndGoToNew()")
	p, b := Register(user, pass, repeatPass)
	p.Navigate(os.Getenv("GOAPP_URL") + "/new")
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("return from RegisterAndGoToNew()")
	return p, b
}

func RegisterAndGoToSettings(user string, pass string,
	repeatPass string) (*rod.Page, *rod.Browser) {
	log.Println("using RegisterAndGoToSettings()")
	p, b := Register(user, pass, repeatPass)
	p.Navigate(os.Getenv("GOAPP_URL") + "/settings")
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("return from RegisterAndGoToSettings()")
	return p, b
}

func Login(username string, password string) (*rod.Page, *rod.Browser) {
	log.Println("using Login()")
	p, b := LoginPage()
	p.MustElement("#username").MustInput(username)
	p.MustElement("#password").MustInput(password)
	p.MustElement("#login-button").MustClick()
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("return from Login()")
	return p, b
}

func CreateFlush(
	register bool,
	user string,
	pass string,
	timeStart time.Time,
) (*rod.Page, *rod.Browser, error) {
	var p *rod.Page
	var b *rod.Browser
	if register {
		p, b = RegisterAndGoToNew(user, pass, pass)
	} else {
		p, b = Login(user, pass)
		err := p.Navigate(os.Getenv("GOAPP_URL") + "/new")
		if err != nil {
			return p, b, err
		}
		err = p.WaitStable(time.Second * 2)
		if err != nil {
			return p, b, err
		}
	}
	p.MustElement("#new-flush-time-start").MustInputTime(timeStart)
	p.MustElement("#new-flush-time-end").MustInputTime(timeStart.Add(time.Minute * 10))
	p.MustElement("#new-flush-rating").MustInput("5")
	p.MustElement("#new-flush-phone-used").MustClick()
	p.MustElement("#new-flush-note").MustInput("test comment")
	p.MustElement("#submit-flush-button").MustClick()
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		return p, b, err
	}
	if err := CheckErrorDivText(p, "placeholder"); err != nil {
		return p, b, err
	}
	return p, b, nil
}
