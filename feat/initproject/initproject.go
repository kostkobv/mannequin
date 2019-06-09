package initproject

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/kostkobv/mannequin"
	"github.com/kostkobv/mannequin/pkg/docker"
)

type InitProject struct{}

// New is the Init constructor.
func New() *InitProject {
	return &InitProject{}
}

// Name implementation.
func (i *InitProject) Name() string {
	return "init"
}

// Do implementation.
func (i *InitProject) Do(c mannequin.Mnqn, args ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("couldn't get working dir: %s", err)
	}

	lc, err := i.initLConfig(c, wd, mannequin.DefaultLConfigFileName, path.Base(wd))
	if err != nil {
		return fmt.Errorf("couldn't initialize reference file: %s", err)
	}

	if err := c.Config.Register(mannequin.Project{Path: wd, Name: lc.Name}); err != nil {
		return fmt.Errorf("couldn't register project to the configuration: %s", err)
	}

	fmt.Fprintf(c, "Project \"%s\" is successfully registered\n", lc.Name)

	return nil
}

func (i *InitProject) initLConfig(c mannequin.Mnqn, wd, fname, pname string) (mannequin.LConfig, error) {
	switch {
	case wd == "":
		return mannequin.LConfig{}, errors.New("working directory is required")
	case fname == "":
		return mannequin.LConfig{}, errors.New("local configuration file name is required")
	case pname == "":
		return mannequin.LConfig{}, errors.New("project name is required")
	}

	path := wd + string(os.PathSeparator) + fname
	if _, err := os.Stat(path); os.IsNotExist(err) {
		lc, err := i.lConfigSurvey(c, wd, pname)
		if err != nil {
			return mannequin.LConfig{}, err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0775)
		if err != nil {
			return mannequin.LConfig{}, err
		}
		defer f.Close()

		return lc, lc.Save(f)
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0775)
	if err != nil {
		return mannequin.LConfig{}, err
	}
	defer f.Close()

	return mannequin.NewLConfigFromFile(f)
}

// Info implementation.
func (i *InitProject) Info() io.Reader {
	return strings.NewReader("Runs the survey to register the project " +
		"(it is presumed that the command is ran within the project root folder)")
}

func (i *InitProject) lConfigSurvey(c mannequin.Mnqn, wd, pname string) (mannequin.LConfig, error) {
	lc, err := mannequin.NewLConfig()
	if err != nil {
		return mannequin.LConfig{}, err
	}

	lc.Name = pname
	lc.Docker.File = docker.DefaultFilePath

	err = mannequin.FileExists(lc.Docker.File)
	if err != nil {
		fmt.Fprintln(c, "Dockerfile is not found")
		for {
			fmt.Fprintf(c, "Please provide path to the working Dockerfile (example: %s):\n", docker.DefaultFilePath)
			if _, err := fmt.Scanln(&lc.Docker.File); err != nil {
				fmt.Fprintf(c, "Couldn't read Dockerfile: %s\n", err)
				continue
			}

			if err := mannequin.FileExists(lc.Docker.File); err != nil {
				fmt.Fprintf(c, "Couldn't find Dockerfile: %s\n", err)
				continue
			}

			break
		}
	}

	for {
		fmt.Fprintf(c, "Please provide relative path to the Helm chart:\n")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(c, "Couldn't read input: %s\n", err)
			continue
		}
		lc.Helm.ChartPath = strings.TrimSpace(text)
		if lc.Helm.ChartPath == "" {
			fmt.Fprintln(c, "Value is required")
			continue
		}

		if err := mannequin.DirExists(lc.Helm.ChartPath); err != nil {
			fmt.Fprintf(c, "Couldn't find chart dir: %s\n", err)
			continue
		}

		break
	}
	for {
		fmt.Fprintf(c, "Please provide relative path to the Helm values (skip, if you don't need values):\n")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(c, "Couldn't read input: %s\n", err)
			continue
		}
		lc.Helm.ValuesPath = strings.TrimSpace(text)
		if lc.Helm.ValuesPath == "" {
			break
		}

		if err := mannequin.FileExists(lc.Helm.ValuesPath); err != nil {
			fmt.Fprintf(c, "Couldn't find values file: %s\n", err)
			continue
		}

		break
	}
	lc.Helm.ReleaseName = pname

	if err := lc.Validate(); err != nil {
		return mannequin.LConfig{}, err
	}

	return lc, nil
}
