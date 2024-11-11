package intest

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestCheckPhoneUsedAsDefault(t *testing.T) {
	p, b := RegisterAndGoToSettings(
		"user_settings_test",
		"pass_settings_test",
		"pass_settings_test",
	)
	defer b.MustClose()
	defer p.MustClose()
	p.MustElement("#phone-used-default").MustClick()
	p = p.MustNavigate(os.Getenv("GOAPP_URL") + "/new")
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	checked := p.MustElement("#new-flush-phone-used").MustProperty("checked").Bool()
	log.Println("#new-flush-phone-used: ", checked)
	if !checked {
		t.Fatal("phone used not set as default")
	}
}
