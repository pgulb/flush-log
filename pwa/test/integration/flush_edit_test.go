package intest

import (
	"testing"
	"time"
)

func TestEditFlush(t *testing.T) {
	p, b, err := CreateFlush(true, "flush_edit_test", "flush_edit_test", time.Now())
	defer b.MustClose()
	defer p.MustClose()
	if err != nil {
		t.Fatal(err)
	}
	p.MustElementR("button", "Edit 🛠️").MustClick()
	err = p.WaitStable(time.Second * 2)
	if err != nil {
		t.Fatal(err)
	}
	p.MustElement("#edited-time-start").MustInputTime(time.Now())
	p.MustElement("#edited-time-end").MustInputTime(time.Now().Add(time.Hour))
	p.MustElement("#edited-rating").MustInput("9")
	p.MustElement("#edited-phone-used").MustClick()
	p.MustElement("#edited-note").
		MustInput("test test t esttest asdf asedfasdf asedfas deśąśęąśęźżź")
	p.MustElement("#edit-submit-button").MustClick()
	err = p.WaitStable(time.Second * 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckErrorDivText(p, "placeholder"); err != nil {
		t.Fatal(err)
	}
}
