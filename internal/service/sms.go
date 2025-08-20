package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type smsResponse struct {
	SMSMessageData struct {
		Recipients []struct {
			Status string `json:"status"`
		} `json:"Recipients"`
	} `json:"SMSMessageData"`
}

func SendSMS(to, message string) error {
	apiKey := os.Getenv("AFRICASTALKING_API_KEY")
	username := os.Getenv("AFRICASTALKING_USERNAME")

	if apiKey == "" || username == "" {
		return fmt.Errorf("missing Africa's Talking credentials in environment variables")
	}

	baseURL := os.Getenv("AFRICASTALKING_URL")

	form := url.Values{}
	form.Set("username", username)
	form.Set("to", to)
	form.Set("message", message)
	form.Set("enqueue", "1")

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var smsResp smsResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
		return err
	}

	if len(smsResp.SMSMessageData.Recipients) == 0 {
		return fmt.Errorf("no recipients in response")
	}
	if smsResp.SMSMessageData.Recipients[0].Status != "Success" {
		return fmt.Errorf("failed to send SMS: %s", smsResp.SMSMessageData.Recipients[0].Status)
	}

	return nil
}
