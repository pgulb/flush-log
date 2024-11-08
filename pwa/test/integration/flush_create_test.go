package intest

import (
	"testing"
	"time"
)

func TestCreateFlush(t *testing.T) {
	p, b := RegisterAndGoToNew("user_flush_test", "pass_flush_test", "pass_flush_test")
	defer b.MustClose()
	defer p.MustClose()
	p.MustElement("#new-flush-time-start").MustInputTime(time.Now())
	p.MustElement("#new-flush-time-end").MustInputTime(time.Now().Add(time.Hour))
	p.MustElement("#new-flush-rating").MustInput("5")
	p.MustElement("#new-flush-phone-used").MustClick()
	p.MustElement("#new-flush-note").MustInput("test comment")
	p.MustElement("#submit-flush-button").MustClick()
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckErrorDivText(p, "placeholder"); err != nil {
		t.Fatal(err)
	}
}
