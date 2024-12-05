package test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	f "github.com/pgulb/flush-log/flush"
)

func TestFormatFlushTime(t *testing.T) {
	t.Parallel()
	a := make([]int, 31)
	for i := range a {
		a[i] = i
	}
	for _, i := range a {
		now := time.Now()
		subtractedTime := now.Add(time.Duration(i*-24) * time.Hour)
		result := f.FormatFlushTime(subtractedTime, subtractedTime)
		log.Println(result)
		if i == 7 {
			if !strings.Contains(result, "week ago, ") {
				t.Errorf("FormatFlushTime() = %v, want %v", result, "week ago, ")
			}
		}
		if i <= 6 && i >= 2 {
			if !strings.Contains(result, fmt.Sprintf("%v days ago, ", i)) {
				t.Errorf("FormatFlushTime() = %v, want %v", result, fmt.Sprintf("%v days ago, ", i))
			}
		}
		if i == 1 {
			if !strings.Contains(result, "yesterday, ") {
				t.Errorf("FormatFlushTime() = %v, want %v", result, "yesterday, ")
			}
		}
		if i == 0 {
			if !strings.Contains(result, "today, ") {
				t.Errorf("FormatFlushTime() = %v, want %v", result, "today, ")
			}
		}
	}
}
