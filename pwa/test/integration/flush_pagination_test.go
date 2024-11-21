package intest

import (
	"log"
	"testing"
	"time"
)

func TestFlushPagination(t *testing.T) {
	now := time.Now()
	for i := 0; i < 6; i++ {
		register := false
		if i == 0 {
			register = true
		}
		p, b, err := CreateFlush(
			register,
			"flush_pagination_test",
			"flush_pagination_test",
			now.Add(time.Duration(i)*time.Hour),
		)
		defer b.MustClose()
		defer p.MustClose()
		if err != nil {
			t.Fatal(err)
		}
	}
	p, b := Login("flush_pagination_test", "flush_pagination_test")
	defer b.MustClose()
	defer p.MustClose()
	els := p.MustElements(".remove-flush-buttonz-div")
	log.Println("number of flushes: ", len(els))
	if len(els) != 3 {
		t.Fatal("wrong number of flushes")
	}
	p.MustElement("#update-button").MustClick()
	err := p.WaitStable(time.Second * 2)
	if err != nil {
		log.Fatal(err)
	}
	els = p.MustElements(".remove-flush-buttonz-div")
	log.Println("number of flushes: ", len(els))
	if len(els) != 6 {
		t.Fatal("wrong number of flushes")
	}
}
