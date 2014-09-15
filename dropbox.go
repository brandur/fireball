package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	DropboxContentApiUrl = "https://api-content.dropbox.com"
)

type Dropbox struct {
	AccessToken string
}

func (d *Dropbox) GetContents(file string) (string, error) {
	client := http.Client{}

	request, err := http.NewRequest("GET", DropboxContentApiUrl+"/1/files/auto/"+file, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", "Bearer "+d.AccessToken)

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Unexpected status code: %v (expected: %v)",
			resp.StatusCode, 200)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
