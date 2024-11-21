package test

import (
	"log"
	"testing"
	"time"

	f "github.com/pgulb/flush-log/flush"
)

func TestValidateRegistryCreds(t *testing.T) {
	t.Parallel()
	cases := []f.LastTriedCreds{
		{
			User:     "test",
			Password: "test",
		},
		{
			User:     "testtest",
			Password: "testtest",
		},
		{
			User:     "testtesttesttest",
			Password: "testtesttesttest",
		},
	}
	if err := f.ValidateRegistryCreds(
		"test", "test", "test", cases[0]); err == nil {
		t.Fatal("should error on already used creds")
	}
	if err := f.ValidateRegistryCreds(
		"test", "test", "test", cases[1]); err != nil {
		t.Fatal(err)
	}
	for _, c := range [][]interface{}{
		{"", "test", "test", cases[2]},
		{"test", "", "test", cases[2]},
		{"test", "test", "", cases[2]},
		{"", "", "", cases[2]},
		{"test", "", "", cases[2]},
		{"", "test", "", cases[2]},
		{"", "test", "test", cases[2]},
	} {
		if err := f.ValidateRegistryCreds(
			c[0].(string), c[1].(string), c[2].(string), c[3].(f.LastTriedCreds),
		); err == nil {
			t.Fatal("should error on already used creds")
		}
	}
}

func TestValidateLoginCreds(t *testing.T) {
	t.Parallel()
	cases := []f.LastTriedCreds{
		{
			User:     "test",
			Password: "test",
		},
		{
			User:     "testtest",
			Password: "test",
		},
	}

	if err := f.ValidateLoginCreds(
		"test", "test", cases[0]); err == nil {
		t.Fatal("should error on already used creds")
	}
	if err := f.ValidateLoginCreds(
		"testtest", "test", cases[1]); err == nil {
		t.Fatal("should error on already used creds")
	}
	if err := f.ValidateLoginCreds(
		"test", "", f.LastTriedCreds{}); err == nil {
		t.Fatal("should error on empty pass")
	}
	if err := f.ValidateLoginCreds(
		"", "test", f.LastTriedCreds{}); err == nil {
		t.Fatal("should error on empty user")
	}
	if err := f.ValidateLoginCreds(
		"", "", f.LastTriedCreds{}); err == nil {
		t.Fatal("should error on empty creds")
	}
}

func TestValidateFlush(t *testing.T) {
	t.Parallel()
	failCases := []f.Flush{
		{
			TimeStart: time.Now(),
			TimeEnd:   time.Now().Add(-time.Hour), // end before start
			Rating:    5,
			PhoneUsed: true,
			Note:      "test",
		},
		{
			TimeStart: time.Now().Add(time.Hour),
			TimeEnd:   time.Now(),
			Rating:    11, // out of range
			PhoneUsed: true,
			Note:      "test",
		},
		{
			TimeStart: time.Time{},
			TimeEnd:   time.Now(),
			Rating:    0, // out of range
			PhoneUsed: true,
			Note:      "test",
		},
		{
			TimeStart: time.Time{},
			TimeEnd:   time.Now(),
			Rating:    5,
			PhoneUsed: true,
			Note: `asdasdasdasdasdasdasdasdasdasdasda
			sdasdasdasdasdasdasdasdasdasdasdasdasdasdasdas
			dasdasdasdasdasdasdasdasdasdasdasdasdasd`, // too long
		},
	}

	okCases := []f.Flush{
		{
			TimeStart: time.Now(),
			TimeEnd:   time.Now().Add(time.Hour),
			Rating:    5,
			PhoneUsed: true,
			Note:      "test",
		},
		{
			TimeStart: time.Now(),
			TimeEnd:   time.Now().Add(time.Hour),
			Rating:    10,
			PhoneUsed: false,
			Note:      "dfgs lkjghsd uhgriewop hguifd sohbvxpoxchvbp oihs pdfobhpwerb hwpefobhi sdpfbhsdpfboishdrg sdpfgoish",
		},
		{
			TimeStart: time.Now().Add(-time.Hour),
			TimeEnd:   time.Now().Add(time.Hour),
			Rating:    1,
			PhoneUsed: false,
			Note:      `k;\ÎºIàðãÀv,ëm,ÓÅèïº)ÿ7dèÉó×EàrXê3liæWãÎtÅÑ%²Xse²w¥=ðÿë8+Îå6_ÊÁ¶w£,!Iaú¾%¤×øzNíæ¤æ\ØåÐo7\ÿ`,
		},
	}

	for _, c := range failCases {
		log.Println(c)
		if err := f.ValidateFlush(c); err == nil {
			t.Fatal("should error on invalid flush")
		}
	}
	for _, c := range okCases {
		log.Println(c)
		if err := f.ValidateFlush(c); err != nil {
			t.Fatal(err)
		}
	}
}
