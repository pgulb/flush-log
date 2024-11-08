package flush

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func SetLastUsedCredsState(ctx app.Context, user string, pass string) {
	ctx.SetState("lastUsedCredsRegister", LastTriedCreds{
		User:     user,
		Password: pass,
	}).ExpiresIn(time.Second * 10)
}

func ValidateRegistryCreds(
	user string, pass string, repeatPass string, lastCreds LastTriedCreds) error {
	if user == "" || pass == "" || repeatPass == "" {
		return errors.New("fill all required fields")
	}
	if pass != repeatPass {
		return errors.New("passwords don't match")
	}
	if lastCreds.User == user && lastCreds.Password == pass {
		log.Println("Skipping last used credentials...")
		return errors.New("you already tried those credentials")
	}
	return nil
}

func ValidateLoginCreds(
	user string, pass string, lastCreds LastTriedCreds) error {
	if user == "" || pass == "" {
		return errors.New("username and password required")
	}
	if lastCreds.User == user && lastCreds.Password == pass {
		log.Println("Skipping last used credentials...")
		return errors.New("you already tried those credentials")
	}
	return nil
}

func GetRegisterCreds() (string, string, string) {
	user := app.Window().GetElementByID("register-username").Get("value").String()
	pass := app.Window().GetElementByID("register-password").Get("value").String()
	repeatPass := app.Window().GetElementByID("register-password-repeat").Get("value").String()
	return user, pass, repeatPass
}

func GetLoginCreds() (string, string) {
	user := app.Window().GetElementByID("username").Get("value").String()
	pass := app.Window().GetElementByID("password").Get("value").String()
	return user, pass
}

func ShowBadRegisterCredsErr() {
	app.Window().GetElementByID(
		"register-error").Set(
		"innerHTML", "username: up to 60 chars, (letters, numbers and _),<br/>password: 8-60 chars")
}

func DisplayError(err error) {
	app.Window().
		Call("alert", "An error occurred:\n\n"+err.Error()+"\n\nRefresh the page to continue.")
	log.Fatal(err)
}

func ShowErrorDiv(ctx app.Context, err error, seconds time.Duration) {
	log.Println("Error: ", err)
	app.Window().GetElementByID("error").Set("innerHTML", err.Error())
	app.Window().GetElementByID("error").Set("className", ErrorDivCss)
	ctx.Async(func() {
		time.Sleep(time.Second * seconds)
		app.Window().GetElementByID("error").Set("className", InviCss)
	})
	// TODO consider using https://developer.mozilla.org/en-US/docs/Web/API/Node/cloneNode
	// clone the element and set random ID to clone
	// then hide the clone after 2 seconds
}

func NewFLush(ctx app.Context, timeStartRaw string, timeEndRaw string,
	ratingRaw string, phoneUsed bool, note string) (Flush, error) {
	log.Println("timeStart raw: ", timeStartRaw)
	log.Println("timeEnd raw: ", timeEndRaw)
	timeStart, err := time.Parse("2006-01-02T15:04", timeStartRaw)
	if err != nil {
		return Flush{}, errors.New("time-start and time-end are required")
	}
	timeEnd, err := time.Parse("2006-01-02T15:04", timeEndRaw)
	if err != nil {
		return Flush{}, errors.New("time-start and time-end are required")
	}
	rating, err := strconv.Atoi(ratingRaw)
	if err != nil {
		return Flush{}, err
	}
	flush := Flush{
		TimeStart: timeStart,
		TimeEnd:   timeEnd,
		Rating:    rating,
		PhoneUsed: phoneUsed,
		Note:      note,
	}
	log.Println("New flush: ", flush)
	return flush, nil
}

func ValidateFlush(flush Flush) error {
	if flush.TimeStart.After(flush.TimeEnd) {
		return errors.New("start-time must be before end-time")
	}
	if flush.Rating < 1 || flush.Rating > 10 {
		return errors.New("rating must be between 1 and 10")
	}
	if len([]rune(flush.Note)) > 100 {
		return errors.New("note must be at most 100 characters")
	}
	return nil
}
