package flush

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func GetFlushes(ctx app.Context) ([]Flush, error) {
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
		Rating    int       `json:"rating"`
		PhoneUsed bool      `json:"phone_used"`
		Note      string    `json:"note"`
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
	}
	log.Println("temporary flush struct: ", temp)
	log.Println("Flushes: ", flushes)
	return flushes, nil
}
