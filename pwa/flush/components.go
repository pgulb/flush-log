package flush

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	YellowButtonCss  = "font-bold bg-yellow-500 p-2 rounded text-white mx-1"
	ErrorDivCss      = "flex flex-row fixed bottom-4 left-4 bg-red-500 text-white p-4 text-xl rounded-lg"
	CenteringDivCss  = "flex flex-row min-h-screen justify-center items-center"
	RegisterDivCss   = "p-4 text-center text-xl shadow-lg bg-white rounded-lg mx-10"
	InviCss          = "fixed invisible"
	RootContainerCss = "shadow-lg bg-white rounded-lg p-6 min-h-72 relative"
	LoadingCss       = "flex flex-row justify-center items-center"
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
	return app.Button().Text("I need account").OnClick(b.onClick).Class(
		YellowButtonCss + " hover:bg-yellow-700").ID("show-register")
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
		&LoadingWidget{id: "register-loading"},
	)
}

type RootContainer struct {
	app.Compo
	buttonUpdate
	FlushList app.UI
}

func (b *RootContainer) OnMount(ctx app.Context) {
	b.buttonUpdate.RootContainer = b
	var creds Creds
	ctx.GetState("creds", &creds)
	log.Println("Logged in: ", creds.LoggedIn)
	if !creds.LoggedIn {
		log.Println("Not logged in at root...")
	} else {
		app.Window().GetElementByID("root-container").Set("className", RootContainerCss)
		app.Window().GetElementByID("about-container").Set("className", "invisible fixed")
		ShowLoading("flushes-loading")
		defer Hide("flushes-loading")
		flushes, err := GetFlushes(ctx)
		if err != nil {
			ShowErrorDiv(ctx, err, 1)
		} else {
			app.Window().GetElementByID("hidden-hello").Set("innerHTML", "hello!")
			b.FlushList = FLushTable(flushes)
		}
	}
}
func (b *RootContainer) Render() app.UI {
	return app.Div().Body(
		app.P().Text("empty").Class("invisible fixed").ID("hidden-hello"),
		app.Div().Body(
			app.H1().Text("Flush Log").Class("text-2xl"),
			&buttonLogout{},
			app.P().Text("Tracked flushes:").Class("py-2"),
			&LoadingWidget{id: "flushes-loading"},
			b.FlushList,
			app.Div().Body(
				b.buttonUpdate.Render(),
				&LinkButton{
					Text:          "(+)",
					Location:      "new",
					AdditionalCss: "absolute bottom-4 right-4 hover:bg-yellow-700",
				},
			).
				Class("m-10"),
		).Class("invisible fixed").ID("root-container"),
		&AboutContainer{},
		app.Div().Body(&ErrorContainer{}),
	)
}

type buttonUpdate struct {
	app.Compo
	*RootContainer
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
			ShowLoading("flushes-loading")
			defer Hide("flushes-loading")
			flushes, err := GetFlushes(ctx)
			if err != nil {
				ShowErrorDiv(ctx, err, 1)
			} else {
				app.Window().GetElementByID("hidden-hello").Set("innerHTML", "hello!")
				b.RootContainer.FlushList = FLushTable(flushes)
			}
		}
	})
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
			&LoadingWidget{id: "login-loading"},
		),
	).Class("p-4 text-center text-xl shadow-lg bg-white rounded-lg").ID("login-container"),
		&RegisterContainer{},
		app.Div().Body(&ErrorContainer{})).Class(
		CenteringDivCss)
}

type buttonLogin struct {
	app.Compo
}

func (b *buttonLogin) Render() app.UI {
	return app.Button().Text("Log in").OnClick(b.onClick).Class(
		YellowButtonCss + " hover:bg-yellow-700")
}
func (b *buttonLogin) onClick(ctx app.Context, e app.Event) {
	loginSeconds := 60
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
		defer Hide("login-loading")
		status, basic_auth, err := TryLogin(user, pass)
		if err != nil {
			ShowErrorDiv(ctx, err, 1)
			return
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
				User:     user,
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
	return app.Button().Text("Register").OnClick(b.onClick).Class(
		YellowButtonCss + " hover:bg-yellow-700").ID("register-button")
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
		defer Hide("register-loading")
		status, basic_auth, err := TryRegister(user, pass)
		log.Println("register status code: ", status)
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
		"font-bold border-2 border-white p-2 rounded absolute top-4 right-4")
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
			).Class("p-4 text-center text-xl shadow-lg bg-white rounded-lg"),
			app.Br(),
			&LinkButton{Text: "Back to Home Screen", Location: "."},
		).Class("flex flex-col"),
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
		defer Hide("new-flush-loading")
		statusCode, err := TryAddFlush(creds, flush)
		log.Println("Flush add statusCode: ", statusCode)
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
}

type AboutContainer struct {
	app.Compo
}

func (a *AboutContainer) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.H1().Text("About Flush Log").Class("text-2xl bold"),
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
				AdditionalCss: "hover:bg-yellow-700",
			},
		).Class("flex flex-col p-4 shadow-lg rounded-lg"),
	).Class(CenteringDivCss).ID("about-container")
}

func FLushTable(flushes []Flush) app.UI {
	if len(flushes) == 0 {
		return app.Div().Body(app.P().Text("No flushes yet."))
	}
	divs := []app.UI{}
	var totalTime time.Duration
	var count int
	var meanRating int
	var timesWithPhone int
	for _, flush := range flushes {
		totalTime += flush.TimeEnd.Sub(flush.TimeStart)
		count++
		if flush.PhoneUsed {
			timesWithPhone++
		}
		meanRating += flush.Rating
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
				app.P().Text("Phone used: "+phoneUsed),
				app.P().Text("Note: '"+flush.Note+"'"),
			).Class("flex flex-col p-4 border-1 shadow-lg rounded-lg"),
		)
	}
	statsDiv := app.Div().Body(
		app.P().Text("Total flushes: "+strconv.Itoa(count)),
		app.P().Text("Total time: "+strconv.Itoa(int(totalTime.Minutes()))+" minutes"),
		app.P().Text("Mean time: "+strconv.Itoa(int(totalTime.Minutes())/count)+" minutes"),
		app.P().Text("Mean rating: "+strconv.Itoa(meanRating/count)),
		app.P().Text("Times with phone used: "+strconv.Itoa(timesWithPhone)),
		app.P().Text("Percent with phone used: "+strconv.Itoa(timesWithPhone*100/count)+"%"),
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

type LoadingWidget struct {
	app.Compo
	id string
}

func (l *LoadingWidget) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.Span().
				Text("Loading...").Class("font-bold text-black").
				Class("!absolute !-m-px !h-px !w-px !overflow-hidden !whitespace-nowrap !border-0 !p-0 ![clip:rect(0,0,0,0)]"),
		).
			Class("inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-e-transparent align-[-0.125em] text-surface motion-reduce:animate-[spin_1.5s_linear_infinite] text-yellow-500"),
	).Class(InviCss).ID(l.id)
}
