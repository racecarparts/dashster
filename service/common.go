package service

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/cenkalti/backoff/v4"
)

// run a command using the shell; no need to split args
// from https://stackoverflow.com/questions/6182369/exec-a-shell-command-in-go
func runcmd(cmd string, shell bool) []byte {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Printf("error executing command: '%s'", cmd)
		}
		return out
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Println(err)
	}
	return out
}

func getRequestWithBackoff(url string) ([]byte, error) {
	resp := []byte{}
	op := func() error {
		var err error
		resp, err = getRequest(url)
		if err != nil {
			return err
		}
		return nil
	}
	err := backoff.Retry(op, backoff.NewExponentialBackOff())
	if err != nil {
		return []byte{}, err
	}

	return resp, nil
}

func getRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil && resp.StatusCode != http.StatusOK {
		err = errors.New("HTTP status code: " + strconv.Itoa(resp.StatusCode) + " : " + resp.Status)
		return []byte{}, err
	}
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return body, nil
}

func getRequestBearerAuth(url string, authToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+authToken)
	return getRequestWithAuthHeader(req)
}

func getRequestBasicAuth(url string, authToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuthToken(authToken))
	return getRequestWithAuthHeader(req)
}

func getRequestWithAuthHeader(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil && resp.StatusCode != http.StatusOK {
		err = errors.New("HTTP status code: " + strconv.Itoa(resp.StatusCode) + " : " + resp.Status)
		return []byte{}, err
	}
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return body, nil
}

func basicAuthToken(token string) string {
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func padRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func padLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}
