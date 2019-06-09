package feat

import (
	"fmt"
	"io"
)

func info(f *Feats) io.Reader {
	r, w := io.Pipe()
	go func(f *Feats, w io.WriteCloser) {
		fmt.Fprintln(w, "mnqnctl makes local development of microservices within k8s simple")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "For more information - https://github.com/kostkobv/mannequin")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Available commands:")
		fmt.Fprintln(w)

		f.FeatsInfo(w)

		w.Close()
	}(f, w)

	return r
}

// NewMnqnctlFeats constructor.
func NewMnqnctlFeats(fds ...FeatDoer) (*Feats, error) {
	return NewFeats("mnqn", info, fds...)
}
