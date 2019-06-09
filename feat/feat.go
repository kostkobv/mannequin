package feat

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kostkobv/mannequin"
)

// FeatDoer represents main feature interface to interact with.
type FeatDoer interface {
	Name() string
	Do(c mannequin.Mnqn, args ...string) error
	Info() io.Reader
}

// Feats is a list of registered features.
// Used to execute commands and to print out the available info.
type Feats struct {
	fs   map[string]FeatDoer
	name string
	info func(*Feats) io.Reader
}

// Control on compile level if Feats implements FeatDoer.
var _ FeatDoer = (*Feats)(nil)

// NewFeats is a constructor for Feats.
// Also another way to register FeatDoer.
func NewFeats(name string, info func(*Feats) io.Reader, fds ...FeatDoer) (*Feats, error) {
	switch {
	case name == "":
		return nil, errors.New("name is required")
	case info == nil:
		return nil, errors.New("info is required")
	}
	f := Feats{name: name, fs: map[string]FeatDoer{}, info: info}
	for _, fd := range fds {
		f.fs[fd.Name()] = fd
	}

	return &f, nil
}

// Name of the Feats.
func (f *Feats) Name() string {
	return "mnqn"
}

// Do handles the call.
func (f *Feats) Do(c mannequin.Mnqn, args ...string) error {
	// no arguments - print the info.
	if len(args) == 0 {
		io.Copy(os.Stdout, f.Info()) // nolint: errcheck
		return nil
	}

	fname := args[0]
	for n, fd := range f.fs {
		if n == fname {
			return fd.Do(c, args[1:]...)
		}
	}

	fmt.Fprintf(&c, "Unknown command \"%s\"\n", fname)

	// print info since we didn't find matching feature.
	io.Copy(&c, f.Info()) // nolint: errcheck
	return nil
}

// Register the FeatDoer.
func (f *Feats) Register(fd FeatDoer) error {
	f.fs[fd.Name()] = fd
	return nil
}

// ByName returns FeatDoer by the provided name.
// Returns error in case if FeatDoer with provided name is not registered.
func (f *Feats) ByName(name string) (FeatDoer, error) {
	fd, ok := f.fs[name]
	if !ok {
		return nil, errors.New("feature is not registered")
	}

	return fd, nil
}

// Info returns reader with the available information on the registered features.
func (f *Feats) Info() io.Reader {
	return f.info(f)
}

func (f *Feats) FeatsInfo(w io.Writer) {
	for _, fd := range f.fs {
		fmt.Fprintf(w, "%s\t\t", fd.Name())
		_, err := io.Copy(w, fd.Info())
		if err != nil {
			fmt.Fprintf(w, "Couldn't read info of \"%s\": %s", fd.Name(), err.Error()) // nolint: errcheck
		}
		fmt.Fprintln(w)
	}
}
