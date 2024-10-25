package intest

import (
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func TestLoadPages(t *testing.T)() {
	for _, url := range Endpoints() {
		p, err := rod.New().MustConnect().Page(proto.TargetCreateTarget{
			URL: url})
		if err != nil {
			t.Fatal(err)
		}
		p.MustWaitStable()
		p.MustClose()
	}
}

func TestCheckForErrorDivId(t *testing.T)() {
	for _, url := range Endpoints() {
		p, err := rod.New().MustConnect().Page(proto.TargetCreateTarget{
			URL: url})
		if err != nil {
			t.Fatal(err)
		}
		p = p.MustWaitStable()
		err = p.MustElement("#error").WaitInvisible()
		if err != nil {
			t.Fatal(err)
		}
		p.MustClose()
	}
}
