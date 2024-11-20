package flush

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	YellowButtonCss  = "font-bold bg-amber-600 p-2 rounded mx-1"
	ErrorDivCss      = "flex flex-row fixed bottom-4 left-4 bg-red-500 p-4 text-xl rounded-lg"
	CenteringDivCss  = "flex flex-row min-h-screen justify-center items-center"
	WindowDivCss     = "p-4 text-center text-xl shadow-lg bg-zinc-800 rounded-lg mx-10"
	InviCss          = "fixed invisible"
	RootContainerCss = "shadow-lg bg-zinc-800 rounded-lg p-6 min-h-72 relative"
	LoadingCss       = "flex flex-row justify-center items-center"
	RemoveButtonCss  = "font-bold bg-red-500 p-2 rounded hover:bg-red-700 m-1"
	LogoutButtonCss  = "font-bold bg-amber-700 p-2 rounded mx-1 hover:bg-amber-900"
	UpdateButtonCss  = "bg-green-600 hover:bg-green-800 text-xl p-2 rounded bottom-4 right-4 fixed"
	InstallButtonCss = "bg-green-600 hover:bg-green-800 p-2 rounded m-2"
)

type ErrorContainer struct {
	app.Compo
}

func (e *ErrorContainer) Render() app.UI {
	return app.Div().Body(app.Div().Body(
		app.P().Text("placeholder error")).Class(
		"p-8 text-center text-xl shadow-lg bg-zinc-800 rounded-lg",
	)).Class(
		InviCss,
	).ID("error")
}

type buttonShowRegister struct {
	app.Compo
}

func (b *buttonShowRegister) Render() app.UI {
	return app.Button().Text("I need account").OnClick(b.onClick).Class(
		YellowButtonCss + " hover:bg-amber-800").ID("show-register")
}
func (b *buttonShowRegister) onClick(ctx app.Context, e app.Event) {
	app.Window().GetElementByID("register-container").Set("className", WindowDivCss)
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
				"m-2",
			),
			app.Br(),
			app.Input().Type("password").ID("register-password").Placeholder("Password").Class(
				"m-2",
			),
			app.Br(),
			app.Input().Type("password").ID("register-password-repeat").Placeholder(
				"Repeat password").Class(
				"m-2 my-4",
			),
			app.Br(),
			app.P().Text("No way to reset password"),
			app.P().Text("if you forget it."),
			app.Br(),
			&buttonRegister{},
			app.P().Text("").Class("text-red-500").ID("register-error"),
		).Class(InviCss).ID("register-container"),
		&LoadingWidget{id: "register-loading"},
	)
}

type RootContainer struct {
	app.Compo
	buttonUpdate
	FlushList app.UI
}

