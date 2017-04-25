# restmail - Go Package to talk to restmail.net

[![Build Status](https://travis-ci.org/st3fan/restmail.svg?branch=master)](https://travis-ci.org/st3fan/restmail)

*Stefan Arentz, April 2017*

This is a Go client for [restmail.net](https://restmail.net). This is
useful for situations where you want to test a live email flow.

## Example usage

```
client := restmail.NewClient()
if messages, err := client.GetMessages("myAccountName"); err == nil {
   for _, message := range messages {
      if strings.Contains(message.Subject, "Welcome to Fiz Buzr") {
         signupLink := parseConfirmationLinkFromEmail(message.Text)
         // Profit
      }
   }
}
```
