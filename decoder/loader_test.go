package decoder

import (
	"os"
	"testing"
)

type Secrets struct {
	AccountSID string `secrets:"TWILIO_ACCOUNT_SID"`
	AuthToken  string `secrets:"TWILIO_AUTH_TOKEN"`
}

func TestSecretsEnvVarLoader(t *testing.T) {
	var s Secrets
	os.Setenv("TWILIO_ACCOUNT_SID", "test-sid")
	os.Setenv("TWILIO_AUTH_TOKEN", "test-token")
	l := NewLoader(&s)
	l.With(EnvironmentVariableLoader)
	err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestSecretsGCPLoader(t *testing.T) {
	// MUST export PROJECT_ID=your-project
	var s Secrets
	l := NewLoader(&s)
	l.With(GoogleCloudLoader)
	err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestSecretsGCPLoaderWithOptions(t *testing.T) {
	// MUST export PROJECT_ID=your-project
	var s Secrets
	gcp := NewGoogleCloudLoader(WithProject("holy-diver-297719"))
	l := NewLoader(&s)
	l.With(gcp)
	err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}
