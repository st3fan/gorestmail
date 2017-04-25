// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package restmail_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/textproto"
	"testing"

	"github.com/st3fan/restmail"
)

func randomHexString(length int) (string, error) {
	data := make([]byte, length)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	return "gorestmailtest-" + hex.EncodeToString(data), nil
}

func randomAccountName() (string, error) {
	randomString, err := randomHexString(10)
	if err != nil {
		return "", nil
	}
	return "gorestmail-" + randomString, nil
}

func sendEmailCommand(c *textproto.Conn, expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := c.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}
	c.StartResponse(id)
	defer c.EndResponse(id)
	code, msg, err := c.ReadResponse(expectCode)
	return code, msg, err
}

func sendEmail(from, account, subject, body string, headers map[string]string) error {
	conn, err := textproto.Dial("tcp", "restmail.net:smtp")
	if err != nil {
		return err
	}

	defer conn.Close()

	// Parse 220 restmail.net

	_, _, err = conn.ReadCodeLine(220)
	if err != nil {
		return err
	}

	// Send HELO

	_, _, err = sendEmailCommand(conn, 250, "HELO %s", "localhost")
	if err != nil {
		return err
	}

	_, _, err = sendEmailCommand(conn, 250, "MAIL FROM:<%s>", from)
	if err != nil {
		return err
	}

	_, _, err = sendEmailCommand(conn, 250, "RCPT TO:<%s@restmail.net>", account)
	if err != nil {
		return err
	}

	_, _, err = sendEmailCommand(conn, 354, "DATA")
	if err != nil {
		return err
	}

	// Send body

	w := conn.DotWriter()

	fmt.Fprintln(w, "From: "+from)
	fmt.Fprintln(w, "Subject: "+subject)
	for name, value := range headers {
		fmt.Fprintf(w, "%s: %s\n", name, value)
	}
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, body)

	if err := w.Close(); err != nil {
		return err
	}

	_, _, err = sendEmailCommand(conn, 250, "QUIT")
	if err != nil {
		// Restmail does not send a proper response to QUIT :-/
		if err.Error() == "short response: 250" {
			return nil
		}
		return err
	}

	return nil
}

func sendTestMessages(accountName string) error {
	headers1 := map[string]string{"X-Hello": "Hello, one!"}
	if err := sendEmail("example@example.com", accountName, "This is message one", "And this is body one.", headers1); err != nil {
		return err
	}
	headers2 := map[string]string{"X-Hello": "Hello, two!"}
	if err := sendEmail("example@example.com", accountName, "This is message two", "And this is body two.", headers2); err != nil {
		return err
	}
	return nil
}

func Test_RestMailClient(t *testing.T) {
	client := restmail.NewClient()

	accountName, err := randomAccountName()
	if err != nil {
		t.Error("Cannot randomAccountName(): ", err)
	}

	// Send some test messages to a new account with unique name

	if err := sendTestMessages(accountName); err != nil {
		t.Error("Cannot sendTestMessages(): ", err)
	}

	// Get the messages. Last two should be the ones we just sent.

	messages, err := client.GetMessages(accountName)
	if err != nil {
		t.Error("Cannot get messages: ", err)
	}

	if len(messages) != 2 {
		t.Error("len(messages) != 2")
	}

	if messages[0].Subject != "This is message one" {
		t.Error("Unexpected messages[0].Subject")
	}

	if messages[0].Headers["x-hello"] != "Hello, one!" {
		t.Error("Unexpected messages[0].Headers[x-hello]", messages[0].Headers["x-hello"])
	}

	if messages[0].Text != "And this is body one.\n" {
		t.Error("Unexpected messages[0].Text: ", messages[0].Text)
	}

	if messages[1].Subject != "This is message two" {
		t.Error("Unexpected messages[1].Subject")
	}

	if messages[1].Headers["x-hello"] != "Hello, two!" {
		t.Error("Unexpected messages[1].Headers[x-cheese]: ", messages[1].Headers["x-hello"])
	}

	if messages[1].Text != "And this is body two.\n" {
		t.Error("Unexpected messages[1].Text: ", messages[1].Text)
	}

	// Delete the account

	if err := client.DeleteAccount(accountName); err != nil {
		t.Error("Cannot delete account: ", err)
	}

	// Get the messages. There should be none.

	messages, err = client.GetMessages(accountName)
	if err != nil {
		t.Error("Cannot get messages: ", err)
	}

	if len(messages) != 0 {
		t.Error("len(messages) != 2")
	}
}
