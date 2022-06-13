package secretsloader

import (
	"context"
	"errors"
	"log"
	"os"
	"reflect"
	"strings"

	sm "cloud.google.com/go/secretmanager/apiv1"
	smpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

const (
	SecretTag = "secrets"
)

var (
	ErrNotAddressable   = errors.New("reflection: Input not addressable")
	ErrNotSettable      = errors.New("reflection: Input not settable")
	ErrNotMutable       = errors.New("reflection: Input not mutable")
	ErrTypeNotSupported = errors.New("reflection: Input type not supported")
	ErrSecretLoaderErr  = errors.New("SecretLoader error")
)

var (
	EnvironmentVariableLoader SecretLoader = NewEnvVarLoader()
	GoogleCloudLoader         SecretLoader = NewGoogleCloudLoader()
)

type Loader struct {
	Struct any
	Loader SecretLoader
}

func NewLoader(s any) *Loader {
	return &Loader{
		Struct: s,
	}
}

func (l *Loader) With(loader SecretLoader) *Loader {
	l.Loader = loader
	return l
}

func (l *Loader) Load() error {
	v := reflect.ValueOf(l.Struct)
	t := reflect.TypeOf(l.Struct)
	if v.Kind() != reflect.Pointer {
		return ErrNotMutable
	}
	v = v.Elem()
	t = t.Elem()
	if !v.CanAddr() {
		return ErrNotAddressable
	}
	if !v.CanSet() {
		return ErrNotSettable
	}
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		tag := ft.Tag.Get(SecretTag)
		if tag == "" {
			continue
		}
		fv := v.Field(i)

		switch sfk := fv.Kind(); sfk {
		case reflect.String:
			val, err := l.Loader.GetString(tag)
			if err != nil {
				return ErrSecretLoaderErr
			}
			fv.SetString(val)
		case reflect.Slice:
			// If it's a slice, we need to check it's element kind (Slice of what?)
			st := fv.Type()
			switch sk := st.Elem().Kind(); sk {
			case reflect.Uint8:
				val, err := l.Loader.GetBytes(tag)
				if err != nil {
					return ErrSecretLoaderErr
				}
				fv.SetBytes([]byte(val))
			default:
				return ErrTypeNotSupported
			}
		default:
			return ErrTypeNotSupported
		}
	}
	return nil
}

type SecretLoader interface {
	GetString(string) (string, error)
	GetBytes(string) ([]byte, error)
}

type EnvVarLoader struct{}

func NewEnvVarLoader() *EnvVarLoader {
	return &EnvVarLoader{}
}

func (l *EnvVarLoader) GetString(s string) (string, error) {
	value := os.Getenv(s)
	return value, nil
}

func (l *EnvVarLoader) GetBytes(s string) ([]byte, error) {
	value := os.Getenv(s)
	return []byte(value), nil
}

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
		return []byte{}, ErrSecretLoaderErr
	}
	value := resp.Payload.GetData()
	return value, nil

}
