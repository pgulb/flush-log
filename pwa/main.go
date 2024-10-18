package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type creds struct {
	UserColonPass string
	LoggedIn bool
}

type hello struct {
	app.Compo
}
func (h *hello) Render() app.UI {
	return app.H1().Text("Flush Log").Class("text-3xl")
}

type paragraphLoremIpsum struct {
	app.Compo
}
func (p *paragraphLoremIpsum) Render() app.UI {
	return app.P().Text(`Data from API will appear below.`).Class(
		"py-2",
	)
}

type rootContainer struct {
	app.Compo
	hello
	paragraphLoremIpsum
	messageFromApi string
	buttonUpdate
}
func (b *rootContainer) OnMount(ctx app.Context) {
	var creds creds
	ctx.GetState("creds", &creds)
	// ctx.GetState("loggedIn", &loggedIn)
	log.Println("Logged in: ", creds.LoggedIn)
	if !creds.LoggedIn {
		app.Window().Set("location", "/login")
	} else {
		b.messageFromApi = getDataFromApi(ctx)
	}
}
func (b *rootContainer) Render() app.UI {
	b.buttonUpdate.message = &b.messageFromApi
	return app.Div().Body(
		b.hello.Render(),
		b.paragraphLoremIpsum.Render(),
		app.P().Text(b.messageFromApi).Class(
			"py-2",
		),
		b.buttonUpdate.Render(),
		&buttonLogout{},
	).Class("shadow-lg bg-white rounded-lg p-6 min-h-72 relative")
}

type buttonUpdate struct {
	app.Compo
	message *string
}
func (b *buttonUpdate) Render() app.UI {
	return app.Button().Text("Update").OnClick(b.onClick).Class(
		"bg-yellow-500 hover:bg-yellow-700 text-black font-bold py-2 px-4 rounded absolute bottom-4 left-4",)
}
func (b *buttonUpdate) onClick(ctx app.Context, e app.Event) {
	var creds creds
	ctx.GetState("creds", &creds)
	if creds.LoggedIn {
		log.Println("Getting new API response...")
		m := getDataFromApi(ctx)
		*b.message = m
	}
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
			&buttonLogin{},
		),
	).Class("p-4 text-center text-xl shadow-lg bg-white rounded-lg",)).Class(
		"flex flex-row min-h-screen justify-center items-center")
}

type buttonLogin struct {
	app.Compo
}
func (b *buttonLogin) Render() app.UI {
	return app.Button().Text("Log in").OnClick(b.onClick).Class(
		"font-bold bg-yellow-500 p-2 rounded text-white",)
}
func (b *buttonLogin) onClick(ctx app.Context, e app.Event) {
	log.Println("Logging in...")
	user := app.Window().GetElementByID("username").Get("value").String()
	pass := app.Window().GetElementByID("password").Get("value").String()

	if user == "admin" && pass == "admin" {
		ctx.SetState("creds", creds{
			UserColonPass: base64.StdEncoding.EncodeToString([]byte(user + ":" + pass)),
			LoggedIn: true,
			}).ExpiresIn(time.Second * 60).PersistWithEncryption()
		app.Window().Set("location", "/")
	}
}

type buttonLogout struct {
	app.Compo
}
func (b *buttonLogout) Render() app.UI {
	return app.Button().Text("Log out").OnClick(b.onClick).Class(
		"font-bold border-2 border-white p-2 rounded absolute bottom-4 right-4",)
}
func (b *buttonLogout) onClick(ctx app.Context, e app.Event) {
	ctx.SetState("creds", creds{LoggedIn: false}).PersistWithEncryption()
	app.Window().Set("location", "/")
}

func getDataFromApi(ctx app.Context) string {
	r, err := http.Get("/web/apiurl")
	if err != nil {
		displayError(err)
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		displayError(err)
	}
	apiUrl := string(bytes)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		displayError(err)
	}
	var c creds
	ctx.GetState("creds", &c)
	req.Header.Add("Authorization", "Basic "+c.UserColonPass)

	r, err = http.DefaultClient.Do(req)
	if r.StatusCode == 401 {
		ctx.SetState("creds", creds{LoggedIn: false}).PersistWithEncryption()
		app.Window().Set("location", "/login")
	}
	if err != nil {
		displayError(err)
	}
	defer func ()  {
		err := r.Body.Close()
		if err != nil {
			displayError(err)
		}
	}()
	bytes, err = io.ReadAll(r.Body)
	if err != nil {
		displayError(err)
	}
	return string(bytes)
}

func displayError(err error) {
	app.Window().Call("alert", "An error occurred:\n\n"+err.Error())
	log.Fatal(err)
}

func main() {
	app.Route("/", func() app.Composer { return &rootContainer{
	} })
	app.Route("/login", func() app.Composer { return &loginContainer{} })
	app.RunWhenOnBrowser()

	if os.Getenv("BUILD_STATIC") == "true" {
		err := app.GenerateStaticWebsite(".", &app.Handler{
			Name:        "Flush-Log",
			Description: "bowel tracking app",
			Resources:   app.GitHubPages("flush-log"),
		})
		if err != nil {
			log.Fatal(err)
		}
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

	log.Println("Listening on "+port)
	log.Println("API url: ", apiUrl)
	http.Handle("/", &app.Handler{
		Name: "Flush-Log",
		Description: "bowel tracking app",
		Scripts: []string{
			"https://cdn.tailwindcss.com",
		},
	})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
