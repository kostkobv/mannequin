package mannequin

import (
	"errors"
	"os"
)

// FileExists returns true if passed file exists.
func FileExists(path string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	i, err := os.Stat(wd + string(os.PathSeparator) + path)
	if err != nil {
		return err
	}

	if i.IsDir() {
		return errors.New("path points to folder")
	}
	return nil
}

// DirExists returns true if passed dir exists.
func DirExists(path string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	i, err := os.Stat(wd + string(os.PathSeparator) + path)
	if err != nil {
		return err
	}

	if !i.IsDir() {
		return errors.New("path points to file")
	}
	return nil
}
