package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	yellowButtonCss = "font-bold bg-yellow-500 p-2 rounded text-white mx-1"
	errorDivCss = "flex flex-row fixed bottom-4 left-4 bg-red-500 text-white p-4 text-xl rounded-lg"
	centeringDivCss = "flex flex-row min-h-screen justify-center items-center"
	registerDivCss = "p-4 text-center text-xl shadow-lg bg-white rounded-lg mx-10"
	inviCss = "fixed invisible"
)

type creds struct {
	UserColonPass string
	LoggedIn      bool
}

type lastTriedCreds struct {
	User string
	Password string
}

type ErrorContainer struct {
	app.Compo
}
func (e *ErrorContainer) Render() app.UI {
	return app.Div().Body(app.Div().Body(
		app.P().Text("placeholder error")).Class(
			"p-8 text-center text-xl shadow-lg bg-white rounded-lg",
		)).Class(
			inviCss,
			).ID("error")
}

type buttonShowRegister struct {
	app.Compo
}
func (b *buttonShowRegister) Render() app.UI {
	return app.Button().Text("I need account").OnClick(b.onClick).Class(yellowButtonCss)
}
func (b *buttonShowRegister) onClick(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("register-container").Set("className", registerDivCss)
}

type registerContainer struct {
	app.Compo
}
func (r *registerContainer) Render() app.UI {
	return app.Div().Body(
		app.P().Text("Register").Class("font-bold"),
		app.Input().Type("text").ID("register-username").Placeholder("Username").Class(
			"m-2 placeholder-gray-800",
		),
		app.Br(),
		app.Input().Type("password").ID("register-password").Placeholder("Password").Class(
			"m-2 placeholder-gray-800",
		),
		app.Br(),
		app.Input().Type("password").ID("register-password-repeat").Placeholder(
			"Repeat password").Class(
			"m-2 placeholder-gray-800 my-4",
		),
		app.Br(),
		&buttonRegister{},
	).Class(inviCss).ID("register-container")
}

type rootContainer struct {
	app.Compo
	buttonUpdate
}
func (b *rootContainer) OnMount(ctx app.Context) {
	var creds creds
	ctx.GetState("creds", &creds)
	log.Println("Logged in: ", creds.LoggedIn)
	if !creds.LoggedIn {
		app.Window().Set("location", "login")
	} else {
		ctx.Async(func() {
			m, err := getFlushesFromApi(ctx)
			if err != nil {
				showErrorDiv(ctx, err)
			} else {
				app.Window().GetElementByID("fetched-flushes").Set("innerHTML", m)
			}
		})
	}
}
func (b *rootContainer) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Flush Log").Class("text-2xl"),
		app.P().Text("Tracked flushes:").Class("py-2"),
		app.P().Text("").Class(
			"py-2",
		).ID("fetched-flushes"),
		b.buttonUpdate.Render(),
		&buttonLogout{},
		app.Div().Body(&ErrorContainer{}),
	).Class("shadow-lg bg-white rounded-lg p-6 min-h-72 relative")
}

type buttonUpdate struct {
	app.Compo
}
func (b *buttonUpdate) Render() app.UI {
	return app.Button().Text("Update").OnClick(b.onClick).Class(
		"bg-yellow-500 hover:bg-yellow-700 text-black font-bold py-2 px-4 rounded absolute bottom-4 left-4")
}
func (b *buttonUpdate) onClick(ctx app.Context, e app.Event) {
	var creds creds
	ctx.GetState("creds", &creds)
	ctx.Async(func() {
		if creds.LoggedIn {
			log.Println("Getting new API response...")
			m, err := getFlushesFromApi(ctx)
			if err != nil {
				
			} else {
				app.Window().GetElementByID("fetched-flushes").Set("innerHTML", m)
			}
		}})
}

type loginContainer struct {
	app.Compo
}
func (l *loginContainer) Render() app.UI {
	return app.Div().Body(app.Div().Body(
		app.P().Text("Log in to continue.").Class("font-bold"),
		app.Div().Body(
			app.Input().Type("text").ID("username").Placeholder("Username").Class(
				"m-2 placeholder-gray-800",
			),
			app.Br(),
			app.Input().Type("password").ID("password").Placeholder("Password").Class(
				"m-2 placeholder-gray-800",
			),
			app.Br(),
			app.Div().Body(
				app.Input().Type("checkbox").ID("remember-me").Class("m-2"),
				app.Label().For("remember-me").Text("Remember me").Class("p-2"),
			),
			app.Br(),
			app.Div().Body(
				&buttonLogin{},
				&buttonShowRegister{},
			),
		),
	).Class("p-4 text-center text-xl shadow-lg bg-white rounded-lg"),
	&registerContainer{},
	app.Div().Body(&ErrorContainer{})).Class(
		centeringDivCss)
}

