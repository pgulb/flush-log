package intest

import (
	"os"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func endpoints() []string {
	base := os.Getenv("GOAPP_URL")
	return []string{
		base,
		base + "/login",
	}
}

func TestEnvIsSet(t *testing.T) {
	if os.Getenv("GOAPP_URL") == "" {
		t.Fatal("GOAPP_URL not set")
	}
}

func TestLoadPages(t *testing.T)() {
	for _, url := range endpoints() {
		p, err := rod.New().MustConnect().Page(proto.TargetCreateTarget{
			URL: url})
		if err != nil {
			t.Fatal(err)
		}
		p.MustWaitStable()
		err = p.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCheckForErrorDivId(t *testing.T)() {
	for _, url := range endpoints() {
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
		err = p.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}
