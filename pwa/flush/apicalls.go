package flush

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func CloseBody(r *http.Response) {
	if err := r.Body.Close(); err != nil {
		DisplayError(err)
	}
}

func GetApiUrl() (string, error) {
	r, err := http.Get("web/apiurl")
	if err != nil {
		return "", err
	}
	defer CloseBody(r)
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func TryLogin(username string, password string) (int, string, error) {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return 0, "", err
	}
	basicAuth := BasicAuth(username, password)
	req, err := AuthorizedRequest("GET", apiUrl, basicAuth)
	if err != nil {
		return 0, "", err
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer CloseBody(r)
	return r.StatusCode, basicAuth, nil
}

func TryRegister(username string, password string) (int, string, error) {
	apiUrl, err := GetApiUrl()
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
	basicAuth := BasicAuth(username, password)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer CloseBody(r)
	return r.StatusCode, basicAuth, nil
}

func TryAuthentication(ctx app.Context) (string, error) {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return "", err
	}
	var c Creds
	ctx.GetState("creds", &c)
	req.Header.Add("Authorization", "Basic "+c.UserColonPass)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer CloseBody(r)
	if r.StatusCode >= 400 {
		ctx.SetState("creds", Creds{LoggedIn: false}).PersistWithEncryption()
		app.Window().Set("location", "login")
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func BasicAuth(username string, password string) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)
}

func AuthorizedRequest(method string, url string,
	basicAuth string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth)
	return req, nil
}

func TryAddFlush(creds Creds, flush Flush) (int, error) {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return 0, err
	}
	js, err := json.Marshal(flush)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("PUT", apiUrl+"/flush", bytes.NewBuffer(js))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Basic "+creds.UserColonPass)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer CloseBody(r)
	return r.StatusCode, nil
}

func GetFlushes(ctx app.Context, skip int) ([]Flush, bool, error) {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return nil, false, err
	}
	req, err := http.NewRequest("GET", apiUrl+"/flushes", nil)
	if err != nil {
		return nil, false, err
	}
	var c Creds
	ctx.GetState("creds", &c)
	req.Header.Add("Authorization", "Basic "+c.UserColonPass)
	log.Println("flush fetch skip: ", skip)
	log.Println("adding skip to /flushes...")
	q := url.Values{}
	q.Add("skip", strconv.Itoa(skip))
	req.URL.RawQuery = q.Encode()
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer CloseBody(r)
	if r.StatusCode >= 400 {
		ctx.SetState("creds", Creds{LoggedIn: false}).PersistWithEncryption()
		app.Window().Set("location", "login")
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, false, err
	}
	var temp TempFlushes
	err = json.Unmarshal(bytes, &temp)
	if err != nil {
		return nil, false, err
	}
	log.Println("more data available: ", temp.MoreDataAvailable)
	flushes := Flushes{
		Flushes:           make([]Flush, len(temp.Flushes)),
		MoreDataAvailable: temp.MoreDataAvailable,
	}
	for i := range temp.Flushes {
		flushes.Flushes[i].TimeStart, err = time.Parse(
			"2006-01-02T15:04:05",
			temp.Flushes[i].TimeStart,
		)
		if err != nil {
			return nil, false, err
		}
		flushes.Flushes[i].TimeEnd, err = time.Parse("2006-01-02T15:04:05", temp.Flushes[i].TimeEnd)
		if err != nil {
			return nil, false, err
		}
		flushes.Flushes[i].Rating = temp.Flushes[i].Rating
		flushes.Flushes[i].PhoneUsed = temp.Flushes[i].PhoneUsed
		flushes.Flushes[i].Note = temp.Flushes[i].Note
		flushes.Flushes[i].ID = temp.Flushes[i].ID
	}
	log.Println("temporary flush struct: ", temp)
	log.Println("Flushes: ", flushes)
	if len(flushes.Flushes) > 0 {
		ctx.SetState("skip", skip+3)
	}
	if temp.MoreDataAvailable {
		return flushes.Flushes, true, nil
	}
	return flushes.Flushes, false, nil
}

func ChangePass(newPass string, currentCreds string) error {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return err
	}
	body := []byte(fmt.Sprintf(`
	{
		"username": "placeholder",
		"password": "%s"
	}`, newPass))
	req, err := http.NewRequest("PUT", apiUrl+"/pass_change", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+currentCreds)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer CloseBody(resp)
	if resp.StatusCode >= 400 {
		return errors.New("failed to change password")
	}
	return nil
}

