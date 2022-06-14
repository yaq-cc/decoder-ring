# Go's Decoder Ring: Cloud based Secrets Management Made Easy

```go
package main

import (
    "log"

    secrets "github.com/yaq-cc/decoder-ring"
)

// Create a struct to hold your secrets
type Secrets struct {
	AccountSID string `secrets:"TWILIO_ACCOUNT_SID"`
	AuthToken  string `secrets:"TWILIO_AUTH_TOKEN"`
}

// Load your secrets 
func main() {
    var s Secrets
	l := secrets.NewLoader(&s)
	l.With(secrets.GoogleCloudLoader)
	err := l.Load()
	if err != nil {
		log.Fatal(err)
	}
	t.Println(s)
}
```
