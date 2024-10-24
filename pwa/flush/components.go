package flush

import (
	"errors"
	"log"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	YellowButtonCss = "font-bold bg-yellow-500 p-2 rounded text-white mx-1"
	ErrorDivCss = "flex flex-row fixed bottom-4 left-4 bg-red-500 text-white p-4 text-xl rounded-lg"
	CenteringDivCss = "flex flex-row min-h-screen justify-center items-center"
	RegisterDivCss = "p-4 text-center text-xl shadow-lg bg-white rounded-lg mx-10"
	InviCss = "fixed invisible"
)

type ErrorContainer struct {
	app.Compo
}
func (e *ErrorContainer) Render() app.UI {
	return app.Div().Body(app.Div().Body(
		app.P().Text("placeholder error")).Class(
			"p-8 text-center text-xl shadow-lg bg-white rounded-lg",
		)).Class(
			InviCss,
			).ID("error")
}

type buttonShowRegister struct {
	app.Compo
}
func (b *buttonShowRegister) Render() app.UI {
	return app.Button().Text("I need account").OnClick(b.onClick).Class(YellowButtonCss)
}
func (b *buttonShowRegister) onClick(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("register-container").Set("className", RegisterDivCss)
	app.Window().GetElementByID("login-container").Set("className", InviCss)
}

type RegisterContainer struct {
	app.Compo
}
func (r *RegisterContainer) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
		app.P().Text("Register").Class("font-bold"),
		app.Input().Type("text").ID("register-username").Placeholder("Username").Class(
			"m-2 placeholder-gray-500",
		),
		app.Br(),
		app.Input().Type("password").ID("register-password").Placeholder("Password").Class(
			"m-2 placeholder-gray-500",
		),
		app.Br(),
		app.Input().Type("password").ID("register-password-repeat").Placeholder(
			"Repeat password").Class(
			"m-2 placeholder-gray-500 my-4",
		),
		app.Br(),
		&buttonRegister{},
		app.P().Text("").Class("text-red-500").ID("register-error"),
	).Class(InviCss).ID("register-container"),
	)
}

type RootContainer struct {
	app.Compo
	buttonUpdate
}
func (b *RootContainer) OnMount(ctx app.Context) {
	var creds Creds
	ctx.GetState("creds", &creds)
	log.Println("Logged in: ", creds.LoggedIn)
	if !creds.LoggedIn {
		app.Window().Set("location", "login")
	} else {
		ctx.Async(func() {
			m, err := GetFlushesFromApi(ctx)
			if err != nil {
				ShowErrorDiv(ctx, err, 1)
			} else {
				app.Window().GetElementByID("fetched-flushes").Set("innerHTML", m)
			}
		})
	}
}
func (b *RootContainer) Render() app.UI {
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
	var creds Creds
	ctx.GetState("creds", &creds)
	ctx.Async(func() {
		if creds.LoggedIn {
			log.Println("Getting new API response...")
			m, err := GetFlushesFromApi(ctx)
			if err != nil {
				
			} else {
				app.Window().GetElementByID("fetched-flushes").Set("innerHTML", m)
			}
		}})
}

type LoginContainer struct {
	app.Compo
}
func (l *LoginContainer) Render() app.UI {
	return app.Div().Body(app.Div().Body(
		app.P().Text("Log in to continue.").Class("font-bold"),
		app.Div().Body(
			app.Input().Type("text").ID("username").Placeholder("Username").Class(
				"m-2 placeholder-gray-500",
			),
			app.Br(),
			app.Input().Type("password").ID("password").Placeholder("Password").Class(
				"m-2 placeholder-gray-500",
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
	).Class("p-4 text-center text-xl shadow-lg bg-white rounded-lg",
	).ID("login-container"),
	&RegisterContainer{},
	app.Div().Body(&ErrorContainer{})).Class(
		CenteringDivCss)
}

type buttonLogin struct {
	app.Compo
}
func (b *buttonLogin) Render() app.UI {
	return app.Button().Text("Log in").OnClick(b.onClick).Class(
		YellowButtonCss)
}
func (b *buttonLogin) onClick(ctx app.Context, e app.Event) {
	loginSeconds := 60
	log.Println("Trying to log in...")
	if app.Window().GetElementByID("remember-me").Get("checked").Bool() {
		log.Println("remember-me checked")
		loginSeconds = 604800 // week
	}
	lastCreds := LastTriedCreds{}
	ctx.GetState("lastUsedCreds", &lastCreds)
	user, pass := GetLoginCreds()
	err := ValidateLoginCreds(user, pass, lastCreds)
	if err != nil {
		ShowErrorDiv(ctx, err, 1)
		return
	}
	ctx.Async(func() {
		status, basic_auth, err := TryLogin(user, pass)
		if err != nil {
			ShowErrorDiv(ctx, err, 1)
		}
		switch status {
		case 200:
			ctx.SetState("creds", Creds{
				UserColonPass: basic_auth,
				LoggedIn:      true,
			}).ExpiresIn(time.Second * time.Duration(loginSeconds)).PersistWithEncryption()
			log.Printf("Logged in as %s\n for %v seconds!", user, loginSeconds)
			app.Window().Set("location", ".")
			ctx.DelState("lastUsedCreds")
		case 401:
			ShowErrorDiv(ctx, errors.New("invalid credentials"), 1)
			ctx.SetState("lastUsedCreds", LastTriedCreds{
				User: user,
				Password: pass,
			}).ExpiresIn(time.Second * 10)
		default:
			ShowErrorDiv(ctx, errors.New("login failed"), 1)
			ctx.DelState("lastUsedCreds")
		}
})
}

type buttonRegister struct {
	app.Compo
}
func (b *buttonRegister) Render() app.UI {
	return app.Button().Text("Register").OnClick(b.onClick).Class(YellowButtonCss)
}
func (b *buttonRegister) onClick(ctx app.Context, e app.Event) {
	log.Println("Trying to register...")
	lastCreds := LastTriedCreds{}
	ctx.GetState("lastUsedCredsRegister", &lastCreds)
	user, pass, repeatPass := GetRegisterCreds()
	err := ValidateRegistryCreds(user, pass, repeatPass, lastCreds)
	if err != nil {
		ShowErrorDiv(ctx, err, 1)
		return
	}
	ctx.Async(func() {
		status, basic_auth, err := TryRegister(user, pass)
		if err != nil {
			ShowErrorDiv(ctx, err, 1)
		}
		switch status {
		case 201:
			ctx.SetState("creds", Creds{
				UserColonPass: basic_auth,
				LoggedIn:      true,
			}).ExpiresIn(time.Second * time.Duration(604800)).PersistWithEncryption()
			ctx.DelState("lastUsedCredsRegister")
			app.Window().Set("location", ".")
		case 422:
			ShowBadRegisterCredsErr()
			SetLastUsedCredsState(ctx, user, pass)
		case 409:
			ShowErrorDiv(ctx, errors.New("username already exists"), 1)
			SetLastUsedCredsState(ctx, user, pass)
		default:
			ShowErrorDiv(ctx, errors.New("register failed"), 1)
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
	ctx.SetState("creds", Creds{LoggedIn: false}).PersistWithEncryption()
	app.Window().Set("location", ".")
}