func (b *RootContainer) OnMount(ctx app.Context) {
	b.buttonUpdate.parent = b
	var creds Creds
	ctx.GetState("creds", &creds)
	log.Println("Logged in: ", creds.LoggedIn)
	if !creds.LoggedIn {
		log.Println("Not logged in at root...")
	} else {
		Hide("update-button")
		app.Window().GetElementByID("root-container").Set("className", RootContainerCss)
		app.Window().GetElementByID("about-container").Set("className", "invisible fixed")
		ShowLoading("flushes-loading")
		ctx.Async(func() {
			result := GetFlushesFromOID(ctx)
			ctx.Dispatch(func(ctx app.Context) {
				defer Hide("flushes-loading")
				if result == nil {
					ShowErrorDiv(ctx, errors.New("Error while fetching flushes"), 2)
					return
				}
				app.Window().GetElementByID("hidden-hello").Set("innerHTML", "hello!")
				b.SetList(result)
				var isEmpty string
				ctx.GetState("no-flushes", &isEmpty)
				log.Println("no-flushes: ", isEmpty)
				if isEmpty != "true" {
					log.Println("Viewing update button")
					app.Window().GetElementByID("update-button").
						Set("className", YellowButtonCss+" hover:bg-amber-800 align-middle")
				} else {
					log.Println("No flushes, hiding update button")
					app.Window().GetElementByID("update-button").Set("className", InviCss)
				}
			})
		})
	}
}
func (b *RootContainer) Render() app.UI {
	return app.Div().Body(
		app.P().Text("empty").Class("invisible fixed").ID("hidden-hello"),
		&UpdateButton{},
		app.Div().Body(
			app.H1().Text("Flush Log").Class("text-2xl"),
			app.Div().Body(
				&buttonLogout{},
				&LinkButton{
					Text:          "Settings",
					Location:      "settings",
					AdditionalCss: "hover:bg-amber-800 m-1",
				},
				&LinkButton{
					Text:          "(+)",
					Location:      "new",
					AdditionalCss: "hover:bg-amber-800",
				},
				&InstallButton{},
			).ID("root-buttons-container").Class("flex flex-col absolute top-4 right-4"),
			app.P().Text("Tracked flushes:").Class("py-2"),
			&LoadingWidget{id: "flushes-loading"},
			b.FlushList,
			app.Div().Body(
				&b.buttonUpdate,
				&LoadingWidget{id: "flushes-loading-update"},
			).
				Class(LoadingCss+" m-2"),
		).Class("invisible fixed").ID("root-container"),
		&AboutContainer{},
		app.Div().Body(&ErrorContainer{}),
	)
}
func (b *RootContainer) SetList(list app.UI) {
	b.FlushList = list
}

type buttonUpdate struct {
	app.Compo
	parent *RootContainer
}

func (b *buttonUpdate) Render() app.UI {
	return app.Button().Text("Load More").OnClick(b.onClick).Class(
		YellowButtonCss + " hover:bg-amber-800 align-middle").ID("update-button")
}
func (b *buttonUpdate) onClick(ctx app.Context, e app.Event) {
	ShowLoading("flushes-loading-update")
	Hide("update-button")
	ctx.Async(func() {
		result := GetFlushesFromOID(ctx)
		ctx.Dispatch(func(ctx app.Context) {
			defer Hide("flushes-loading-update")
			if result == nil {
				ShowErrorDiv(ctx, errors.New("Error while fetching flushes"), 2)
				return
			}
			b.parent.SetList(result)
			var isEmpty string
			ctx.GetState("no-flushes", &isEmpty)
			if isEmpty == "true" {
				log.Println("No flushes, hiding update button")
				app.Window().GetElementByID("update-button").Set("className", InviCss)
				b.parent.SetList(app.Div().
					Body(app.P().Text("No more flushes to show.").Class("py-2")))
			} else {
				app.Window().GetElementByID("update-button").
					Set("className", YellowButtonCss+" hover:bg-amber-800 align-middle")
			}
		})
	})
}

type LoginContainer struct {
	app.Compo
}

func (l *LoginContainer) Render() app.UI {
	return app.Div().Body(
		&UpdateButton{},
		app.Div().Body(
			app.P().Text("Log in to continue.").Class("font-bold"),
			app.Div().Body(
				app.Input().Type("text").ID("username").Placeholder("Username").Class(
					"m-2",
				),
				app.Br(),
				app.Input().Type("password").ID("password").Placeholder("Password").Class(
					"m-2",
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
				&LoadingWidget{id: "login-loading"},
			),
		).Class("p-4 text-center text-xl shadow-lg bg-zinc-800 rounded-lg").ID("login-container"),
		&RegisterContainer{},
		app.Div().Body(&ErrorContainer{})).Class(
		CenteringDivCss)
}

type buttonLogin struct {
	app.Compo
}

func (b *buttonLogin) Render() app.UI {
	return app.Button().Text("Log in").OnClick(b.onClick).Class(
		YellowButtonCss + " hover:bg-amber-800").ID("login-button")
}
func (b *buttonLogin) onClick(ctx app.Context, e app.Event) {
	loginSeconds := 600
	log.Println("Trying to log in...")
	ShowLoading("login-loading")
	if app.Window().GetElementByID("remember-me").Get("checked").Bool() {
		log.Println("remember-me checked")
		loginSeconds = 604800 // week
	}
	lastCreds := LastTriedCreds{}
	ctx.GetState("lastUsedCreds", &lastCreds)
	user, pass := GetLoginCreds()
	err := ValidateLoginCreds(user, pass, lastCreds)
	if err != nil {
		Hide("login-loading")
		ShowErrorDiv(ctx, err, 1)
		return
	}
	ctx.Async(func() {
		status, basic_auth, err := TryLogin(user, pass)
		if err != nil {
			ShowErrorDiv(ctx, err, 1)
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			defer Hide("login-loading")
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
					User:     user,
					Password: pass,
				}).ExpiresIn(time.Second * 10)
			default:
				ShowErrorDiv(ctx, errors.New("login failed"), 1)
				ctx.DelState("lastUsedCreds")
			}
		})
	})
}

