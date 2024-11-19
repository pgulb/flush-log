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

func GetFlushes(ctx app.Context, skip int) ([]Flush, error) {
	apiUrl, err := GetApiUrl()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", apiUrl+"/flushes", nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer CloseBody(r)
	if r.StatusCode >= 400 {
		ctx.SetState("creds", Creds{LoggedIn: false}).PersistWithEncryption()
		app.Window().Set("location", "login")
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	temp := []struct {
		TimeStart string `json:"time_start"`
		TimeEnd   string `json:"time_end"`
		Rating    int    `json:"rating"`
		PhoneUsed bool   `json:"phone_used"`
		Note      string `json:"note"`
		ID        string `json:"_id"`
	}{}
	err = json.Unmarshal(bytes, &temp)
	if err != nil {
		return nil, err
	}
	flushes := make([]Flush, len(temp))
	for i := range temp {
		flushes[i].TimeStart, err = time.Parse("2006-01-02T15:04:05", temp[i].TimeStart)
		if err != nil {
			return nil, err
		}
		flushes[i].TimeEnd, err = time.Parse("2006-01-02T15:04:05", temp[i].TimeEnd)
		if err != nil {
			return nil, err
		}
		flushes[i].Rating = temp[i].Rating
		flushes[i].PhoneUsed = temp[i].PhoneUsed
		flushes[i].Note = temp[i].Note
		flushes[i].ID = temp[i].ID
	}
	log.Println("temporary flush struct: ", temp)
	log.Println("Flushes: ", flushes)
	if len(flushes) > 0 {
		ctx.SetState("skip", skip+3)
	}
	return flushes, nil
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