func RemoveAccount(currentCreds string) error {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", apiUrl+"/user", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+currentCreds)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer CloseBody(resp)
	if resp.StatusCode >= 400 {
		return errors.New("failed to remove user account")
	}
	return nil
}

func RemoveFlush(ID string, currentCreds string) error {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", apiUrl+"/flush/"+ID, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+currentCreds)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer CloseBody(resp)
	if resp.StatusCode >= 400 {
		return errors.New("failed to remove flush")
	}
	return nil
}

func GetStats(ctx app.Context) (FlushStatsInt, error) {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return FlushStatsInt{}, err
	}
	req, err := http.NewRequest("GET", apiUrl+"/stats", nil)
	if err != nil {
		return FlushStatsInt{}, err
	}
	var c Creds
	ctx.GetState("creds", &c)
	req.Header.Add("Authorization", "Basic "+c.UserColonPass)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlushStatsInt{}, err
	}
	defer CloseBody(r)
	if r.StatusCode >= 400 {
		ctx.SetState("creds", Creds{LoggedIn: false}).PersistWithEncryption()
		app.Window().Set("location", "login")
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return FlushStatsInt{}, err
	}
	log.Println("raw Stats: ", string(bytes))
	var stats FlushStats
	err = json.Unmarshal(bytes, &stats)
	if err != nil {
		return FlushStatsInt{}, err
	}
	statsInt := FlushStatsInt{}
	statsInt.FlushCount = int(stats.FlushCount)
	statsInt.TotalTime = int(stats.TotalTime)
	statsInt.MeanTime = int(stats.MeanTime)
	statsInt.MeanRating = int(stats.MeanRating)
	statsInt.PhoneUsedCount = int(stats.PhoneUsedCount)
	statsInt.PercentPhoneUsed = int(stats.PercentPhoneUsed)
	return statsInt, nil
}

func GiveFeedback(creds Creds, note string) (int, error) {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("POST", apiUrl+"/feedback", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Basic "+creds.UserColonPass)
	q := url.Values{}
	q.Add("note", note)
	req.URL.RawQuery = q.Encode()
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer CloseBody(r)
	return r.StatusCode, nil
}

func SubmitEditedFlush(creds Creds, flush Flush) (int, error) {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("PUT", apiUrl+"/flush/"+flush.ID, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Basic "+creds.UserColonPass)
	q := url.Values{}
	q.Add("time_start", flush.TimeStart.Format("2006-01-02 15:04:05"))
	q.Add("time_end", flush.TimeEnd.Format("2006-01-02 15:04:05"))
	if flush.PhoneUsed {
		q.Add("phone_used", "true")
	} else {
		q.Add("phone_used", "false")
	}
	q.Add("rating", strconv.Itoa(flush.Rating))
	q.Add("note", flush.Note)
	req.URL.RawQuery = q.Encode()
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer CloseBody(r)
	if r.StatusCode >= 400 {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("error while reading body")
			log.Println(err)
		} else {
			log.Println(string(b))
		}
	}
	return r.StatusCode, nil
}

func GetFlushByID(ctx app.Context, flushID string) (Flush, int, error) {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return Flush{}, 0, err
	}
	req, err := http.NewRequest("GET", apiUrl+"/flush/"+flushID, nil)
	if err != nil {
		return Flush{}, 0, err
	}
	var c Creds
	ctx.GetState("creds", &c)
	req.Header.Add("Authorization", "Basic "+c.UserColonPass)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return Flush{}, 0, err
	}
	defer CloseBody(r)
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return Flush{}, 0, err
	}
	var temp TempFlush
	err = json.Unmarshal(bytes, &temp)
	if err != nil {
		return Flush{}, 0, err
	}
	flush := Flush{
		Rating:    temp.Rating,
		PhoneUsed: temp.PhoneUsed,
		Note:      temp.Note,
		ID:        temp.ID,
	}
	flush.TimeStart, err = time.Parse(
		"2006-01-02T15:04:05",
		temp.TimeStart,
	)
	if err != nil {
		return Flush{}, 0, err
	}
	flush.TimeEnd, err = time.Parse(
		"2006-01-02T15:04:05",
		temp.TimeEnd,
	)
	if err != nil {
		return Flush{}, 0, err
	}
	log.Println("temporary flush struct: ", temp)
	log.Println("Flush: ", flush)
	return flush, r.StatusCode, nil
}
