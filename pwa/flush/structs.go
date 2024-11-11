package flush

import "time"

type Creds struct {
	UserColonPass string
	LoggedIn      bool
}

type LastTriedCreds struct {
	User     string
	Password string
}

type Flush struct {
	ID        string    `json:"_id"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Rating    int       `json:"rating"`
	PhoneUsed bool      `json:"phone_used"`
	Note      string    `json:"note"`
}
