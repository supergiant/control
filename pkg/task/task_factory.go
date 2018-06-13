package task

import (
	"github.com/RichardKnop/machinery/v1/tasks"
	"fmt"
	"github.com/satori/go.uuid"
)

type Option func(signature *tasks.Signature) (*tasks.Signature, error)

func CreateTask(name string, opts ...Option) (*tasks.Signature, error) {
	t := &tasks.Signature{
		Name: name,
		Args: []tasks.Arg{},
	}

	signatureID := uuid.NewV4()
	t.UUID = fmt.Sprintf("task_%v", signatureID)

	var err error
	for _, opt := range opts {
		t, err = opt(t)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func WithHostPort(host string, port int) Option {
	return func(sig *tasks.Signature) (*tasks.Signature, error) {
		sig.Args = append(sig.Args,
			tasks.Arg{
				Name:  "host",
				Value: host,
				Type:  "string",
			}, tasks.Arg{
				Name:  "port",
				Value: port,
				Type:  "int",
			})
		return sig, nil
	}
}

func WithScript(script string) Option {
	return func(sig *tasks.Signature) (*tasks.Signature, error) {
		sig.Args = append(sig.Args,
			tasks.Arg{
				Name:  "script",
				Value: script,
				Type:  "string",
			},
		)
		return sig, nil
	}
}

func WithSSHCert(cert []byte) Option {
	return func(sig *tasks.Signature) (*tasks.Signature, error) {
		sig.Args = append(sig.Args,
			tasks.Arg{
				Name:  "script",
				Value: string(cert), // machinery can't accept []byte
				Type:  "string",
			},
		)
		return sig, nil
	}
}