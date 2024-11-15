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

type FlushStats struct {
	FlushCount       float64 `json:"flushCount"`
	TotalTime        float64 `json:"totalTime"`
	MeanTime         float64 `json:"meanTime"`
	MeanRating       float64 `json:"meanRating"`
	PhoneUsedCount   float64 `json:"phoneUsedCount"`
	PercentPhoneUsed float64 `json:"percentPhoneUsed"`
}

type FlushStatsInt struct {
	FlushCount       int
	TotalTime        int
	MeanTime         int
	MeanRating       int
	PhoneUsedCount   int
	PercentPhoneUsed int
}
