package service

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
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

func getRequestBasicAuth(url string, authToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuthToken(authToken))

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
