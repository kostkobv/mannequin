package deploy

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kostkobv/mannequin/pkg/docker"

	"github.com/kostkobv/mannequin/pkg/kubectl"

	"github.com/kostkobv/mannequin"
	"github.com/kostkobv/mannequin/feat"
	"github.com/kostkobv/mannequin/pkg/helm"
)

// Deploy feature.
type Deploy struct {
	SubFeats *feat.Feats
}

// New is a constructor for Deploy.
func New(fds ...feat.FeatDoer) (*Deploy, error) {
	f, err := feat.NewFeats("deploy", info, fds...)
	if err != nil {
		return nil, err
	}

	return &Deploy{SubFeats: f}, nil
}

// Name impl.
func (d *Deploy) Name() string {
	return "deploy"
}

// Do impl.
func (d *Deploy) Do(c mannequin.Mnqn, args ...string) error {
	// check if there are subfeatures called.
	if len(args) != 0 {
		if err := d.SubFeats.Do(c, args...); err != nil {
			return err
		}
	}

	fmt.Fprintln(c, "Reading local configuration.")
	lcfile, err := os.OpenFile(mannequin.DefaultLConfigFileName, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	lc, err := mannequin.NewLConfigFromFile(lcfile)
	if err != nil {
		return err
	}
	fmt.Fprintf(c, "Found local configuration for \"%s\".\n", lc.Name)

	fmt.Fprintln(c, "Checking global dependencies.")
	if err := lc.CheckGlobalDepsReady(); err != nil {
		return err
	}

	fmt.Fprintf(c, "Setting kubectl context to \"%s\".\n", c.K8SContext)
	if err := kubectl.UseContext(c.K8SContext); err != nil {
		return err
	}

	fmt.Fprintln(c, "Building image.")
	if err := lc.Docker.GenerateImageName(lc.Name); err != nil {
		return fmt.Errorf("couldn't generage image name: %s", err)
	}
	if err := lc.Docker.GenerateVer(); err != nil {
		return fmt.Errorf("couldn't generate image version: %s", err)
	}

	if err := docker.BuildImage(c, &c.LocalVars, lc.Docker); err != nil {
		return fmt.Errorf("couldn't build image: %s", err)
	}

	fmt.Fprintln(c, "Ready to deploy.")
	return helm.Deploy(c, &c.LocalVars, lc.Helm)
}

// Info impl.
func (d *Deploy) Info() io.Reader {
	return strings.NewReader("Deploys project with the configuration in the same folder via selected kubernetes context")
}

func info(f *feat.Feats) io.Reader {
	r, w := io.Pipe()
	go func(f *feat.Feats, w io.WriteCloser) {
		fmt.Fprintln(w, "deploy project with the configuration in the same folder via selected kubernetes context")
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