type buttonLogin struct {
	app.Compo
}
func (b *buttonLogin) Render() app.UI {
	return app.Button().Text("Log in").OnClick(b.onClick).Class(
		yellowButtonCss)
}
func (b *buttonLogin) onClick(ctx app.Context, e app.Event) {
	loginSeconds := 60
	log.Println("Trying to log in...")
	if app.Window().GetElementByID("remember-me").Get("checked").Bool() {
		log.Println("remember-me checked")
		loginSeconds = 604800 // week
	}
	lastCreds := lastTriedCreds{}
	ctx.GetState("lastUsedCreds", &lastCreds)
	user := app.Window().GetElementByID("username").Get("value").String()
	pass := app.Window().GetElementByID("password").Get("value").String()
	if user == "" || pass == "" {
		showErrorDiv(ctx, errors.New("username and password required"))
		return
	}
	if lastCreds.User == user && lastCreds.Password == pass {
		log.Println("Skipping last used credentials...")
		showErrorDiv(ctx, errors.New("you already tried those credentials"))
		return
	}
	ctx.Async(func() {
		status, basic_auth, err := tryLogin(user, pass)
		if err != nil {
			showErrorDiv(ctx, err)
		}
		switch status {
		case 200:
			ctx.SetState("creds", creds{
				UserColonPass: basic_auth,
				LoggedIn:      true,
			}).ExpiresIn(time.Second * time.Duration(loginSeconds)).PersistWithEncryption()
			log.Printf("Logged in as %s\n for %v seconds!", user, loginSeconds)
			app.Window().Set("location", ".")
			ctx.DelState("lastUsedCreds")
		case 401:
			showErrorDiv(ctx, errors.New("invalid credentials"))
			ctx.SetState("lastUsedCreds", lastTriedCreds{
				User: user,
				Password: pass,
			}).ExpiresIn(time.Second * 10)
		default:
			showErrorDiv(ctx, errors.New("login failed"))
			ctx.DelState("lastUsedCreds")
		}
})
}

type buttonRegister struct {
	app.Compo
}
func (b *buttonRegister) Render() app.UI {
	return app.Button().Text("Register").OnClick(b.onClick).Class(yellowButtonCss)
}
func (b *buttonRegister) onClick(ctx app.Context, e app.Event) {
	log.Println("Trying to register...")
	lastCreds := lastTriedCreds{}
	ctx.GetState("lastUsedCredsRegister", &lastCreds)
	user := app.Window().GetElementByID("register-username").Get("value").String()
	pass := app.Window().GetElementByID("register-password").Get("value").String()
	repeatPass := app.Window().GetElementByID("register-password-repeat").Get("value").String()
	if user == "" || pass == "" || repeatPass == "" {
		showErrorDiv(ctx, errors.New("fill all required fields"))
		return
	}
	if pass != repeatPass {
		showErrorDiv(ctx, errors.New("passwords don't match"))
		return
	}
	if lastCreds.User == user && lastCreds.Password == pass {
		log.Println("Skipping last used credentials...")
		showErrorDiv(ctx, errors.New("you already tried those credentials"))
		return
	}
	ctx.Async(func() {
		status, basic_auth, err := tryRegister(user, pass)
		if err != nil {
			showErrorDiv(ctx, err)
		}
		switch status {
		case 201:
			ctx.SetState("creds", creds{
				UserColonPass: basic_auth,
				LoggedIn:      true,
			}).ExpiresIn(time.Second * time.Duration(604800)).PersistWithEncryption()
			ctx.DelState("lastUsedCredsRegister")
			app.Window().Set("location", ".")
		case 422:
			showErrorDiv(ctx, errors.New("invalid username or password"))
			ctx.SetState("lastUsedCredsRegister", lastTriedCreds{
				User: user,
				Password: pass,
			}).ExpiresIn(time.Second * 10)
		case 409:
			showErrorDiv(ctx, errors.New("username already exists"))
			ctx.SetState("lastUsedCredsRegister", lastTriedCreds{
				User: user,
				Password: pass,
			}).ExpiresIn(time.Second * 10)
		default:
			showErrorDiv(ctx, errors.New("register failed"))
			ctx.DelState("lastUsedCredsRegister")
		}
})
}

