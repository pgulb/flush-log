package intest

import (
	"os"
	"testing"
)

func TestEnvIsSet(t *testing.T) {
	t.Parallel()
	if os.Getenv("GOAPP_URL") == "" {
		t.Fatal("GOAPP_URL not set")
	}
}