type buttonRegister struct {
	app.Compo
}

func (b *buttonRegister) Render() app.UI {
	return app.Button().Text("Register").OnClick(b.onClick).Class(
		YellowButtonCss + " hover:bg-amber-800").ID("register-button")
}
func (b *buttonRegister) onClick(ctx app.Context, e app.Event) {
	log.Println("Trying to register...")
	ShowLoading("register-loading")
	lastCreds := LastTriedCreds{}
	ctx.GetState("lastUsedCredsRegister", &lastCreds)
	user, pass, repeatPass := GetRegisterCreds()
	err := ValidateRegistryCreds(user, pass, repeatPass, lastCreds)
	if err != nil {
		Hide("register-loading")
		ShowErrorDiv(ctx, err, 1)
		return
	}
	ctx.Async(func() {
		status, basic_auth, err := TryRegister(user, pass)
		log.Println("register status code: ", status)
		if err != nil {
			ShowErrorDiv(ctx, err, 1)
		}
		ctx.Dispatch(func(ctx app.Context) {
			defer Hide("register-loading")
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
	})
}

type buttonLogout struct {
	app.Compo
}

func (b *buttonLogout) Render() app.UI {
	return app.Button().Text("Log out").OnClick(b.onClick).Class(
		LogoutButtonCss)
}
func (b *buttonLogout) onClick(ctx app.Context, e app.Event) {
	ctx.SetState("creds", Creds{LoggedIn: false}).PersistWithEncryption()
	app.Window().Set("location", ".")
}

type NewFlushContainer struct {
	app.Compo
}

func (c *NewFlushContainer) Render() app.UI {
	return app.Div().Body(
		&UpdateButton{},
		app.Div().Body(
			app.Div().Body(
				app.P().Text("Add new flush").Class("font-bold"),
				app.Br(),
				app.Label().For("new-flush-time-start").Text("Start:").Class("m-2"),
				app.Input().Type("datetime-local").ID("new-flush-time-start").Class("m-2"),
				app.Br(),
				app.Label().For("new-flush-time-end").Text("End:").Class("m-2"),
				app.Input().Type("datetime-local").ID("new-flush-time-end").Class("m-2"),
				app.Br(),
				app.Label().For("new-flush-rating").Text("Rating (1-worst, 10-best)").Class("m-2"),
				app.Select().ID("new-flush-rating").Class("m-2").Body(
					app.Option().Value("1").Text("1"),
					app.Option().Value("2").Text("2"),
					app.Option().Value("3").Text("3"),
					app.Option().Value("4").Text("4"),
					app.Option().Value("5").Text("5"),
					app.Option().Value("6").Text("6"),
					app.Option().Value("7").Text("7"),
					app.Option().Value("8").Text("8"),
					app.Option().Value("9").Text("9"),
					app.Option().Value("10").Text("10"),
				),
				app.Br(),
				app.Label().For("new-flush-phone-used").Text("Phone used").Class("m-2"),
				app.Input().Type("checkbox").ID("new-flush-phone-used").Class("m-2"),
				app.Br(),
				app.Hr(),
				app.Textarea().Placeholder("notes").ID(
					"new-flush-note").MaxLength(100),
				app.Br(),
				&SubmitFlushButton{},
				&LoadingWidget{id: "new-flush-loading"},
			).Class("p-4 text-center text-xl shadow-lg bg-zinc-800 rounded-lg"),
			app.Br(),
			&LinkButton{
				Text:          "Back to Home Screen",
				Location:      ".",
				AdditionalCss: "hover:bg-amber-800",
			},
		).
			Class("flex flex-col"),
		app.Div().Body(&ErrorContainer{}),
	).Class(CenteringDivCss)
}
func (c *NewFlushContainer) OnMount(ctx app.Context) {
	var creds Creds
	ctx.GetState("creds", &creds)
	log.Println("Logged in: ", creds.LoggedIn)
	if !creds.LoggedIn {
		app.Window().Set("location", "login")
		return
	}
	var set string
	ctx.GetState("phoneUsedDefault", &set)
	if set == "true" {
		app.Window().GetElementByID("new-flush-phone-used").Set("checked", true)
	}
}

type LinkButton struct {
	app.Compo
	Text          string
	Location      string
	AdditionalCss string
}

func (b *LinkButton) Render() app.UI {
	if b.AdditionalCss != "" {
		return app.Button().
			Text(b.Text).
			Class(b.AdditionalCss + " " + YellowButtonCss).
			OnClick(b.onClick)
	}
	return app.Button().Text(b.Text).Class(YellowButtonCss).OnClick(b.onClick)
}
func (b *LinkButton) onClick(ctx app.Context, e app.Event) {
	app.Window().Set("location", b.Location)
}

type SubmitFlushButton struct {
	app.Compo
}

func (b *SubmitFlushButton) Render() app.UI {
	return app.Button().
		Text("Submit").
		Class(YellowButtonCss).
		ID("submit-flush-button").
		OnClick(b.onClick)
}
func (b *SubmitFlushButton) onClick(ctx app.Context, e app.Event) {
	var creds Creds
	ShowLoading("new-flush-loading")
	ctx.GetState("creds", &creds)
	log.Println("Logged in: ", creds.LoggedIn)
	if !creds.LoggedIn {
		Hide("new-flush-loading")
		app.Window().Set("location", "login")
		return
	}
	flush, err := NewFLush(ctx,
		app.Window().GetElementByID("new-flush-time-start").Get("value").String(),
		app.Window().GetElementByID("new-flush-time-end").Get("value").String(),
		app.Window().GetElementByID("new-flush-rating").Get("value").String(),
		app.Window().GetElementByID("new-flush-phone-used").Get("checked").Bool(),
		app.Window().GetElementByID("new-flush-note").Get("value").String())
	if err != nil {
		Hide("new-flush-loading")
		ShowErrorDiv(ctx, err, 2)
		return
	}
	err = ValidateFlush(flush)
	if err != nil {
		Hide("new-flush-loading")
		ShowErrorDiv(ctx, err, 1)
		return
	}
	ctx.Async(func() {
		statusCode, err := TryAddFlush(creds, flush)
		log.Println("Flush add statusCode: ", statusCode)
		ctx.Dispatch(func(ctx app.Context) {
			defer Hide("new-flush-loading")
			if err != nil {
				ShowErrorDiv(ctx, err, 2)
				return
			}
			switch statusCode {
			case 201, 204:
				app.Window().Set("location", ".")
			default:
				ShowErrorDiv(ctx, errors.New("Unexpected error while adding flush"), 2)
			}
		})
	})
}

type AboutContainer struct {
	app.Compo
}

func (a *AboutContainer) Render() app.UI {
	return app.Div().Body(
		&UpdateButton{},
		app.Div().Body(
			app.Div().Body(
				app.H1().Text("About Flush Log").Class("text-2xl bold inline"),
				app.Img().Src("/web/paper.png").Class("inline").Width(100).Height(100),
			),
			app.Br(),
			app.P().Text("Application for tracking toilet flushes."),
			app.P().Text("You can use it to save them, rate them, check your toilet stats."),
			app.P().
				Text("The app will show you statistics like total time spent, average time spent, % times with phone used etc."),
			app.Br(),
			app.P().Text("App is still under development. New features can be added."),
			app.P().
				Text("App can be 'installed' - it will appear on computer's program list or on phone home screen."),
			app.Br(),
			&LinkButton{
				Text:          "Login/Register",
				Location:      "login",
				AdditionalCss: "hover:bg-amber-800",
			},
			&InstallButton{},
		).Class("flex flex-col p-4 shadow-lg rounded-lg bg-zinc-800"),
	).Class(CenteringDivCss).ID("about-container")
}

func FlushTable(flushes []Flush, ctx app.Context) app.UI {
	if len(flushes) == 0 {
		ctx.SetState("no-flushes", "true")
		return app.Div().Body(app.P().Text("No flushes yet."))
	}
	divs := []app.UI{}
	for _, flush := range flushes {
		var phoneUsed string
		if flush.PhoneUsed {
			phoneUsed = "Yes"
		} else {
			phoneUsed = "No"
		}
		divs = append(divs,
			app.Div().Body(
				timeDiv(flush),
				app.P().Text("Rating: "+strconv.Itoa(flush.Rating)),
				app.Div().Body(
					&RemoveFlushButton{ID: flush.ID},
					&ConfirmRemoveFlushButton{ID: flush.ID},
					&CancelRemoveFlushButton{ID: flush.ID},
				).Class("max-w-1/6 remove-flush-buttonz-div"),
				app.P().Text("Phone used: "+phoneUsed),
				app.P().Text("Note: '"+flush.Note+"'").Class("break-all"),
			).Class("flex flex-col p-4 border-1 shadow-lg rounded-lg").ID("div-"+flush.ID),
		)
	}
	stats, err := GetStats(ctx)
	if err != nil {
		ShowErrorDiv(ctx, err, 3)
		return app.Div()
	}
	statsDiv := app.Div().Body(
		app.P().Text("Total flushes: "+strconv.Itoa(stats.FlushCount)),
		app.P().Text("Total time: "+strconv.Itoa(stats.TotalTime)+" minutes"),
		app.P().Text("Mean time: "+strconv.Itoa(stats.MeanTime)+" minutes"),
		app.P().Text("Mean rating: "+strconv.Itoa(stats.MeanRating)),
		app.P().Text("Times with phone used: "+strconv.Itoa(stats.PhoneUsedCount)),
		app.P().Text("Percent with phone used: "+strconv.Itoa(stats.PercentPhoneUsed)+"%"),
	).Class("flex flex-col p-4 border-1 shadow-lg rounded-lg font-bold")
	divs = append([]app.UI{statsDiv}, divs...)
	return app.Div().Body(divs...)
}

func timeDiv(flush Flush) app.UI {
	flushDuration := strconv.FormatFloat(
		flush.TimeEnd.Sub(flush.TimeStart).Minutes(),
		'f',
		0,
		64,
	)
	if flush.TimeStart.Day() == flush.TimeEnd.Day() {
		return app.Div().Body(
			app.P().Text("Time: ").Class("font-bold inline"),
			app.P().Text(flushDuration+" minutes, "+flush.TimeStart.Format(
				"2006-01-02 15:04")+"-"+flush.TimeEnd.Format("15:04")).Class("inline"),
		)
	} else {
		return app.Div().Body(
			app.P().Text("Time: ").Class("font-bold inline"),
			app.P().Text(flushDuration+" minutes, "+flush.TimeStart.Format(
				"2006-01-02 15:04")+" - "+flush.TimeEnd.Format("2006-01-02 15:04")).Class("inline"),
		)
	}
}

type RemoveFlushButton struct {
	app.Compo
	ID string
}

func (b *RemoveFlushButton) Render() app.UI {
	return app.Button().Text("🗑️").Class(RemoveButtonCss).ID(b.ID).OnClick(b.onClick)
}
func (b *RemoveFlushButton) onClick(ctx app.Context, e app.Event) {
	log.Printf("Flush remove button pressed (%s)...\n", b.ID)
	app.Window().GetElementByID(b.ID).Set("className", InviCss)
	app.Window().GetElementByID(b.ID+"-confirm").Set("className", RemoveButtonCss)
	app.Window().GetElementByID(b.ID+"-cancel").Set("className", YellowButtonCss)
}

type ConfirmRemoveFlushButton struct {
	app.Compo
	ID string
}

func (b *ConfirmRemoveFlushButton) Render() app.UI {
	return app.Button().
		Text("CONFIRM DELETE").
		Class(InviCss).
		ID(b.ID + "-confirm").
		OnClick(b.onClick)
}
func (b *ConfirmRemoveFlushButton) onClick(ctx app.Context, e app.Event) {
	log.Printf("Confirm remove button pressed (%s)...\n", b.ID)
	log.Println("removing flush " + b.ID + "...")
	ShowLoading("flushes-loading")
	ctx.Async(func() {
		var creds Creds
		ctx.GetState("creds", &creds)
		err := RemoveFlush(b.ID, creds.UserColonPass)
		ctx.Dispatch(func(ctx app.Context) {
			defer Hide("flushes-loading")
			if err != nil {
				ShowErrorDiv(ctx, err, 2)
				return
			}
			Hide("div-" + b.ID)
		})
	})
}

type CancelRemoveFlushButton struct {
	app.Compo
	ID string
}

func (b *CancelRemoveFlushButton) Render() app.UI {
	return app.Button().
		Text("CANCEL").
		Class(InviCss).
		ID(b.ID + "-cancel").
		OnClick(b.onClick)
}
func (b *CancelRemoveFlushButton) onClick(ctx app.Context, e app.Event) {
	log.Printf("Cancel remove button pressed (%s)...\n", b.ID)
	app.Window().GetElementByID(b.ID).Set("className", RemoveButtonCss)
	app.Window().GetElementByID(b.ID+"-confirm").Set("className", InviCss)
	app.Window().GetElementByID(b.ID+"-cancel").Set("className", InviCss)
}

type LoadingWidget struct {
	app.Compo
	id string
}

func (l *LoadingWidget) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.Span().
				Text("Loading...").Class("font-bold").
				Class("!absolute !-m-px !h-px !w-px !overflow-hidden !whitespace-nowrap !border-0 !p-0 ![clip:rect(0,0,0,0)]"),
		).
			Class("inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-e-transparent align-[-0.125em] text-surface motion-reduce:animate-[spin_1.5s_linear_infinite] text-amber-600"),
	).Class(InviCss).ID(l.id)
}

