package intest

import (
	"testing"
)

func TestFeedback(t *testing.T) {
	user, pass := "feedback_create_test", "feedback_create_test"
	p, b, err := CreateFeedback(
		true,
		user,
		pass,
		"testing feedback text lolololololol",
		"placeholder",
	)
	defer b.MustClose()
	defer p.MustClose()
	if err != nil {
		t.Fatal(err)
	}
	p2, b2, err := CreateFeedback(
		false,
		user,
		pass,
		"too short",
		"Feedback too short",
	)
	defer b2.MustClose()
	defer p2.MustClose()
	if err != nil {
		t.Fatal(err)
	}
}
