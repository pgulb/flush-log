package intest

import (
	"testing"

	"github.com/go-rod/rod"
)

func TestLoadPages(t *testing.T)() {
	for _, url := range Endpoints() {
		b := rod.New().MustConnect()
		p := b.MustPage(url)
		p.MustWaitStable()
		p.MustClose()
		b.MustClose()
	}
}

func TestCheckForErrorDivId(t *testing.T)() {
	for _, url := range Endpoints() {
		b := rod.New().MustConnect()
		p := b.MustPage(url)
		p = p.MustWaitStable()
		p.MustElement("#error").MustWaitInvisible()
		p.MustClose()
		b.MustClose()
	}
}
