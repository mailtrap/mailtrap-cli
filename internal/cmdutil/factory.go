package cmdutil

import (
	"io"
	"os"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/config"
)

type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

type Factory struct {
	Config         func() *config.Config
	IOStreams       *IOStreams
	ClientOverride *client.Client // for testing
}

func NewFactory() *Factory {
	return &Factory{
		Config: func() *config.Config {
			return config.Load()
		},
		IOStreams: &IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		},
	}
}

func (f *Factory) NewClient() (*client.Client, error) {
	if f.ClientOverride != nil {
		return f.ClientOverride, nil
	}
	token, err := config.RequireAPIToken()
	if err != nil {
		return nil, err
	}
	return client.New(token), nil
}
