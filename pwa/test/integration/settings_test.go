package intest

import (
	"log"
	"os"
	"strings"
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

func TestChangePassword(t *testing.T) {
	p, b := RegisterAndGoToSettings(
		"user_settings_test2",
		"pass_settings_test2",
		"pass_settings_test2",
	)
	defer b.MustClose()
	defer p.MustClose()
	p.MustElement("#chp-password").MustInput("new_pass_settings_test2")
	p.MustElement("#chp-password-repeat").MustInput("new_pass_settings_test2")
	p.MustElement("#chp-button").MustClick()
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	p2, b2 := Login("user_settings_test2", "new_pass_settings_test2")
	defer b2.MustClose()
	defer p2.MustClose()
	e := p2.MustElement("#hidden-hello")
	if !strings.Contains(e.MustProperty("innerHTML").String(), "hello!") {
		t.Fatal("user did not login successfully")
	}
}

func TestAccountRemoval(t *testing.T) {
	p, b := RegisterAndGoToSettings(
		"user_settings_test3",
		"pass_settings_test3",
		"pass_settings_test3",
	)
	defer b.MustClose()
	defer p.MustClose()
	p.MustElement("#remove-account-button").MustClick()
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	if err := CheckErrorDivText(p, "Type 'byebye' before deleting account"); err != nil {
		t.Fatal(err)
	}
	p.MustElement("#remove-account-byebye").MustInput("byebye")
	p.MustElement("#remove-account-button").MustClick()
	err = p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	if err := CheckErrorDivText(p, "placeholder"); err != nil {
		t.Fatal(err)
	}
}