type SettingsContainer struct {
	app.Compo
}

func (s *SettingsContainer) Render() app.UI {
	return app.Div().Body(
		&UpdateButton{},
		app.Div().Body(
			app.Div().Body(
				app.H1().Text("App Settings").Class("text-2xl m-2"),
				&PassChangeContainer{},
				app.Br(),
				app.Hr(),
				app.Br(),
				app.P().
					Text("You can export your flushes into formats that can be read by other apps").
					Class("m-1"),
				&ExportButton{ExportFormat: "JSON"},
				app.Br(),
				&ExportButton{ExportFormat: "CSV"},
				app.Br(),
				app.Hr(),
				app.Br(),
				app.P().
					Text("Below settings are stored in your browser only").
					Class("font-bold m-2"),
				app.Label().
					Text("Check 'phone used' option by default").
					For("phone-used-default").Class("m-2"),
				&PhoneUsedDefaultCheckbox{},
				app.Br(),
				app.Hr(),
				app.Br(),
				&RemoveAccountContainer{},
			).
				Class(WindowDivCss),
			app.Br(),
			&LinkButton{
				Text:          "Back to Home Screen",
				Location:      ".",
				AdditionalCss: "hover:bg-amber-800",
			},
			&LoadingWidget{id: "settings-loading"},
		).
			ID("settings-container").Class(CenteringDivCss+" flex-col"),
		&ErrorContainer{},
	)
}
func (s *SettingsContainer) OnMount(ctx app.Context) {
	var creds Creds
	ctx.GetState("creds", &creds)
	log.Println("Logged in: ", creds.LoggedIn)
	if !creds.LoggedIn {
		app.Window().Set("location", "login")
		return
	}
}

