package gcp_decoder

import (
	"context"
	"log"
	"os"
	"strings"

	decoder "github.com/yaq-cc/decoder-ring/decoder"

	sm "cloud.google.com/go/secretmanager/apiv1"
	smpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var GoogleCloudLoader = NewGoogleCloudLoader()

type GCPLoader struct {
	Context   context.Context
	ProjectID string
	Client    *sm.Client
}

type GCPLoaderOption func(*GCPLoader)

func FromEnvironment(key string) GCPLoaderOption {
	return func(l *GCPLoader) {
		projectID := os.Getenv(key)
		l.ProjectID = projectID
	}
}

func WithProject(projectId string) GCPLoaderOption {
	return func(l *GCPLoader) {
		l.ProjectID = projectId
	}
}

func NewGoogleCloudLoader(opts ...GCPLoaderOption) *GCPLoader {
	ctx := context.Background()
	client, err := sm.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	loader := &GCPLoader{
		Context: ctx,
		Client:  client,
	}
	if opts == nil {
		opt := FromEnvironment("PROJECT_ID")
		opt(loader)
	} else {
		for _, opt := range opts {
			opt(loader)
		}
	}
	return loader
}

func (l *GCPLoader) GetString(s string) (string, error) {
	b, err := l.GetBytes(s)
	return string(b), err
}

func (l *GCPLoader) GetBytes(s string) ([]byte, error) {
	var secretName strings.Builder
	secretName.WriteString("projects/")
	secretName.WriteString(l.ProjectID)
	secretName.WriteString("/secrets/")
	secretName.WriteString(s)
	secretName.WriteString("/versions/")
	secretName.WriteString("latest")

	req := &smpb.AccessSecretVersionRequest{
		Name: secretName.String(),
	}
	resp, err := l.Client.AccessSecretVersion(l.Context, req)
	if err != nil {
		log.Println(err)
		return []byte{}, decoder.ErrSecretLoaderErr
	}
	value := resp.Payload.GetData()
	return value, nil

}
