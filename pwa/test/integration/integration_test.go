package intest

import (
	"os"
	"testing"

	"github.com/go-rod/rod"
)

func TestEnvIsSet(t *testing.T) {
	if os.Getenv("GOAPP_URL") == "" {
		t.Fatal("GOAPP_URL not set")
	}
}

func TestLoginPage(t *testing.T)() {
    page := rod.New().MustConnect().MustPage(os.Getenv("GOAPP_URL"))
    page.MustWaitStable().MustElementByJS("document.querySelector('body')")
}
