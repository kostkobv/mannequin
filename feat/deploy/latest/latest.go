package latest

import (
	"io"
	"strings"

	"github.com/kostkobv/mannequin"
)

// Latest deployment feature.
type Latest struct{}

// New is a constructor for New.
func New() *Latest {
	return &Latest{}
}

// Name impl.
func (l *Latest) Name() string {
	return "latest"
}

// Do impl.
func (l *Latest) Do(c mannequin.Mnqn, args ...string) error {
	return nil
}

// Info impl.
func (l *Latest) Info() io.Reader {
	return strings.NewReader("Before deploying, fetch the latest versions for all the project dependencies and redeploy it")
}
