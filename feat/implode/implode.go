package implode

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kostkobv/mannequin"
)

type Implode struct{}

func New() *Implode {
	return &Implode{}
}

// Name impl.
func (p *Implode) Name() string {
	return "implode"
}

// Do impl.
func (p *Implode) Do(c mannequin.Mnqn, args ...string) error {
	fmt.Fprintln(c, "Are you sure you want to purge current configuration? (y/N):")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(c, "Something went wrong: %s\n", err)
		return err
	}
	text = strings.TrimSpace(text)
	if text != "y" && text != "Y" {
		fmt.Fprintln(c, "Skipping.")
		return nil
	}

	path, err := mannequin.ConfigFolderPath()
	if err != nil {
		fmt.Fprintf(c, "Couldn't fetch configuration folder path: %s\n", err)
		return err
	}

	if err := os.RemoveAll(path); err != nil {
		return err
	}

	fmt.Fprintln(c, "Purged.")
	return nil
}

// Info impl.
func (p *Implode) Info() io.Reader {
	return strings.NewReader("Implodes previously made global configuration")
}
