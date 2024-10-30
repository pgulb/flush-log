package intest

import (
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func TestLoadPages(t *testing.T)() {
	for _, url := range Endpoints() {
		u := launcher.New().Bin(LauncherSystemBrowser()).MustLaunch()
		b := rod.New().ControlURL(u).MustConnect()
		p := b.MustPage(url)
		defer b.MustClose()
		defer p.MustClose()
	}
}

func TestCheckForErrorDivId(t *testing.T)() {
	for _, url := range Endpoints() {
		u := launcher.New().Bin(LauncherSystemBrowser()).MustLaunch()
		b := rod.New().ControlURL(u).MustConnect()
		p := b.MustPage(url)
		defer b.MustClose()
		defer p.MustClose()
		p.MustElement("#error").MustWaitInvisible()
	}
}