type PhoneUsedDefaultCheckbox struct {
	app.Compo
	Check app.UI
}

func (c *PhoneUsedDefaultCheckbox) Render() app.UI {
	c.Check = app.Input().
		Type("checkbox").
		ID("phone-used-default").
		Class("m-2").
		OnClick(c.OnClick)
	return c.Check
}
func (c *PhoneUsedDefaultCheckbox) OnClick(ctx app.Context, e app.Event) {
	var set string
	log.Println("Getting phoneUsedDefault...")
	ctx.GetState("phoneUsedDefault", &set)
	if set == "true" {
		log.Println("Setting phoneUsedDefault to false")
		ctx.SetState("phoneUsedDefault", "false").Persist()
	} else {
		log.Println("Setting phoneUsedDefault to true")
		ctx.SetState("phoneUsedDefault", "true").Persist()
	}
}
func (c *PhoneUsedDefaultCheckbox) OnMount(ctx app.Context) {
	var set string
	log.Println("Getting phoneUsedDefault in OnLoad...")
	ctx.GetState("phoneUsedDefault", &set)
	if set == "true" {
		app.Window().GetElementByID("phone-used-default").Set("checked", true)
	}
}

type ExportButton struct {
	app.Compo
	ExportFormat string
}

func (b *ExportButton) Render() app.UI {
	return app.Button().
		Text(fmt.Sprintf("Export to %s", b.ExportFormat)).
		Class(YellowButtonCss + " hover:bg-amber-800 m-1").
		OnClick(b.OnClick)
}
func (b *ExportButton) OnClick(ctx app.Context, e app.Event) {
	var creds Creds
	ctx.GetState("creds", &creds)
	decoded, err := base64.StdEncoding.DecodeString(creds.UserColonPass)
	if err != nil {
		ShowErrorDiv(ctx, err, 1)
		return
	}
	apiUrl, err := GetApiUrl()
	if err != nil {
		ShowErrorDiv(ctx, err, 1)
		return
	}
	apiUrl = strings.Replace(apiUrl, "http://", fmt.Sprintf("http://%s@", decoded), 1)
	apiUrl = strings.Replace(apiUrl, "https://", fmt.Sprintf("https://%s@", decoded), 1)
	completeUrl := fmt.Sprintf(
		"%s/flushes?export_format=%s",
		apiUrl,
		strings.ToLower(b.ExportFormat),
	)
	app.Window().
		Call("open", completeUrl)
}

