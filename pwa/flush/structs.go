package flush

type Creds struct {
	UserColonPass string
	LoggedIn      bool
}

type LastTriedCreds struct {
	User     string
	Password string
}
