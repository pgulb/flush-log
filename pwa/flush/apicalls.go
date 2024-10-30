package flush

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func GetFlushesFromApi(ctx app.Context) (string, error) {
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
