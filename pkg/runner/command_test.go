package runner

import (
	"testing"
	"context"
	"bytes"
	"io"
)

func TestNewCommand(t *testing.T) {
	testCases := []struct{
		ctx context.Context
		w1 io.Writer
		w2 io.Writer
		script string
		err error
	}{
		{
			ctx: context.Background(),
			w1: &bytes.Buffer{},
			w2: &bytes.Buffer{},
			script: "echo 'hello, world'",
			err: nil,
		},
		{
			ctx: context.Background(),
			w1: &bytes.Buffer{},
			w2: nil,
			script: "echo 'hello, world'",
			err: ErrNilWriter,
		},
		{
			ctx: nil,
			err:ErrNilContext,
		},
	}
	
	for _, testCase := range testCases {
		cmd, err := NewCommand(testCase.ctx, testCase.script, testCase.w1, testCase.w2)

		if err != testCase.err {
			t.Errorf("expected error %v actual %v", testCase.err, err)
		}

		if err == nil && cmd == nil {
			t.Errorf("cmd must not be nil")
		}
	}
}
