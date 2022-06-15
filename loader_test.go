package main

import (
	"os"
	"testing"

	gcp "github.com/yaq-cc/decoder-ring/gcp-decoder"
	"github.com/yaq-cc/decoder-ring/loader"
)

type Secrets struct {
	AccountSID string `secrets:"TWILIO_ACCOUNT_SID"`
	AuthToken  string `secrets:"TWILIO_AUTH_TOKEN"`
}

func TestSecretsEnvVarLoader(t *testing.T) {
	var s Secrets
	os.Setenv("TWILIO_ACCOUNT_SID", "test-sid")
	os.Setenv("TWILIO_AUTH_TOKEN", "test-token")
	l := loader.NewLoader(&s)
	l.With(loader.EnvironmentVariableLoader)
	err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestSecretsGCPLoader(t *testing.T) {
	// MUST export PROJECT_ID=your-project
	var s Secrets
	l := loader.NewLoader(&s)
	l.With(gcp.GoogleCloudLoader)
	err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestSecretsGCPLoaderWithOptions(t *testing.T) {
	// MUST export PROJECT_ID=your-project
	var s Secrets
	gl := gcp.NewGoogleCloudLoader(gcp.WithProject("holy-diver-297719"))
	l := loader.NewLoader(&s)
	l.With(gl)
	err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}
