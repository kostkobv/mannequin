package mannequin

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents data that is persisted and used later.
type Config struct {
	Version  string    `yaml:"version,flow"`
	Projects []Project `yaml:"projects,flow"`
	file     *os.File
}

// NewConfig is a constructor for Config.
func NewConfig(ver string) (Config, error) {
	if ver == "" {
		return Config{}, errors.New("version is required")
	}

	return Config{Version: ver}, nil
}

// NewConfigFromFile reads Config from .yaml file provided in as the file.
func NewConfigFromFile(file *os.File) (Config, error) {
	if file == nil {
		return Config{}, errors.New("file is required")
	}

	c := Config{file: file}
	if err := yaml.NewDecoder(file).Decode(&c); err != nil {
		return Config{}, err
	}

	return c, nil
}

// Register the Project.
func (c *Config) Register(p Project) error {
	if err := p.Validate(); err != nil {
		return err
	}

	if c.Projects != nil {
		for _, pp := range c.Projects {
			if pp.Name == p.Name {
				if pp.Path == p.Path {
					return fmt.Errorf("project with name \"%s\" is already registered", p.Name)
				}

				return fmt.Errorf("project \"%s\" is already registered but on different path (%s): re-register the project or move it back to the original path", pp.Name, pp.Path)
			}
		}
	}

	c.Projects = append(c.Projects, p)

	return c.Save(c.file)
}

// Validate the Config.
func (c *Config) Validate() error {
	if c.Version == "" {
		return errors.New("version is required")
	}
	return nil
}

// Save the Config to the provided path as a .yaml file.
// If file is not provided but Config was created from file -
// config would be written to the file.
func (c *Config) Save(file *os.File) error {
	if file == nil && c.file != nil {
		file = c.file
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

// Project data that is registered for further use.
type Project struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// NewProject is a constructor.
func NewProject(path string) (Project, error) {
	if path == "" {
		return Project{}, errors.New("path is required")
	}

	return Project{Path: path}, nil
}

// Validate the Project.
func (p *Project) Validate() error {
	switch {
	case p.Name == "":
		return errors.New("project name is required")
	case p.Path == "":
		return errors.New("path is required")
	}

	return nil
}
