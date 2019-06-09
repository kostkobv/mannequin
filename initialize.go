package mannequin

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
)

const (
	configPath = ".mnqn"
	configFile = "config.yaml"

	pathSep = os.PathSeparator
)

// ConfigFolderPath returns expected path for the configuration folder.
func ConfigFolderPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir + string(pathSep) + configPath + string(pathSep), nil
}

// ConfigPath returns expected path for the configuration file.
func ConfigPath() (string, error) {
	folderPath, err := ConfigFolderPath()
	if err != nil {
		return "", err
	}

	return folderPath + configFile, err
}

// CheckInited configuration for the whole client.
func CheckInited() error {
	absFolderPath, err := ConfigFolderPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(absFolderPath); os.IsNotExist(err) {
		return errors.New("config folder does not exist")
	}

	absPath, err := ConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return errors.New("config file does not exist")
	}

	return nil
}

// Initialize the client.
func Initialize(out io.Writer) error {
	absFolderPath, err := ConfigFolderPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(absFolderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(absFolderPath, 0775); err != nil {
			return fmt.Errorf("couldn't create config folder: %s", err)
		}
	}

	absPath, err := ConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		var f *os.File
		if f, err = os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE, 0775); err != nil {
			return fmt.Errorf("couldn't create config file: %s", err)
		}
		defer f.Close()

		mnqn, err := New(out, nil)
		if err != nil {
			return err
		}

		return mnqn.Config.Save(f)
	}

	return nil
}
