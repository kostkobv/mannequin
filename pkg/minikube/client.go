package minikube

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

var (
	checkVer = regexp.MustCompile(`minikube version: (v\d+.\d+.\d+)`)
	checkRun = regexp.MustCompile(`host: Running\nkubelet: Running\napiserver: Running\nkubectl: Correctly Configured: pointing to minikube-vm at.*`)
)

func CheckInstalled() (string, error) {
	cmd := exec.Command("minikube", "version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	res := checkVer.FindStringSubmatch(string(out))
	if len(res) < 1 {
		return "", errors.New("minikube is not installed")
	}

	return res[1], nil
}

func CheckRunning() error {
	cmd := exec.Command("minikube", "status")
	out, err := cmd.Output()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			return err
		}

		return fmt.Errorf("check if minikube is running: %s (error code %d)", ee.Stderr, ee.ExitCode())
	}

	res := checkRun.FindStringSubmatch(string(out))
	if len(res) < 1 {
		return errors.New("check if minikube is running")
	}

	return nil
}
