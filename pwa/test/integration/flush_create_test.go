package intest

import (
	"testing"
	"time"
)

func TestCreateFlush(t *testing.T) {
	p, b, err := CreateFlush(true, "flush_create_test", "flush_create_test", time.Now())
	defer b.MustClose()
	defer p.MustClose()
	if err != nil {
		t.Fatal(err)
	}
}
