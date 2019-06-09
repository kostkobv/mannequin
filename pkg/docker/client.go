package docker

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kostkobv/mannequin/pkg"
)

const DefaultImageNameTmplt = "mnqn.local/%s"

// VarName of the docker variables that are provided into VarSetter.
type VarName = string

// Available docker vars.
const (
	VarDockerImageTag     VarName = "DOCKER_IMAGE_TAG"
	VarDockerImageName    VarName = "DOCKER_IMAGE_NAME"
	VarDockerImageVersion VarName = "DOCKER_IMAGE_VERSION"
)

const DefaultFilePath = "./Dockerfile"

// BuildImage and tag it using the image name and version.
// Dockerfile from provided filepath would be used.
// DefaultFilePath would be used otherwise.
// w is used to print output.
func BuildImage(w io.Writer, vars pkg.VarStorer, lc LConfig) error {
	switch {
	case w == nil:
		return errors.New("writer is required")
	case vars == nil:
		return errors.New("variable store is required")
	}

	if err := lc.Validate(); err != nil {
		return err
	}

	// check dockerfile.
	_, err := os.Stat(lc.File)
	switch {
	case os.IsNotExist(err):
		return errors.New("dockerfile cannot be found")
	case err != nil:
		return fmt.Errorf("could not check if dockerfile exists: %s", err)
	}

	// generate image tag.
	tag, err := lc.ImageTag()
	if err != nil {
		return err
	}

	// set vars.
	if err := vars.Register(VarDockerImageTag, tag); err != nil {
		return err
	}
	if err := vars.Register(VarDockerImageVersion, lc.Version); err != nil {
		return err
	}
	if err := vars.Register(VarDockerImageName, lc.ImageName); err != nil {
		return err
	}

	fmt.Fprintln(w, "Building:")
	fmt.Fprintln(w, "----------------------------------------------------")

	// run the command.
	cmd := exec.Command("docker", "build", "-t", tag, "-f", lc.File, ".")
	cmd.Stderr = w
	cmd.Stdout = w

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			return err
		}

		return fmt.Errorf("failed: error code %d: %s", ee.ExitCode(), strings.TrimSpace(string(ee.Stderr)))
	}

	fmt.Fprintln(w, "----------------------------------------------------")
	fmt.Fprintln(w, "Succefully built!")

	return nil
}
