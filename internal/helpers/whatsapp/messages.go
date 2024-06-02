package whatsapp_helper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	config "github.com/ochiengotieno304/oneotp/internal/configs"
	"github.com/ochiengotieno304/oneotp/internal/types"

)

func SendTextMessage(rq *types.SendWhatsappTextMessageRequest) (*types.SendWhatsappTextMessageResponse, error) {

	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/messages", config.FacebookGraphEndpoint, config.PhoneNumberID)
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
		"messaging_product": "whatsapp",
		"recipient_type": "individual",
		"to": "%s",
		"type": "text",
		"text": {
			"body": "%s"
		}
	}`, rq.Recipient, rq.Message))

	fmt.Println(payload)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println("ERR 1", err)
		return nil, err
	}

	bearerToken := fmt.Sprintf("Bearer %s", config.WhatsappToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearerToken)

	var fbRes types.SendWhatsappTextMessageResponse

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("ERR 2", err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("ERR 3", err)

		return nil, err
	}

	fmt.Println(string(body))

	err = json.Unmarshal(body, &fbRes) // unmarshal meta response into fbRes

	if err != nil {
		return nil, err
	}

	return &fbRes, nil
}