type PassChangeContainer struct {
	app.Compo
}

func (p *PassChangeContainer) Render() app.UI {
	return app.Div().Body(
		app.P().
			Text("Change password").
			Class("m-1"),
		app.P().
			Text("You will be prompted to log in again after changing").
			Class("m-1"),
		app.Input().Type("password").ID(
			"chp-password").Placeholder("New password").Class(
			"m-1",
		),
		app.Input().Type("password").ID(
			"chp-password-repeat").Placeholder("Repeat password").Class(
			"m-1",
		),
		&ChangePassButton{},
	).ID("passchange-container").Class("flex flex-col")
}

type ChangePassButton struct {
	app.Compo
}

func (c *ChangePassButton) Render() app.UI {
	return app.Button().
		Text("Change").
		Class(YellowButtonCss + " hover:bg-amber-800 m-1").
		OnClick(c.OnClick).ID("chp-button")
}
func (c *ChangePassButton) OnClick(ctx app.Context, e app.Event) {
	ShowLoading("settings-loading")
	newPass := app.Window().GetElementByID("chp-password").Get("value").String()
	repeatPass := app.Window().GetElementByID("chp-password-repeat").Get("value").String()
	var creds Creds
	ctx.GetState("creds", &creds)
	if err := ValidateChangePass(newPass, repeatPass); err != nil {
		Hide("settings-loading")
		ShowErrorDiv(ctx, err, 1)
		return
	}
	if err := ChangePass(newPass, creds.UserColonPass); err != nil {
		Hide("settings-loading")
		ShowErrorDiv(ctx, err, 1)
		return
	}
	ctx.SetState("creds", Creds{LoggedIn: false}).PersistWithEncryption()
	Hide("settings-loading")
	app.Window().Set("location", "login")
}

