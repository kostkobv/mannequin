package docker

import (
	"errors"
	"fmt"
	"time"
)

// LConfig for Docker.
type LConfig struct {
	ImageName string `yaml:"image_name,flow"`
	Version   string `yaml:"-"`
	File      string `yaml:"file,omitempty,flow"`
}

// Validate the LConfig.
func (lc *LConfig) Validate() error {
	return nil
}

// ImageTag based on image name and version.
func (lc *LConfig) ImageTag() (string, error) {
	switch {
	case lc.ImageName == "":
		return "", errors.New("image name is required")
	case lc.Version == "":
		return "", errors.New("version is required")
	}

	return lc.ImageName + ":" + lc.Version, nil
}

// GenerateVer for the docker image and set it to the LConfig.
func (lc *LConfig) GenerateVer() error {
	if lc.Version != "" {
		return errors.New("version is already set")
	}

	lc.Version = time.Now().Format("2006-02-01-15-04-05")
	return nil
}

func (lc *LConfig) GenerateImageName(name string) error {
	switch {
	case lc.ImageName != "":
		return nil
	case name == "":
		return errors.New("name is required")
	}

	lc.ImageName = fmt.Sprintf(DefaultImageNameTmplt, name)

	return nil
}
