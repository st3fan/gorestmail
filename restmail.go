// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package restmail

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Endpoint string
}

type Message struct {
	Text    string                 `json:"text"`
	Subject string                 `json:"subject"`
	Headers map[string]interface{} `json:"headers"`
}

func NewClient() *Client {
	return &Client{Endpoint: "https://restmail.net"}
}

func (c *Client) DeleteAccount(account string) error {
	client := http.DefaultClient

	req, err := http.NewRequest("DELETE", c.Endpoint+"/mail/"+account, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetMessages(account string) ([]Message, error) {
	req, err := http.NewRequest("GET", c.Endpoint+"/mail/"+account, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var messages []Message
	if err = json.Unmarshal(body, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