type RemoveAccountContainer struct {
	app.Compo
}

func (r *RemoveAccountContainer) Render() app.UI {
	return app.Div().Body(
		app.P().
			Text("Remove account").
			Class("m-1"),
		app.Input().Placeholder("Type 'byebye' here").ID("remove-account-byebye").Class("m-1"),
		&RemoveAccountButton{},
	).ID("remove-account-container").Class("flex flex-col")
}

type RemoveAccountButton struct {
	app.Compo
}

func (c *RemoveAccountButton) Render() app.UI {
	return app.Button().
		Text("Remove account").
		Class(RemoveButtonCss).
		OnClick(c.OnClick).ID("remove-account-button")
}
func (c *RemoveAccountButton) OnClick(ctx app.Context, e app.Event) {
	var creds Creds
	ShowLoading("settings-loading")
	ctx.GetState("creds", &creds)
	if app.Window().GetElementByID("remove-account-byebye").Get("value").String() != "byebye" {
		Hide("settings-loading")
		ShowErrorDiv(ctx, errors.New("Type 'byebye' before deleting account"), 1)
		return
	}
	if err := RemoveAccount(creds.UserColonPass); err != nil {
		ShowErrorDiv(ctx, err, 1)
		Hide("settings-loading")
		return
	}
	ctx.SetState("creds", Creds{LoggedIn: false}).PersistWithEncryption()
	app.Window().Set("location", ".")
}