type buttonLogout struct {
	app.Compo
}
func (b *buttonLogout) Render() app.UI {
	return app.Button().Text("Log out").OnClick(b.onClick).Class(
		"font-bold border-2 border-white p-2 rounded absolute bottom-4 right-4")
}
func (b *buttonLogout) onClick(ctx app.Context, e app.Event) {
	ctx.SetState("creds", creds{LoggedIn: false}).PersistWithEncryption()
	app.Window().Set("location", ".")
}

func closeBody(r *http.Response) {
	if err := r.Body.Close(); err != nil { displayError(err) } 
}

func getApiUrl() (string, error) {
	r, err := http.Get("web/apiurl")
	if err != nil {
		return "", err
	}
	defer closeBody(r)
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func tryLogin(username string, password string) (int, string, error) {
	apiUrl, err := getApiUrl()
	if err != nil {
		return 0, "", err
	}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return 0, "", err
	}
	basic_auth := base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)
	req.Header.Add("Authorization", "Basic "+basic_auth)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer closeBody(r)
	return r.StatusCode, basic_auth, nil
}

func tryRegister(username string, password string) (int, string, error) {
	apiUrl, err := getApiUrl()
	if err != nil {
		return 0, "", err
	}
	js := []byte(fmt.Sprintf(`
	{
		 "username": "%s",
		  "password": "%s"
	}
	`, username, password))
	req, err := http.NewRequest("POST", apiUrl+"/user", bytes.NewBuffer(js))
	if err != nil {
		return 0, "", err
	}
	basic_auth := base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer closeBody(r)
	return r.StatusCode, basic_auth, nil
}

func getFlushesFromApi(ctx app.Context) (string, error) {
	apiUrl, err := getApiUrl()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return "", err
	}
	var c creds
	ctx.GetState("creds", &c)
	req.Header.Add("Authorization", "Basic "+c.UserColonPass)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer closeBody(r)
	if r.StatusCode >= 400 {
		ctx.SetState("creds", creds{LoggedIn: false}).PersistWithEncryption()
		app.Window().Set("location", "login")
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func displayError(err error) {
	app.Window().Call("alert", "An error occurred:\n\n"+err.Error()+"\n\nRefresh the page to continue.")
	log.Fatal(err)
}

func showErrorDiv(ctx app.Context, err error) {
	app.Window().GetElementByID("error").Set("innerHTML", err.Error())
	app.Window().GetElementByID("error").Set("className", errorDivCss)
	ctx.Async(func() {
		time.Sleep(time.Second * 1)
		app.Window().GetElementByID("error").Set("className", inviCss)
	})
	// TODO consider using https://developer.mozilla.org/en-US/docs/Web/API/Node/cloneNode
	// clone the element and set random ID to clone
	// then hide the clone after 2 seconds
}

func main() {
	app.Route("/", func() app.Composer {
		return &rootContainer{}
	})
	app.Route("/login", func() app.Composer { return &loginContainer{} })
	app.RunWhenOnBrowser()

	if os.Getenv("BUILD_STATIC") == "true" {
		err := app.GenerateStaticWebsite(".", &app.Handler{
			Name:        "Flush-Log",
			Description: "bowel tracking app",
			Resources:   app.GitHubPages("flush-log"),
			Scripts: []string{
				"https://cdn.tailwindcss.com",
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	port := os.Getenv("GOAPP_PORT")
	if port == "" {
		log.Fatal("GOAPP_PORT not set")
	}
	apiUrlBytes, err := os.ReadFile("web/apiurl")
	if err != nil {
		log.Fatal(err)
	}
	apiUrl := string(apiUrlBytes)

	log.Println("Listening on " + port)
	log.Println("API url: ", apiUrl)
	http.Handle("/", &app.Handler{
		Name:        "Flush-Log",
		Description: "bowel tracking app",
		Scripts: []string{
			"https://cdn.tailwindcss.com",
		},
	})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
