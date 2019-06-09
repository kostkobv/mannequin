package version

import (
	"fmt"
	"io"
	"strings"

	"github.com/kostkobv/mannequin"
)

type Version struct{}

func New() *Version {
	return &Version{}
}

func (v *Version) Name() string {
	return "version"
}

func (v *Version) Do(c mannequin.Mnqn, args ...string) error {
	_, err := fmt.Fprintf(c, "Mannequin version: %s", c.Version)
	return err
}

func (v *Version) Info() io.Reader {
	return strings.NewReader("Returns application version information")
}
