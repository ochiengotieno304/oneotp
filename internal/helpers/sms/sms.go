package sms

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	config "github.com/ochiengotieno304/oneotp/internal/configs"
)

func SendSMS(message string, phone string) {

	configs, err := config.LoadConfig()
	if err != nil {
		return
	}

	data := url.Values{}
	data.Set("username", configs.Username)
	data.Set("to", phone)
	data.Set("from", configs.ShortCode)
	data.Set("message", message)

	urlString := configs.ATSMSEndpoint

	payload := strings.NewReader(data.Encode())

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodPost, urlString, payload)

	if err != nil {
		fmt.Println(err, "1")
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("apiKey", configs.ATAPIKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err, "2")
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err, "3")
		return
	}
	fmt.Println(string(body))
}