func GetFlushesFromOID(ctx app.Context) app.UI {
	var skip int
	ctx.GetState("skip", &skip)
	fls, err := GetFlushes(ctx, skip)
	if err != nil {
		return nil
	}
	result := FlushTable(fls, ctx)
	return result
}

type UpdateButton struct {
	app.Compo
	updateAvailable bool
	Css             string
}

func (c *UpdateButton) OnAppUpdate(ctx app.Context) {
	c.updateAvailable = ctx.AppUpdateAvailable()
}
func (c *UpdateButton) Render() app.UI {
	if c.updateAvailable {
		c.Css = UpdateButtonCss
		log.Println("There is an update available!")
	} else {
		c.Css = InviCss
	}
	return app.Button().
		Text("Update App ⬇️").
		OnClick(c.onUpdateClick).
		Class(c.Css)
}
func (c *UpdateButton) onUpdateClick(ctx app.Context, e app.Event) {
	ctx.Reload()
}

type InstallButton struct {
	app.Compo
	name             string
	isAppInstallable bool
}

func (b *InstallButton) OnMount(ctx app.Context) {
	b.isAppInstallable = ctx.IsAppInstallable()
}
func (b *InstallButton) OnAppInstallChange(ctx app.Context) {
	b.isAppInstallable = ctx.IsAppInstallable()
}
func (b *InstallButton) Render() app.UI {
	return app.Div().
		Body(
			app.If(b.isAppInstallable, func() app.UI {
				return app.Button().
					Text("Install App").
					OnClick(b.onInstallButtonClicked).
					Class(InstallButtonCss)
			}),
		).Class("flex flex-col")
}
func (b *InstallButton) onInstallButtonClicked(ctx app.Context, e app.Event) {
	ctx.ShowAppInstallPrompt()
}
