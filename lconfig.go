package mannequin

import (
	"errors"
	"fmt"
	"os"

	"github.com/kostkobv/mannequin/pkg/docker"
	"github.com/kostkobv/mannequin/pkg/helm"
	"github.com/kostkobv/mannequin/pkg/kubectl"
	"github.com/kostkobv/mannequin/pkg/minikube"

	"gopkg.in/yaml.v2"
)

// DefaultLConfigFileName represents the name of the file that could be found in
// the root of the project.
const DefaultLConfigFileName = ".mnqn.yaml"

// DepType represents type of the Dep.
type DepType string

// Available DepTypes.
const (
	DepProject DepType = "project"
	DepService DepType = "service"
)

var validDepTypes = []DepType{DepProject, DepService}

// LConfig represents local configuration of the project.
type LConfig struct {
	Version string         `yaml:"version,flow"`
	Name    string         `yaml:"name,flow"`
	Docker  docker.LConfig `yaml:"docker,flow"`
	Helm    helm.LConfig   `yaml:"helm,flow"`
	Deps    Deps           `yaml:"deps,omitempty,flow"`
	file    *os.File
}

// NewConfig is a constructor for Config.
func NewLConfig(deps ...Dep) (LConfig, error) {
	for _, d := range deps {
		if err := d.Validate(); err != nil {
			return LConfig{}, fmt.Errorf("\"%s\" is not a valid dependency: %s", d.Name, err)
		}
	}

	return LConfig{Version: ver, Deps: deps}, nil
}

// NewConfigFromFile reads Config from .yaml file provided in as the file.
func NewLConfigFromFile(file *os.File) (LConfig, error) {
	if file == nil {
		return LConfig{}, errors.New("file is required")
	}

	lc := LConfig{file: file}
	if err := yaml.NewDecoder(file).Decode(&lc); err != nil {
		return LConfig{}, err
	}

	return lc, nil
}

// Validate the Config.
func (c *LConfig) Validate() error {
	switch {
	case c.Version == "":
		return errors.New("version is required")
	case c.Name == "":
		return errors.New("name is required")
	}

	if err := c.Docker.Validate(); err != nil {
		return fmt.Errorf("docker configuration is invalid: %s", err)
	}

	if err := c.Helm.Validate(); err != nil {
		return fmt.Errorf("helm configuration is invalid: %s", err)
	}

	return nil
}

// Save the LConfig to the provided path as a .yaml file.
// If file is not provided but Config was created from file -
// config would be written to the file.
func (c *LConfig) Save(file *os.File) error {
	if file == nil && c.file != nil {
		file = c.file
	}

	if err := c.Validate(); err != nil {
		return fmt.Errorf("local configuration is not valid: %s", err)
	}

	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	if err := yaml.NewEncoder(file).Encode(c); err != nil {
		return err
	}

	return file.Sync()
}

// CheckGlobalDeps if all the required services are installed.
func (lc *LConfig) CheckGlobalDeps() error {
	if _, err := kubectl.CheckInstalled(); err != nil {
		return fmt.Errorf("kubectl: %s", err)
	}

	if _, err := helm.CheckInstalled(lc.Helm.BinaryPath); err != nil {
		return fmt.Errorf("helm: %s", err)
	}

	if _, err := minikube.CheckInstalled(); err != nil {
		return fmt.Errorf("minikube: %s", err)
	}

	return nil
}

func (lc *LConfig) CheckGlobalDepsReady() error {
	if err := lc.CheckGlobalDeps(); err != nil {
		return fmt.Errorf("global dependency returned error: %s", err)
	}

	if err := minikube.CheckRunning(); err != nil {
		return fmt.Errorf("minikube: %s", err)
	}

	return nil
}

// Dep represents Dependency of the Project.
type Dep struct {
	Name    string  `yaml:"name"`
	Type    DepType `yaml:"type"`
	Prepare func() error
}

// Validate the dep.
func (d *Dep) Validate() error {
	switch {
	case d == nil:
		return errors.New("dependency is required")
	case d.Name == "":
		return errors.New("name is required")
	case d.Type == "":
		return errors.New("type is required")
	case d.Prepare == nil:
		return errors.New("prepare func is required")
	}

	var validType bool
	for _, t := range validDepTypes {
		if d.Type == t {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("%s is not a valid dep type", d.Type)
	}

	return nil
}

type Deps []Dep

// Register the Dependency.
func (ds *Deps) Register(d Dep) error {
	if err := d.Validate(); err != nil {
		return err
	}

	*ds = append(*ds, d)
	return nil
}

// Validate the Deps.
func (ds *Deps) Validate() error {
	if ds == nil {
		return errors.New("no deps to validate")
	}

	for _, d := range *ds {
		if err := d.Validate(); err != nil {
			return err
		}
	}

	return nil
}
